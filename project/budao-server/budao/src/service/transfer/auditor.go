package transfer

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	_ "net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	sj "github.com/bitly/go-simplejson"
	"github.com/sumaig/glog"

	"common"
	"db"
	av "service/transfer/gen-go/mms_report_interface"
)

const (
	MAX_AUDIT_NUM = 100
)

//TODO 修改序列化逻辑，使用字符串，增加类型标识，来源标识，id标识
var gVSerial uint64 = 1
var gCSerial uint64 = 1

var (
	gAuditor    *Auditor
	onceAuditor sync.Once
)

func GetAuditor() *Auditor {
	onceAuditor.Do(func() {
		gAuditor = &Auditor{}
		gAuditor.initAuditor()
	})
	return gAuditor
}

type Auditor struct {
	mysqlClient *db.MysqlClient

	//用于视频和图片的审核
	videoConf       *common.AuditConf
	signParNames    []string
	signParNamesMap map[string][]string
	thriftClient    *av.MmsReportServClient
	transport       thrift.TTransport
}

func (s *Auditor) Close() {
	s.transport.Close()
}

func (s *Auditor) IgnoreVideoAudit() bool {
	return s.videoConf.IsDisable
}

func (s *Auditor) initAuditor() error {
	s.videoConf = common.GetConfig().Audit["video"]

	//init thrift client
	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	var transportFactory thrift.TTransportFactory
	transportFactory = thrift.NewTTransportFactory()
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

	var err error
	s.transport, err = thrift.NewTSocket(s.videoConf.AuditAddr)
	if err != nil {
		glog.Error("Error opening socket:", err)
		return err
	}
	s.transport, err = transportFactory.GetTransport(s.transport)
	if err != nil {
		return err
	}
	if err := s.transport.Open(); err != nil {
		return err
	}
	iprot := protocolFactory.GetProtocol(s.transport)
	oprot := protocolFactory.GetProtocol(s.transport)

	s.thriftClient = av.NewMmsReportServClient(thrift.NewTStandardClient(iprot, oprot))
	err = s.thriftClient.Ping(context.Background())
	if err != nil {
		glog.Error("ping auditVideo server failed. err:%v", err)
	} else {
		glog.Debug("ping auditVideo server ok")
	}

	//prepare signParNames
	s.signParNames = []string{
		"MmsReportReq.chid",
		"MmsReportReq.appid",
		"MmsReport.serial",
		"MmsReport.uid",
		"MmsReport.reportTime",
		"MmsReportAttc.attcType",
		//"MmsReportAttc.attcText",
		"MmsReportAttc.attcUrl",
	}
	s.signParNamesMap = map[string][]string{
		"MmsReportReq":  {"chid", "appid"},
		"MmsReport":     {"serial", "uid", "reportTime"},
		"MmsReportAttc": {"attcType", "attcText", "attcUrl"},
	}

	return nil
}

func (s *Auditor) generateVideoSign(chid, appid string, reports []*av.MmsReport) (sign string) {
	var reportStr string
	for _, report := range reports {
		reportStr += report.Serial
		reportStr += strconv.FormatInt(report.UID, 10)
		reportStr += report.ReportTime
		for _, attc := range report.Attachments {
			reportStr += attc.AttcType
			reportStr += attc.AttcUrl
		}
	}
	data := fmt.Sprintf("%s%s%s%s", chid, appid, reportStr, s.videoConf.SecretKey)
	//TODO string []btye转换优化
	arr := sha1.Sum([]byte(data))
	sign = hex.EncodeToString(arr[:])
	return
}

func (s *Auditor) getNextSerial(gSerial *uint64) (serial uint64) {
	serial = atomic.LoadUint64(gSerial)
	for {
		if ok := atomic.CompareAndSwapUint64(gSerial, serial, serial+1); !ok {
			serial = atomic.LoadUint64(gSerial)
		} else {
			break
		}
	}
	return
}

func (s *Auditor) composeVideoAuditReq(urlMap map[string]uint64) (mmsReportReq *av.MmsReportReq, serialMap map[string]uint64) {
	videoConf := s.videoConf
	serialMap = make(map[string]uint64, 0)

	mmsReportReq = av.NewMmsReportReq()

	var reports []*av.MmsReport
	reports = make([]*av.MmsReport, 0)
	for videoUrl, vid := range urlMap {
		report := av.NewMmsReport()
		//gen serial
		serial := s.getNextSerial(&gVSerial)
		report.Serial = strconv.FormatUint(serial, 10)
		serialMap[report.Serial] = vid

		report.UID = 123
		report.ReportTime = time.Now().Format("2006-01-02 15:04:05")
		//MmsReportAttc
		attachments := make([]*av.MmsReportAttc, 0)
		attachment := av.NewMmsReportAttc()
		attachment.AttcType = "VIDEO_FILE"
		attachment.AttcUrl = videoUrl
		attachments = append(attachments, attachment)
		report.Attachments = attachments
		report.ExtPar = strconv.FormatUint(vid, 10)

		reports = append(reports, report)
	}

	//sign
	mmsSign := av.NewMmsSign()
	mmsSign.AppKey = videoConf.SecretId
	mmsSign.Sign = s.generateVideoSign(videoConf.Chid, videoConf.Appid, reports)
	mmsSign.SignParNames = s.signParNames

	mmsReportReq.Chid = videoConf.Chid
	mmsReportReq.Appid = videoConf.Appid
	mmsReportReq.MmsSign = mmsSign
	mmsReportReq.Reports = reports

	return
}

func (s *Auditor) AuditVideo(urlMap map[string]uint64) (reportMap map[uint64]int32, err error) {
	glog.Debug("audit video is disable:%v", s.videoConf.IsDisable)
	if s.videoConf.IsDisable {
		reportMap = make(map[uint64]int32, 0)
		for _, vid := range urlMap {
			reportMap[vid] = 1
		}
		return
	}
	//prepare req args
	if len(urlMap) > MAX_AUDIT_NUM {
		return nil, errors.New("audit video failed. too many url.")
	}
	mmsReportReq, serialMap := s.composeVideoAuditReq(urlMap)
	glog.Debug("video audit req:[%s]", mmsReportReq.String())

	//send req
	mmsReportRsp, err := s.thriftClient.PushReports(context.Background(), mmsReportReq)
	if err != nil {
		glog.Error("audit video failed. mmsReportRsp:%v, err:%v", mmsReportRsp, err)
		return
	} else if mmsReportRsp.Code <= 0 {
		glog.Error("audit video failed. mmsReportRsp:%v", mmsReportRsp)
		err = errors.New("failed")
		return
	} else {
		glog.Debug("audit video success. mmsReportRsp:%v", mmsReportRsp)
	}

	//parse resp
	reportMap = make(map[uint64]int32, 0)
	for _, mmsReportRspRec := range mmsReportRsp.MmsReportRspRecs {
		if vid, ok := serialMap[mmsReportRspRec.Serial]; ok {
			reportMap[vid] = mmsReportRspRec.Code
		}
	}

	return
}

func (s *Auditor) AuditPicture(pictureUrl string) (err error) {
	return nil
}

type AuditContentExtPar struct {
	ContentId   string `json:"contentId"`
	ContentType string `json:"contentType"`
}

func (s *Auditor) AuditContent(contentType, content, contentId string) (status int, matchs []string, err error) {
	//prepare req args
	contentConf := common.GetConfig().Audit[contentType]

	var reqArgs = []string{"appid", "secretId", "timestamp", "random", "serial", "content", "account", "callback"}

	args := make(url.Values)
	args["appid"] = []string{contentConf.Appid}
	args["secretId"] = []string{contentConf.SecretId}
	secretKey := contentConf.SecretKey
	timestamp := time.Now().Unix()
	args["timestamp"] = []string{strconv.FormatInt(timestamp, 10)}
	random := rand.Int31()
	args["random"] = []string{strconv.Itoa(int(random))}
	serial := s.getNextSerial(&gCSerial)
	args["serial"] = []string{strconv.FormatUint(serial, 10)}
	args["content"] = []string{content}
	args["account"] = []string{contentId}
	auditContentExtPar := &AuditContentExtPar{
		ContentId:   contentId,
		ContentType: contentType,
	}
	extPar, err := json.Marshal(auditContentExtPar)
	if err != nil {
		return
	}
	//args["extPar"] = []string{string(extPar)}
	args["callback"] = []string{string(extPar)}

	var signBuf bytes.Buffer
	sort.Strings(reqArgs)
	for _, arg := range reqArgs {
		signBuf.WriteString(arg)
		signBuf.WriteString(args[arg][0])
	}
	signBuf.WriteString(secretKey)
	glog.Debug("signBuf:%s", signBuf.String())
	signStr := signBuf.String()
	signMd5 := md5.Sum([]byte(signStr))
	sign := hex.EncodeToString(signMd5[:])
	glog.Debug("sign:%s", string(sign[:]))

	args["sign"] = []string{string(sign[:])}
	glog.Debug("text audit args:%v", args)

	//post req
	resp, err := http.PostForm(contentConf.AuditAddr, args)
	if err != nil {
		glog.Debug("post err:%v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		glog.Error("resp:[%v]", resp)
		return
	}

	//parse resp
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	bodyJson, err := sj.NewJson([]byte(body))
	if err != nil {
		return
	}
	glog.Debug("bodyJson:[%v]", bodyJson)
	if bodyJson.Get("code").MustInt() != 100 {
		glog.Error("req failed. resp code:%d message:%s", bodyJson.Get("code").MustInt(), bodyJson.Get("message").MustString())
		return
	}
	result := bodyJson.Get("result")
	if result == nil {
		glog.Error("result is nil")
		return
	}
	/*
		status:
		1	正常	满足正常条件
		2	不通过	满足不通过条件
		3	待确认	满足待确认条件（根据业务配置是否需要该状态）,机器推过来的待确认暂时当成正常处理，等到人工推审核结果过来再更新文字状态，因为人工只推送审核不通过的，同时防止人工推送失败
		-1	未匹配标准	未匹配以上所有条件
		-2	任务调用失败	系统内部代码
	*/
	status = result.Get("status").MustInt()
	if status != 2 {
		return
	}

	matchs = result.Get("matchs").MustStringArray()
	glog.Debug("contentId:%s, contentType:%s, matchs:[%v]", contentId, contentType, matchs)
	taskItems := result.Get("taskItems").MustArray()
	for index, taskItemIter := range taskItems {
		taskItem := taskItemIter.(map[string]interface{})
		glog.Debug("index:%d, taskItem:[%v]", index, taskItem)
	}

	return
}

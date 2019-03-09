package transfer

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sumaig/glog"

	"common"
	"db"
	pb "twirprpc"
)

var (
	gTransfer    *Transfer
	onceTransfer sync.Once
)

func GetTransfer() *Transfer {
	onceTransfer.Do(func() {
		gTransfer = &Transfer{}
		gTransfer.initTransfer()
	})
	return gTransfer
}

//TODO 配置自己的log目录，原则上每个服务都用自己的日志目录
type Transfer struct {
	auditor             *Auditor
	idGenerator         *IdGenerator
	pushPunishCmdServer *PushPunishCmdServer
}

func (s *Transfer) initTransfer() {
	s.auditor = GetAuditor()
	s.idGenerator = GetIdGenerator()
	s.pushPunishCmdServer = GetPushPunishCmdServer()
	go s.pushPunishCmdServer.Start()
}

func (s *Transfer) Close() {
	s.pushPunishCmdServer.Close()
	s.auditor.Close()
}

func (s *Transfer) PostVideos(ctx context.Context, req *pb.PostVideosRequest) (resp *pb.PostVideosResponse, err error) {
	//parse req
	glog.Debug("begin post video ------------------\nreq:[%v]", req)

	resp = &pb.PostVideosResponse{}
	tResults := make([]*pb.PostVideoResult, 0)
	defer func() {
		resp.TResults = tResults
		if err != nil {
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	//check req arg
	for _, postVideo := range req.TVideos {
		if postVideo.SourceVid == "" ||
			postVideo.Pic == "" ||
			postVideo.VideoUrl == "" ||
			postVideo.Title == "" ||
			postVideo.Duration == 0 {
			glog.Error("bad req args. req:%v", req)

			tResult := &pb.PostVideoResult{}
			tResult.SourceVid = postVideo.SourceVid
			tResult.EVideoSource = postVideo.EVideoSource
			tResult.EVideoParseRule = postVideo.EVideoParseRule
			tResult.Result = pb.EPostResult_E_PostResult_Params_Error
			return
		}
	}

	var tVideos []*pb.PostVideo = req.GetTVideos()
	glog.Debug("video size:%d", len(tVideos))
	urlMap := make(map[string]uint64, 0)
	videoMap := make(map[uint64]*pb.PostVideo, 0)
	titleMap := make(map[uint64]*pb.AuditContentRequest, 0)
	for _, tVideo := range tVideos {
		tResult := &pb.PostVideoResult{}
		tResult.SourceVid = tVideo.SourceVid
		tResult.EVideoSource = tVideo.EVideoSource
		tResult.EVideoParseRule = tVideo.EVideoParseRule
		//1.判断重复，过滤重复视频。
		if isRepeat, err := s.ChechVideoIsRepeat(tVideo.SourceVid, int(tVideo.EVideoSource)); isRepeat {
			glog.Debug("video is repeat. sourceVid: %s, videoSource:%d. err:%v", tVideo.SourceVid, tVideo.EVideoSource, err)
			tResult.Result = pb.EPostResult_E_PostResult_Repeated
			tResults = append(tResults, tResult)
			continue
		}
		//2.生成唯一id
		var videoId, autoIncreId uint64
		videoId, autoIncreId, err = s.idGenerator.GetItemId("video_")
		if err != nil {
			glog.Error("generator video id failed. sourceVid:%s, type:%d, err:%v", tVideo.SourceVid, tVideo.EVideoSource, err)
			tResult.Result = pb.EPostResult_E_PostResult_Server_Error
			tResults = append(tResults, tResult)
			continue
		}

		//2.2 上传视频图片,替换现有url为yy云图片连接
		newUrl, err := UploadPic(tVideo.Pic)
		if err != nil {
			glog.Error("upload pic failed.  sourceVid:%s, type:%d, err:%v", tVideo.SourceVid, tVideo.EVideoSource, err)
			tResult.Result = pb.EPostResult_E_PostResult_Server_Error
			tResults = append(tResults, tResult)
			continue
		}
		tVideo.Pic = newUrl

		//3.将视频插入到mysql中，先设置状态为未审核
		state := 0
		if s.auditor.IgnoreVideoAudit() {
			//忽略视频审核时，直接审核通过
			state = 2
		}
		tableName, _ := db.GetTableName("video_", videoId)
		title, _ := db.MysqlEscapeString(tVideo.GetTitle())
		pic, _ := db.MysqlEscapeString(tVideo.GetPic())
		videoUrl, _ := db.MysqlEscapeString(tVideo.GetVideoUrl())
		execSql := fmt.Sprintf("update %s set vid=%d, source_vid=%s, title='%s', coverurl='%s', videourl='%s', duration=%d, width=%d, height=%d, type=%d, parse_type=%d, state=%d where id=%d", tableName, videoId, tVideo.GetSourceVid(), title, pic, videoUrl, tVideo.GetDuration(), tVideo.GetWidth(), tVideo.GetHeight(), tVideo.GetEVideoSource(), tVideo.GetEVideoParseRule(), state, autoIncreId)
		_, err = db.Exec(common.BUDAODB, execSql)
		if err != nil {
			glog.Error("insert video failed. execSql:[%s] err:%v", execSql, err)
			tResult.Result = pb.EPostResult_E_PostResult_Server_Error
			tResults = append(tResults, tResult)
			continue
		}
		tResult.Result = pb.EPostResult_E_PostResult_OK
		tResult.Vid = strconv.FormatUint(videoId, 10)
		tResults = append(tResults, tResult)
		urlMap[tVideo.GetVideoUrl()] = videoId
		videoMap[videoId] = tVideo
		titleMap[videoId] = &pb.AuditContentRequest{
			ContentId:   videoId,
			Content:     title,
			ContentType: "title",
		}
	}

	//4.audit video
	if len(urlMap) == 0 || s.auditor.IgnoreVideoAudit() {
		glog.Debug("len(urlMap):%d, video audit is disable:%v", len(urlMap), s.auditor.IgnoreVideoAudit())
		return
	}
	reportMap, err := s.auditor.AuditVideo(urlMap)
	if err != nil {
		glog.Error("audit video failed. urlMap:[%v], err:%v", urlMap, err)
	}

	//5.compose audit video resp
	for videoId, code := range reportMap {
		tResult := &pb.PostVideoResult{}
		tVideo, ok := videoMap[videoId]
		if !ok {
			continue
		}
		tResult.SourceVid = tVideo.SourceVid
		tResult.EVideoSource = tVideo.EVideoSource
		tResult.EVideoParseRule = tVideo.EVideoParseRule
		//检查审核结果
		if code <= 0 {
			glog.Error("audit video failed. sourceVid:%s, type:%d, code:%d", tVideo.SourceVid, tVideo.EVideoSource, code)
			tResult.Result = pb.EPostResult_E_PostResult_Server_Error
			tResults = append(tResults, tResult)
			continue
		}

		//update video state, 设置状态为审核中
		tableName, _ := db.GetTableName("video_", videoId)
		execSql := fmt.Sprintf("update %s set state=1 where vid=%d", tableName, videoId)
		_, err := db.Exec(common.BUDAODB, execSql)
		if err != nil {
			glog.Error("insert video failed. execSql:[%s] err:%v", execSql, err)
			tResult.Result = pb.EPostResult_E_PostResult_Server_Error
			tResults = append(tResults, tResult)
			continue
		}
		tResult.Result = pb.EPostResult_E_PostResult_OK
		tResults = append(tResults, tResult)
	}

	glog.Debug("end post video ------------------\nresp:[%v]", resp)

	//审核视频标题
	go s.AuditVideoTitleMulti(titleMap)

	return
}

func (s *Transfer) AuditVideoTitleMulti(titleMap map[uint64]*pb.AuditContentRequest) {
	for _, req := range titleMap {
		resp, err := s.AuditContent(context.Background(), req)
		if err != nil {
			glog.Error("audit video title failed. req:%v,resp:%v,err:%v", req, resp, err)
		}
	}
}

const SourceVidPrefix string = "sourcevid"

func (s *Transfer) ChechVideoIsRepeat(sourceVid string, eVideoSource int) (bool, error) {
	//compose key
	sourceVidNum, err := strconv.ParseUint(sourceVid, 10, 64)
	if err != nil {
		return false, err
	}
	key := fmt.Sprintf("%s_%d", SourceVidPrefix, sourceVidNum%1000)
	field := fmt.Sprintf("%s_%d", sourceVid, eVideoSource)
	if r, err := db.HSetNX(key, field, "1"); err != nil || r != 1 {
		return true, err
	}

	return false, nil
}

func (s *Transfer) AuditContent(ctx context.Context, req *pb.AuditContentRequest) (resp *pb.AuditContentResponse, err error) {
	//parse req
	glog.Debug("begin audit content------------------\nreq:[%v]", req)
	resp = &pb.AuditContentResponse{}
	resp.ContentId = req.GetContentId()
	resp.Result = pb.EPostResult_E_PostResult_Server_Error

	//audit content
	contentType := req.GetContentType()
	content := req.GetContent()
	contentId := req.GetContentId()
	contentIdStr := strconv.FormatUint(contentId, 10)
	state, matchs, err := s.auditor.AuditContent(contentType, content, contentIdStr)
	if err != nil {
		glog.Error("audit failed. err:%v", err)
		return
	}
	glog.Debug("audit content result. state:%d, matchs:%v", state, matchs)
	/*
		status:
		1	正常	满足正常条件
		2	不通过	满足不通过条件
		3	待确认	满足待确认条件（根据业务配置是否需要该状态）,机器推过来的待确认暂时当成正常处理，等到人工推审核结果过来再更新文字状态，因为人工只推送审核不通过的，同时防止人工推送失败
		-1	未匹配标准	未匹配以上所有条件
		-2	任务调用失败	系统内部代码
	*/

	if state == -1 || state == -2 {
		glog.Error("audit failed. state:%d", state)
		return
	}
	if state == 2 {
		resp.AuditResult = pb.EAuditResult_E_AuditResult_Unpass
		state = 3
	} else if state == 3 || state == 1 {
		state = 2
		resp.AuditResult = pb.EAuditResult_E_AuditResult_Pass
	}
	//update content state
	var execSql string
	if contentType == "comment" {
		tableName, _ := db.GetTableName("comment_", uint64(contentId))
		execSql = fmt.Sprintf("update %s set state=%d where cid=%d", tableName, state, contentId)
	} else if contentType == "title" {
		tableName, _ := db.GetTableName("video_", uint64(contentId))
		execSql = fmt.Sprintf("update %s set state=%d where vid=%d", tableName, state, contentId)
	} else if contentType == "question" {
		execSql = fmt.Sprintf("update question set state=%d where id=%d", state, contentId)
	}
	result, err := db.Exec(common.BUDAODB, execSql)
	if err != nil {
		glog.Error("update failed. execSql:[%s] result:%v err:%v", execSql, result, err)
		return
	}

	resp.Result = pb.EPostResult_E_PostResult_OK
	glog.Debug("end audit content ------------------resp:[%v]", resp)

	return
}

const (
	Upload_Key    = "ak_yey"
	Upload_Secret = "97cbac39d755ee2c3c63ac79944a2a477fd5faff"
	Upload_Host   = "jxzimg.bs2ul.yy.com"
)

var contentTypeMap = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
}

/**
* 上传图片
* 两种情况：1，图片路径 filepath不为空  2，上传图片文件 filepath为空
*
* @param filepath 图片路径，可以是网络图片 也可以是本地图片
* @param filename 上传图片的名字
* @param buff  上传图片的内容
 */
func UploadPic(filepath string) (newUrl string, err error) {
	var bodyBuff []byte
	var fileSuffix, filename, picType, contentType string

	if !strings.HasPrefix(strings.Trim(filepath, " "), "http") {
		err = errors.New("filepath not have http prefix")
		return
	}

	bodyBuff, picType, err = LoadHttpPic(filepath)
	if err != nil {
		glog.Error("load pic failed. filepath:%s, err:%v", filepath, err)
		return
	}
	filename = path.Base(filepath)
	fileSuffix = path.Ext(filename)
	if v, ok := contentTypeMap[fileSuffix]; ok {
		contentType = v
	} else {
		if picType == "" {
			glog.Error("can not find fileSuffix. filePath:%s", filepath)
			return
		}
		glog.Debug("find pic type from head's content-type:%s", picType)
		contentType = picType
	}
	glog.Debug("filepath:%s,filename:%s,fileSuffix:%s,contentType:%s", filepath, filename, fileSuffix, contentType)

	//加密文件名字
	fname := common.GetMd5Str(bodyBuff)
	newUrl = "http://" + Upload_Host + "/" + fname
	request, err := http.NewRequest("PUT", newUrl, bytes.NewBuffer(bodyBuff))
	if err != nil {
		glog.Error("newRequest error %v", err)
		return
	}

	//加密过程
	expires := strconv.FormatInt(time.Now().Unix()+7200, 10)
	hmacStr := "PUT\njxzimg\n" + fname + "\n" + expires + "\n"
	hm := common.GetSha1Str(Upload_Secret, hmacStr)
	base := base64.URLEncoding.EncodeToString(hm)

	request.Header.Set("Host", Upload_Host)
	request.Header.Set("Date", common.Gmtime())
	request.Header.Set("Authorization", Upload_Key+":"+base+":"+expires)
	request.Header.Set("Content-Length", strconv.Itoa(len(bodyBuff)))
	request.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		glog.Error("respstatusCode not 200")
	} else {
		_, err = ioutil.ReadAll(resp.Body)
	}

	return
}

//获取资源图片内容
func LoadHttpPic(picUrl string) (body []byte, picType string, err error) {
	resp, err := http.Get(picUrl)
	defer resp.Body.Close()
	if err != nil {
		glog.Error("respose http get error %s", err)
		return
	}
	glog.Debug("header:%v", resp.Header)
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		glog.Error("read resp body failed. err:%v, body:%s", err, resp)
	}
	if contentType, ok := resp.Header["Content-Type"]; ok {
		if len(contentType) != 0 {
			picType = resp.Header["Content-Type"][0]
		}
	}
	return
}

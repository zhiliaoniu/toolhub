package transfer

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/sumaig/glog"

	"common"
	"db"
)

/*
var DSA_PUBLIC_KEY = []byte(`-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIBuDCCASwGByqGSM44BAEwggEfAoGBAP1/U4EddRIpUt9KnC7s5Of2EbdSPO9E
AMMeP4C2USZpRV1AIlH7WT2NWPq/xfW6MPbLm1Vs14E7gB00b/JmYLdrmVClpJ+f
6AR7ECLCT7up1/63xhv4O1fnxqimFQ8E+4P208UewwI1VBNaFpEy9nXzrith1yrv
8iIDGZ3RSAHHAhUAl2BQjxUjC8yykrmCouuEC/BYHPUCgYEA9+GghdabPd7LvKtc
NrhXuXmUr7v6OuqC+VdMCz0HgmdRWVeOutRZT+ZxBxCBgLRJFnEj6EwoFhO3zwky
jMim4TwWeotUfI0o4KOuHiuzpnWRbqN/C/ohNWLx+2J6ASQ7zKTxvqhRkImog9/h
WuWfBpKLZl6Ae1UlZAFMO/7PSSoDgYUAAoGBAP1R1jLPc1kikRwexRvKZhmR01hx
FTCYrRaDX8/g+gmQAWWHf0fOrAi0R7dr6BRlT3unfNMgAi8U2+Iet7vpSz1EgG4Z
XRc4XSK704jhMV0FPF98OFKFDBWlxJsNnt/MwKiwIA9KHbC89OzJGSap02Mqfa0f
8LzMUkP848EZDJkD
-----END ENCRYPTED PRIVATE KEY-----`)
*/

var (
	gPushPunishCmdServer    *PushPunishCmdServer
	oncePushPunishCmdServer sync.Once
)

func GetPushPunishCmdServer() *PushPunishCmdServer {
	oncePushPunishCmdServer.Do(func() {
		gPushPunishCmdServer = &PushPunishCmdServer{}
		gPushPunishCmdServer.initInstance()
	})
	return gPushPunishCmdServer
}

type PushPunishCmdServer struct {
	listenAddr  string
	listener    net.Listener
	mysqlClient *db.MysqlClient

	//publicKeyDsa *dsa.PublicKey
}

func (s *PushPunishCmdServer) initInstance() {
	punishConf := common.GetConfig().Audit["video"].PunishConf
	s.listenAddr = punishConf.ListenAddr

	/*
		block, _ := pem.Decode(key)
		if block == nil {
			glog.Error("expected block to be non-nil", block)
			panic("block is nil")
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			glog.Error("could not unmarshall data: `%s`", err)
			panic(err)
		}
		s.publicKeyDsa = *dsa.PublicKey
	*/
}

func (s *PushPunishCmdServer) Start() {
	http.HandleFunc("/push_video_punish_cmd", s.HandleVideoPunishCmd)
	http.HandleFunc("/push_picture_punish_cmd", s.HandlePicturePunishCmd)
	http.HandleFunc("/push_content_punish_cmd", s.HandleContentPunishCmd)

	srv := &http.Server{Addr: s.listenAddr, Handler: nil}
	var err error
	s.listener, err = net.Listen("tcp", s.listenAddr)
	if err != nil {
		panic(err)
	}
	err = srv.Serve(s.listener)
	glog.Debug("PushPunishCmdServer quit ok. addr:%s", s.listenAddr)
}

func (s *PushPunishCmdServer) Close() {
	s.listener.Close()
}

type MmsReportCmdRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (s *PushPunishCmdServer) HandleVideoPunishCmd(w http.ResponseWriter, req *http.Request) {
	glog.Debug("--------------begin HandleVideoPunishCmd-----------\nreq:[%v]", *req)
	resp := &MmsReportCmdRsp{}
	defer func() {
		respByte, err := json.Marshal(resp)
		if err != nil {
			glog.Error("json marshal failed. err:%v", err)
		}
		/*
			respJson2, err := sj.NewJson(respByte)
			if err != nil {
				glog.Error("go-simplejson NewJson failed. err:%v", err)
			}
		*/
		fmt.Fprint(w, string(respByte))
	}()

	err := req.ParseForm()
	if err != nil {
		resp.Code = -99
		resp.Msg = "解析表单失败"
		return
	}
	glog.Debug("PostForm:[%v]", req.PostForm)

	//check req
	if errNum, errStr := s.CheckVideoPunishCmd(req); errNum <= 0 {
		resp.Code = errNum
		resp.Msg = errStr
		return
	}

	//update video state
	err = s.UpdateVideoState(req)
	if err != nil {
		resp.Code = -99
		resp.Msg = "补刀内部错误"
		return
	}
	resp.Code = 1
	resp.Msg = "处罚成功"
	glog.Debug("--------------end HandleVideoPunishCmd-----------resp:%v", resp)
}

func (s *PushPunishCmdServer) UpdateVideoState(req *http.Request) (err error) {
	//处罚回调存在单个流水多次回调情况，多次回调结果覆盖规则,
	//“违规可以覆盖不违规，不违规不能覆盖违规”，无关结果回调时间顺序。
	//state 1：不违规 2：违规 3：无法处理
	pushStateStr := req.PostFormValue("status")
	pushState, err := strconv.Atoi(pushStateStr)
	if err != nil {
		return
	}
	if pushState == 3 {
		pushState = 0
	}

	extParUrlEncoder := req.PostFormValue("extParUrlEncoder")
	videoIdStr, err := url.QueryUnescape(extParUrlEncoder)
	if err != nil {
		return
	}
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 64)
	if err != nil {
		return
	}
	glog.Debug("pushState:%d, extParUrlEncoder:%v, videoId:%d", pushState, extParUrlEncoder, videoId)

	//先获取待审核视频的审核状态
	tableName, _ := db.GetTableName("video_", videoId)
	querySql := fmt.Sprintf("select state from %s where vid=%d", tableName, videoId)
	var state int
	tempRow, err := db.QueryRow(common.BUDAODB, querySql)
	err = tempRow.Scan(&state)
	if err != nil {
		glog.Error("query video audit state failed. vid:%d, err:%v\n", err)
		return
	}
	glog.Debug("video table. vid:%d, state:%d", videoId, state)
	if state == common.VIDEOSTATE_NOT_PASS_AUDIT ||
		state == common.VIDEOSTATE_DELETED {
		return
	}

	//update
	newState := state
	switch state {
	case common.VIDEOSTATE_WAIT_AUDIT:
		newState = pushState
	case common.VIDEOSTATE_AUDITING:
		newState = pushState
	case common.VIDEOSTATE_PASS_AUDIT:
		if pushState == common.VIDEOSTATE_NOT_PASS_AUDIT {
			newState = pushState
		}
	}
	if newState == state {
		return
	}

	execSql := fmt.Sprintf("update %s set state=%d where vid=%d", tableName, newState, videoId)
	result, err := db.Exec(common.BUDAODB, execSql)
	if err != nil {
		glog.Error("insert failed. execSql:[%s] result:%v err:%v", execSql, result, err)
	}
	return
}

func (s *PushPunishCmdServer) CheckVideoPunishCmd(req *http.Request) (errNum int, errStr string) {
	appKey := req.PostFormValue("appKey")

	/*
		serial := req.PostFormValue("serial")
		cmd := req.PostFormValue("cmd")
		reason := req.PostFormValue("reason")
		msg := req.PostFormValue("msg")
		extParUrlEncoder := req.PostFormValue("extParUrlEncoder")
		sign := req.PostFormValue("sign")
		status := req.PostFormValue("status")
	*/

	secretId := common.GetConfig().Audit["video"].SecretId
	if secretId != appKey {
		glog.Error("bad appKey. appKey:%s, secretId:%s", appKey, secretId)
		return -1, "签名认证不通过"
	}

	/*
		//TODO
		//check sign
		//对appKey，serial，cmd，reason，msg，extParUrlEncoder的内容做DSA签名
		plain := fmt.Sprintf("%s%s%s%s%s%s", appKey, serial, cmd, reason, msg, extParUrlEncoder)

		var h hash.Hash
		h = md5.New()
		r := big.NewInt(0)
		s := big.NewInt(0)

		io.WriteString(h, plain)
		signhash := h.Sum(nil)

		r, s, err := dsa.Sign(rand.Reader, privatekey, signhash)
		if err != nil {
			fmt.Println(err)
		}

		// Verify
		verifystatus := dsa.Verify(s.publicKeyDsa, signhash, r, s)
		glog.Debug("sign verify result:%b", verifystatus) // should be true
	*/

	return 1, ""
}

func (s *PushPunishCmdServer) HandlePicturePunishCmd(w http.ResponseWriter, req *http.Request) {
}

//content punish. only notify illegal result and state is 2
func (s *PushPunishCmdServer) HandleContentPunishCmd(w http.ResponseWriter, req *http.Request) {
	glog.Debug("--------------begin HandleContentPunishCmd-----------\nreq:[%v]", *req)
	resp := &MmsReportCmdRsp{}
	defer func() {
		glog.Debug("resp:%v", resp)
		respByte, err := json.Marshal(resp)
		if err != nil {
			glog.Error("json marshal failed. err:%v", err)
		}
		fmt.Fprint(w, string(respByte))
	}()

	//check req
	err := req.ParseForm()
	if err != nil {
		resp.Code = -99
		resp.Msg = "解析表单失败"
		return
	}
	if len(req.Form) == 0 {
		resp.Code = -99
		resp.Msg = "表单为空"
		return
	}
	glog.Debug("Form:%+v", req.Form)

	extParUrlEncoder := req.Form["extParUrlEncoder"][0]
	jsonStr, err := url.QueryUnescape(extParUrlEncoder)
	if err != nil {
		resp.Code = -1
		resp.Msg = "parse extParUrlEncoder failed."
		return
	}
	auditContentExtPar := &AuditContentExtPar{}
	err = json.Unmarshal([]byte(jsonStr), auditContentExtPar)
	if err != nil {
		glog.Error("json unmarshal failed. err:%v", err)
		return
	}
	contentId := auditContentExtPar.ContentId
	contentType := auditContentExtPar.ContentType

	appId := common.GetConfig().Audit[contentType].Appid
	if len(req.Form["appKey"]) < 1 || req.Form["appKey"][0] != appId {
		glog.Error("appKey:%s, appId:%s", req.Form["appKey"], appId)
		resp.Code = -1
		resp.Msg = "签名认证不通过"
		return
	}

	//update content state
	serial := req.Form["serial"][0]
	cmd := req.Form["cmd"][0]
	glog.Debug("contentId:%s, cmd:%v, serial:%s", contentId, cmd, serial)
	hashId, _ := strconv.ParseUint(contentId, 10, 64)
	var execSql string
	if contentType == "comment" {
		tableName, _ := db.GetTableName("comment_", hashId)
		execSql = fmt.Sprintf("update %s set state=%d where cid=%s", tableName, common.VIDEOSTATE_NOT_PASS_AUDIT, contentId)
	} else if contentType == "title" {
		tableName, _ := db.GetTableName("video_", hashId)
		execSql = fmt.Sprintf("update %s set state=%d where vid=%s", tableName, common.VIDEOSTATE_NOT_PASS_AUDIT, contentId)
	}
	_, err = db.Exec(common.BUDAODB, execSql)
	if err != nil {
		resp.Code = -99
		resp.Msg = "补刀内部错误"
		return
	}
	resp.Code = 1
	resp.Msg = "处罚成功"
	glog.Debug("--------------end HandleContentPunishCmd-----------resp:%v", resp)
}

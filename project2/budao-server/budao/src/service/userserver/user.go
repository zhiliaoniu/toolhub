package userserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sumaig/glog"

	"common"
	"db"
	"service/transfer"
	"service/util"
	pb "twirprpc"
)

// GetServer return server of user service
func GetServer() *Server {
	server := &Server{}
	server.initServer()
	return server
}

// WeixinVerify define receive weixin struct
type WeixinVerify struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// QQVerify define receive qq struct
type QQVerify struct {
	Ret                int    `json:"ret"`
	Msg                string `json:"msg"`
	Nickname           string `json:"nickname"`
	Figureurl          string `json:"figureurl"`
	Figureurl_1        string `json:"figureurl_1"`
	Figureurl_2        string `json:"figureurl_2"`
	Figureurl_qq_1     string `json:"figureurl_qq_1"`
	Figureurl_qq_2     string `json:"figureurl_qq_2"`
	Is_yellow_vip      string `json:"is_yellow_vip"`
	Is_yellow_year_vip string `json:"is_yellow_year_vip"`
	Yellow_vip_level   string `json:"yellow_vip_level"`
}

const (
	APPID = "1106124920"
)

// Server identify for user RPC
type Server struct {
	userTokenExpireSec int
}

func (s *Server) initServer() {
	s.userTokenExpireSec = common.USER_TOKEN_EXPIRE_SEC
}

// Login function handle request for user.
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (resp *pb.LoginResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.LoginResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	loginType := int32(req.GetLoginType())
	openid := req.GetOpenid()
	token := req.GetToken()
	userName := req.GetUserName()
	userPhoto := req.GetUserPhoto()

	//2.verify
	var weiResult *WeixinVerify
	var qqResult *QQVerify
	if loginType == int32(pb.LoginType_WEIXIN) {
		// weixin verification
		url := fmt.Sprintf("https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s", token, openid)
		weiResult, err = getWeinxinVerifyData(url)
		if err != nil || weiResult.Errcode != 0 || weiResult.Errmsg != "ok" {
			glog.Error("weixin verify failed. req:%v, err:%v, weiResult:%v", req, err, weiResult)
			resp.Status.Code = pb.Status_SERVER_ERR
			err = errors.New("weixin verify failed")
			return
		}
		glog.Debug("weixin verify success.token:%s, openid:%s", token, openid)
	} else if loginType == int32(pb.LoginType_QQ) {
		// qq verify
		url := fmt.Sprintf("https://graph.qq.com/user/get_simple_userinfo?access_token=%s&oauth_consumer_key=%s&openid=%s&format=json", token, APPID, openid)
		qqResult, err = getQQVerifyData(url)
		if err != nil || qqResult.Ret != 0 || qqResult.Msg != "" {
			glog.Error("qq verify failed. req:%v, err:%v, qqResult:%v", req, err, qqResult)
			resp.Status.Code = pb.Status_SERVER_ERR
			err = errors.New("qq verify failed")
			return
		}
		glog.Debug("qq verify success.token:%s, openid:%s", token, openid)
	}

	openidString, err := db.MysqlEscapeString(openid)
	userNameString, err := db.MysqlEscapeString(userName)
	userPhotoString, err := db.MysqlEscapeString(userPhoto)

	//3.check whether user has login before
	tableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_"]
	tableName := "user_0"
	var i uint64
	for i = 1; i < tableNum; i++ {
		tempStr := fmt.Sprintf(", user_%d", i)
		tableName = tableName + tempStr
	}

	var uid, name, photo, userToken string
	execSQL := fmt.Sprintf("select uid, name, photo from %s where openid='%s'", tableName, openidString)
	rows, err := db.Query(common.BUDAODB, execSQL)
	if err != nil {
		glog.Error("select user info failed. execSQL:%s, err:%v", execSQL, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&uid, &name, &photo); err != nil {
			glog.Error("scan user info failed. err:%v", err)
			return
		}
	}
	if uid != "" {
		//get token
		if userToken, err = getUserToken(uid); err != nil {
			glog.Error("get user token failed. req:%v, err:%v", req, err)
			return
		}

		user := &pb.UserInfo{
			UserId:    uid,
			UserName:  name,
			UserPhoto: photo,
			Token:     userToken,
		}
		glog.Debug("user:%v has login before, not need insert", user)

		resp.Status.Code = pb.Status_OK
		resp.UserInfo = user

		return
	}

	glog.Debug("begin insert new user")
	//4.insert new user
	idGenerate := transfer.GetIdGenerator()
	userID, autoIncreID, err := idGenerate.GetItemId("user_")
	if err != nil {
		glog.Error("generate userid failed. req, err:%v", req, err)
		return
	}
	uid = strconv.FormatUint(userID, 10)
	if userToken, err = getUserToken(uid); err != nil {
		glog.Error("get user token failed. userID:%d, err:%v", userID, err)
		return
	}

	tableName, _ = db.GetTableName("user_", userID)
	sqlString := fmt.Sprintf("update %s set uid=%d, name='%s', photo='%s', login_method=%d, openid='%s', token='%s' where id=%d", tableName, userID, userNameString, userPhotoString, loginType, openidString, userToken, autoIncreID)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("insert user failed. sqlString%s, err:%v", sqlString, err)
		return
	}

	userInfo := &pb.UserInfo{
		UserId:    strconv.FormatUint(userID, 10),
		UserName:  userName,
		UserPhoto: userPhoto,
		Token:     userToken,
	}
	glog.Debug("insert new user success. user:%v", userInfo)

	resp.Status.Code = pb.Status_OK
	resp.UserInfo = userInfo

	return
}

func getUserToken(userId string) (token string, err error) {
	//1.generate token
	token = util.GenerateToken(userId)

	//2.check token is expired
	expire, err := util.CheckTokenExpire(userId, token)
	if err != nil || expire == false {
		return
	}

	//3.extend token expire time
	err = util.ExtendTokenExpire(userId, token)
	return
}

// getWeinxinVerifyData function get weixin verify result
func getWeinxinVerifyData(curl string) (verifyResult *WeixinVerify, err error) {
	resp, err := http.Get(curl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else {
		glog.Debug("read weixin result succeed")
	}

	result := WeixinVerify{}
	if err := json.Unmarshal(body, &result); err != nil {
		glog.Debug("json data transfer failed...")
		return nil, err
	}

	return &WeixinVerify{
		Errcode: result.Errcode,
		Errmsg:  result.Errmsg,
	}, nil
}

func getQQVerifyData(curl string) (verifyResult *QQVerify, err error) {
	resp, err := http.Get(curl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	} else {
		glog.Debug("read qq result succeed")
	}

	result := QQVerify{}
	if err := json.Unmarshal(body, &result); err != nil {
		glog.Debug("json data transfer failed...")
		return nil, err
	}

	return &QQVerify{
		Ret: result.Ret,
		Msg: result.Msg,
	}, nil
}

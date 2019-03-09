package wechatservice

import (
	"common"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sumaig/glog"
	"service"
	"service/api"
	"strconv"
	"time"
)

const (
	CORPID      = "wwd84057ff109fac91"
	SECRET      = "rBdL3IEp-VBhn1fMm-R2ubAnFMaSkffpe7sTMn7TATw"
	VERFITY_STR = "budao.ops"
	AES_KEY     = "sfe023f_9fd&fwfl"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

//微信回调
func (s *Server) WechatRedirect(ctx context.Context, param *api.WechatRedirectRequest) (resp *api.WechatRedirectResponse, err error) {
	resp = &api.WechatRedirectResponse{
		Data: &api.BackStruct{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	var userId string

	if len(param.Code) > 0 {
		//如果存在，授权成功
		//获取用户的userId
		//生成token  md5("") +　userId + time
		accessToken := AccessToken()

		//获取userId
		userInfo := getWechatUserId(accessToken, param.Code)
		glog.Info(userInfo)
		if _, ok := userInfo["UserId"]; ok {
			userId = userInfo["UserId"].(string)
		} else {
			RefreshToken(CORPID, SECRET)
			err = fmt.Errorf(`Get userid is err %v`, userInfo)
			return
		}

		//生成token
		token := getToken(userId)
		glog.Info(token)

		resp.Code = service.SUCCESS_CODE
		resp.Msg = "OK"
		resp.Data.UserId = userId
		resp.Data.Token = token
	} else {
		err = fmt.Errorf("Get code is empty")
	}
	return
}

//获取token
func getToken(userId string) string {
	//生成token  md5("") + userId + time
	ver := common.GetMd5Str([]byte(VERFITY_STR))
	//时间
	t := time.Now().Unix() + (service.LOGIN_EXPIRE * 60 * 60)
	//加密的串
	origData := ver + "_" + userId + "_" + strconv.FormatInt(t, 10)
	result, err := common.AesEncrypt([]byte(origData), []byte(AES_KEY))
	if err != nil {
		panic(err)
	}

	//利用base64转成字符串
	token := base64.StdEncoding.EncodeToString(result)

	return token
}

//获取用户userId
func getWechatUserId(accessToken, code string) map[string]interface{} {
	url := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=" + accessToken + "&code=" + code
	_, buff := common.HttpGet(url)

	m := make(map[string]interface{})
	err := json.Unmarshal(buff, &m)
	if err != nil {
		glog.Error(err)
		return nil
	}

	return m
}

//获取token
func getWechatAccessToken(corpid, corpsecret string) map[string]interface{} {
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + corpsecret
	_, buff := common.HttpGet(url)

	m := make(map[string]interface{})
	err := json.Unmarshal(buff, &m)
	if err != nil {
		glog.Error(err)
		return nil
	}

	return m
}

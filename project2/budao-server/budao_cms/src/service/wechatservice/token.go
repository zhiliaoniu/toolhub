package wechatservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var tokenMutex *sync.RWMutex

type GetToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json: expires_in"`
}

var gtNow GetToken

func RefreshToken(corpId string, corpSecret string) bool {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpId, corpSecret)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	var gt GetToken
	err = json.Unmarshal(body, &gt)
	if err != nil {
		return false
	}

	tokenMutex.Lock()
	gtNow = gt
	tokenMutex.Unlock()

	return true
}

func AccessToken() string {

	tokenMutex.RLock()
	at := gtNow.AccessToken
	tokenMutex.RUnlock()
	return at

}

func InitToken() {
	tokenMutex = new(sync.RWMutex)
}

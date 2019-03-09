package main

import (
	"encoding/json"
	"fmt"
	//	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var tokenMutex *sync.RWMutex

// TODO token class ?
type GetToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json: expires_in"`
}

type TextMsg struct {
	Content string `json:"content"`
}

type SendMsg struct {
	ToUser  string  `json:"touser"`
	ToParty string  `json:"toparty"`
	ToTag   string  `json:"totag"`
	MsgType string  `json:"msgtype"`
	AgentId int     `json:"agentid"`
	Text    TextMsg `json:"text"`
	Safe    int     `json:"safe"`
}

type SendMsgRet struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Invaliduser  string `json:"Invaliduser"`
	invalidparty string `json:"invalidparty"`
	invalidtag   string `json:"invalidtag"`
}

var gtNow GetToken

func refreshToken(corpId string, corpSecret string) bool {
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

// TODO timer for token refresh
func tokenTimer() {
	ticker := time.NewTicker((time.Duration)(gtNow.ExpiresIn-60) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// refreshToken()
			ticker = time.NewTicker((time.Duration)(gtNow.ExpiresIn-60) * time.Second)
		}
	}
}

func accessToken() string {

	tokenMutex.RLock()
	at := gtNow.AccessToken
	tokenMutex.RUnlock()
	return at

}

func initToken() {
	tokenMutex = new(sync.RWMutex)
}

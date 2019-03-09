package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
)

func makeSendMsg(agentId int, toUser string, msg string) []byte {
	s := SendMsg{
		ToUser:  toUser,
		ToParty: ``,
		ToTag:   ``,
		MsgType: `text`,
		AgentId: agentId,
		Text: TextMsg{
			Content: msg,
		},
		Safe: 0,
	}

	b, _ := json.Marshal(s)
	return b
}

func frequencyCheck(agentId int, toUser string) bool {
	//TODO
	return true
}

func send(agentId int, toUser string, msg string) bool {

	if !frequencyCheck(agentId, toUser) {
		return false
	}

	// TODO async???
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", accessToken())

	body := bytes.NewBuffer(makeSendMsg(agentId, toUser, msg))
	res, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		return false
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return false
	}
	var ret SendMsgRet
	err = json.Unmarshal(result, &ret)
	if err != nil {
		return false
	}

	fmt.Printf("[%d, %s]:%s, ret = %v\n", agentId, toUser, msg, ret)
	if ret.Errcode != 0 {
		fmt.Printf("[%d, %s] error: ret code %d\n", agentId, toUser, ret.Errcode)
		return false
	}

	return true
}

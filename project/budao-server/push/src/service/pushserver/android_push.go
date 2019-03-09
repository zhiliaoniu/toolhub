package pushserver

import (
	"common"

	"github.com/sumaig/glog"
	"github.com/ylywyn/jpush-api-go-client"
)

//const (
//	//appKey = "7f8b616f3f78a99bd0a23fe5"
//	//secret = "06f0c168bd87b54f4fdf6654"
//
//	appKey = "ec22af927fc53702169818c0"
//	secret = "e08ce5d0bf5a145c21c2d2dd"
//)

func PushAndroidMessage(pushMsg *PushRecord) (err error) {
	appKey := common.GetConfig().PushConf.AndroidAppKey
	secret := common.GetConfig().PushConf.AndroidSecret
	glog.Debug("begin send android msg. pushMsg:%+v", pushMsg)

	//1.init Platform
	var pf jpushclient.Platform
	pf.Add(jpushclient.ANDROID)

	//2.init Audience
	var ad jpushclient.Audience
	//s := []string{"1", "2", "3"}
	//ad.SetTag(s)
	//ad.SetAlias(s)
	//ad.SetID(s)
	ad.All()

	//3.init Notice
	aps := &APS{
		Alert: &Alert{
			Title: pushMsg.PushTitle,
			Body:  pushMsg.PushContent,
		},
		Badge: 1,
		Sound: "default",
	}

	var msg jpushclient.Message
	msg.Title = pushMsg.PushTitle
	msg.Content = pushMsg.PushContent
	msg.ContentType = "text"
	msg.AddExtras("route_url", pushMsg.PushUrl)
	msg.AddExtras("aps", aps)

	payload := jpushclient.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	payload.SetMessage(&msg)

	bytes, _ := payload.ToBytes()
	glog.Debug("payload:%s", string(bytes))

	//4.push
	c := jpushclient.NewPushClient(secret, appKey)

	str, err := c.Send(bytes)
	if err != nil {
		glog.Error("err:%s", err.Error())
	} else {
		glog.Debug("ok:%s", str)
	}
	glog.Debug("end android push.")

	return
}

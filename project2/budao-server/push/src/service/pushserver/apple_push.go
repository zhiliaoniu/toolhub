package pushserver

import (
	"database/sql"
	"db"
	"encoding/json"
	"fmt"
	"log"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sumaig/glog"

	"common"
)

var GIdentifier uint32 = 0

/**
 * @Breif 将应用程序所推送的消息发送给APNS服务器
 *
 * @Returns 成功返回0，失败返回-1
 */
/*
func PushIOSMessage(pushRecord *PushRecord) (err error) {
	glog.Debug("begin send")
	certFileName := common.GetConfig().PushConf.CertFileName
	keyFileName := common.GetConfig().PushConf.KeyFileName
	apn, err := apns.New(certFileName, keyFileName, "gateway.sandbox.push.apple.com:2195", 1*time.Second)
	if err != nil {
		glog.Debug("connect error: %s\n", err.Error())
		os.Exit(1)
	}
	glog.Debug("connect successed!")
	readError(apn.ErrorChan)
	token := "0817182b73070df36e5c94b66f803513a01aa98e68903ad1c00e87595ae3e15c"

	payload := apns.Payload{}
	payload.Aps.Alert.Body = "hello @zhongfeng. this is from @yangshengzhi"

	notification := apns.Notification{}
	notification.DeviceToken = token
	notification.Identifier = atomic.LoadUint32(&GIdentifier)
	atomic.AddUint32(&GIdentifier, 1)
	notification.Payload = &payload
	err = apn.Send(&notification)
	if err != nil {
		glog.Error("send id(%v). err:%v", notification.Identifier, err)
		return
	}
	glog.Debug("send id(%v)", notification.Identifier)
	return
}
*/

func readError(errorChan <-chan error) {
	//for {
	//	apnerror := <-errorChan
	//	glog.Error("apnerr:%v", apnerror.Error())
	//}
}

/**
 * @Breif 将应用程序所推送的消息发送给APNS服务器
 */
func PushIOSMessageAll(pushRecord *PushRecord) (err error) {
	glog.Debug("begin ios push all")
	tokenMap := make(map[string]bool, 0)
	//1.read all token
	tableNum := int(common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_token_"])
	for i := 0; i < tableNum; i++ {
		tableName := fmt.Sprintf("%s%d", "user_token_", i)
		sqlString := fmt.Sprintf("select device_token from %s", tableName)
		var rows *sql.Rows
		rows, err = db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("read mysql failed. querySql:%s, err:%v", sqlString, err)
			return
		}
		var deviceToken string
		for rows.Next() {
			err = rows.Scan(&deviceToken)
			if err != nil {
				glog.Error("scan failed. err:%v", err)
				return
			}
			tokenMap[deviceToken] = true
		}
	}

	//test code, delete later
	tokenMap = map[string]bool{
		"9a694b78cdcdc3d654739bd092190bdca0920bf1ce02160163745ae57bb8e7b5": true,
	}

	glog.Debug("tokenMap.len:%d tokenMap:%v", len(tokenMap), tokenMap)

	cert, err := certificate.FromPemFile("./conf/sandbox.pem", "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	client := apns2.NewClient(cert).Development()
	//2.send one by one
	for token, _ := range tokenMap {
		notification := &apns2.Notification{}
		notification.DeviceToken = token
		notification.Topic = "com.yy.lasthit"
		aps := &Message{
			APS: &APS{
				Alert: &Alert{
					Title: pushRecord.PushTitle,
					Body:  pushRecord.PushContent,
				},
				Badge: 1,
				Sound: "default",
			},
			RouteUrl: pushRecord.PushUrl,
		}
		var payload []byte
		payload, err = json.Marshal(aps)
		if err != nil {
			glog.Error("marshal failed. err:%v", err)
			return
		}
		glog.Debug("payload:%v", string(payload))
		notification.Payload = payload

		res, err := client.Push(notification)

		if err != nil {
			glog.Error("err:%v", err)
		}

		glog.Debug("push result. code:%v id:%v reason:%v\n", res.StatusCode, res.ApnsID, res.Reason)
	}
	glog.Debug("end ios push all")
	return
}

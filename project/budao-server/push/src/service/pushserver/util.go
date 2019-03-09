package pushserver

import (
	"common"
	"db"
	"fmt"
	"time"

	"github.com/sumaig/glog"
)

//cron使用文档：https://godoc.org/github.com/robfig/cron
/*
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/

type Message struct {
	APS      *APS   `json:"aps"`
	RouteUrl string `json:"route_url"`
}

type Route struct {
	RouteUrl string `json:"route_url"`
}

type APS struct {
	Alert *Alert `json:"alert"`
	Badge int    `json:"badge"`
	Sound string `json:"sound"`
}

type Alert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PushRecord struct {
	Id            uint64    //
	OpUid         uint64    //运营uid
	PushObj       int       //'推送对象：0-全量人群  '
	PushType      int       //'推送方式: 0-立即推送  1-定时推送'
	PushStatus    int       // 推送状态: 0-未发送  1-发送中  2-发送成功  3-发送失败
	PushTitle     string    //push标题
	PushContent   string    //push内容
	PushDescribe  string    //描述
	PushBeginTime time.Time //推送时间
	PushEndTime   time.Time //推送中止时间，即推送消息的存活时间
	PushVideoType int       //推送视频类型: 0-话题  1-视频
	PushUrl       string    //push的超链接地址
	Device        string    //设备:安卓-android  苹果-ios 全部-all
}

/**
 * @Breif 扫描推送任务
 */
func PushTaskScan() {
	glog.Debug("begin PushTaskScan")
	//1.prepare task
	pushRecordMap, err := GetWaitPushTask()
	if err != nil {
		return
	}
	if len(pushRecordMap) == 0 {
		glog.Debug("no push task")
		return
	}
	glog.Debug("push task pushRecordMap. len:%d", len(pushRecordMap))

	//2.get push task
	for id, pushRecord := range pushRecordMap {
		glog.Debug("begin push id:%d, record:%+v", id, pushRecord)
		pushStatus := 2
		if err := PushMessageAll(pushRecord); err != nil {
			pushStatus = 3
			glog.Error("push message all failed. err:%v", err)
		}
		//3.update push task status
		glog.Debug("id:%d, push_status:%d", id, pushStatus)
		//UpdatePushRecordStatus(pushRecord, pushStatus)
	}

	glog.Debug("end PushTaskScan")
}

func UpdatePushRecordStatus(pushRecord *PushRecord, pushStatus int) {
	execSql := fmt.Sprintf("update push set push_status=%d where id=%d", pushStatus, pushRecord.Id)
	_, err := db.Exec(common.BUDAODB, execSql)
	if err != nil {
		glog.Error("update mysql failed. execSql:%s, err:%v", execSql, err)
		return
	}
}

func PushMessageAll(pushRecord *PushRecord) (err error) {
	//1.push andrioid
	if pushRecord.Device == "all" ||
		pushRecord.Device == "android" {
		err = PushAndroidMessage(pushRecord)
		if err != nil {
			return
		}
	}

	//2.push ios
	if pushRecord.Device == "all" ||
		pushRecord.Device == "ios" {
		err = PushIOSMessageAll(pushRecord)
		if err != nil {
			return
		}
	}
	return
}

func GetWaitPushTask() (pushRecordMap map[uint64]*PushRecord, err error) {
	pushRecordMap = make(map[uint64]*PushRecord, 0)
	querySql := fmt.Sprintf("select id, op_uid, push_obj, push_type, push_title, push_content, push_describe, push_time, push_end_time, push_video_type, push_url, device from push where (push_status = 0 or push_status = 3) and status = 0")
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
		return
	}
	now := time.Now()
	for rows.Next() {
		pushRecord := &PushRecord{}
		err = rows.Scan(&pushRecord.Id, &pushRecord.OpUid, &pushRecord.PushObj, &pushRecord.PushType, &pushRecord.PushTitle, &pushRecord.PushContent, &pushRecord.PushDescribe, &pushRecord.PushBeginTime, &pushRecord.PushEndTime, &pushRecord.PushVideoType, &pushRecord.PushUrl, &pushRecord.Device)
		if err != nil {
			glog.Error("scan failed. err:%v", err)
			return
		}
		if pushRecord.PushType == 1 &&
			(now.Before(pushRecord.PushBeginTime) ||
				now.After(pushRecord.PushEndTime)) {
			//glog.Debug("task is stale. now:%v, pushRecord:%+v ", now, pushRecord)
			continue
		}
		pushRecordMap[pushRecord.Id] = pushRecord
	}
	return
}

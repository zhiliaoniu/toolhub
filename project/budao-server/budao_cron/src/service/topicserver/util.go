package topicserver

import (
	"common"
	"db"
	"fmt"
	"time"

	"github.com/sumaig/glog"
)

func InitTopicDynamicInfo() {
	//1.read mysql
	sqlString := fmt.Sprintf("select topic_id, user_num from topic")
	rows, err := db.Query(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("query topic table failed. sqlString:%s, err:%v", sqlString, err)
		return
	}
	defer rows.Close()

	fields := make([]interface{}, 0)
	var topicId, userNum string
	for rows.Next() {
		err = rows.Scan(&topicId, &userNum)
		if err != nil {
			glog.Error("scan topic dynamic info failed. err:%v", err)
			return
		}
		fields = append(fields, topicId, userNum)
	}

	if _, err := db.HMSet(common.TOPIC_DYNAMIC, fields); err != nil {
		glog.Error("hmset failed. err:%v", err)
		return
	}
	glog.Debug("update topic dynamic info success")
}

//InitUserSubscribedTopic update user subscribe topic list
func InitUserSubscribedTopic() {
	userTopicMap := make(map[uint64][]interface{}) //<uid <topicid, update_time>>
	//1.get user subscribe topic from mysql
	tableNUM := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_follow_topic_"]
	var i uint64
	for i = 0; i < tableNUM; i++ {
		tableName := fmt.Sprintf("user_follow_topic_%d", i)
		sqlString := fmt.Sprintf("select uid, topic_id, update_time from %s", tableName)
		rows, err := db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("query %s failed. sqlString:%s, err:%v", tableName, sqlString, err)
			continue
		}
		defer rows.Close()

		var (
			uid        uint64
			topicId    uint64
			updateTime time.Time
		)
		for rows.Next() {
			if err := rows.Scan(&uid, &topicId, &updateTime); err != nil {
				rows.Close()
				break
			}
			userTopicMap[uid] = append(userTopicMap[uid], updateTime.Unix(), topicId)
		}
		//close mysql conn active
		rows.Close()
	}

	//2.set user subscribe topic to redis by zset
	for uid, topicIds := range userTopicMap {
		key := fmt.Sprintf("%s%d", common.USER_SUBSCRIBE_TOPIC_WITH_TIME, uid)
		_, err := db.ZAddMulti(key, topicIds)
		if err != nil {
			glog.Error("zadd multi failed. err:%v", err)
			continue
		}
	}
	glog.Debug("UpdateUserTopic success")
}

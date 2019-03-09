package timelineserver

import (
	"fmt"
	"strconv"

	"github.com/sumaig/glog"

	"common"
	"db"
	"service/util"
)

//call once on server start.read video dynamic info to redis
func InitVideoDynamicInfo() {
	glog.Debug("start init video dynamic info")
	tableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc[common.VIDEO_TABLE_PREFIX]
	var num uint64
	for num = 0; num < tableNum; num++ {
		tableName := fmt.Sprintf("%s%d", common.VIDEO_TABLE_PREFIX, num)
		count := 0
		for {
			sqlString := fmt.Sprintf("select vid, favor_num, fake_favor_num, comment_num, fake_comment_num, view_num, fake_view_num, share_num, fake_share_num, could_favor, could_comment, could_share from %s where state = 2 limit %d,%d", tableName, count, common.READ_MYSQL_MAX_ROWS)
			rows, err := db.Query(common.BUDAODB, sqlString)
			if err != nil {
				glog.Error("query %s table failed. err:%v", tableName, err)
				break
			}
			defer rows.Close()

			for rows.Next() {
				var vid string
				dynamicInfo := &util.VideoDynamicInfo{}
				var CouldLike, CouldComment, CouldShare int
				var FakeViewNum, FakeFavorNum, FakeCommentNum, FakeShareNum uint32
				err = rows.Scan(&vid, &dynamicInfo.LikeCount, &FakeFavorNum, &dynamicInfo.CommentCount, &FakeCommentNum, &dynamicInfo.ViewCount, &FakeViewNum, &dynamicInfo.ShareCount, &FakeShareNum, &CouldLike, &CouldComment, &CouldShare)
				if err != nil {
					glog.Error("scan video dynamic info failed. err:%v", err)
					continue
				}
				//TODO improve
				if CouldLike == 0 {
					dynamicInfo.LikeDisabled = true
				}
				if CouldComment == 0 {
					dynamicInfo.CommentDisabled = true
				}
				if CouldShare == 0 {
					dynamicInfo.ShareDisabled = true
				}
				dynamicInfo.LikeCount += FakeFavorNum
				dynamicInfo.CommentCount += FakeCommentNum
				dynamicInfo.ViewCount += FakeViewNum
				dynamicInfo.ShareCount += FakeShareNum

				fields := make([]interface{}, 0)
				fields = append(fields, common.LIKE_COUNT+"_"+vid, dynamicInfo.LikeCount, common.COMMENT_COUNT+"_"+vid, dynamicInfo.CommentCount, common.VIEW_COUNT+"_"+vid, dynamicInfo.ViewCount, common.SHARE_COUNT+"_"+vid, dynamicInfo.ShareCount, common.COMMENT_DISABLED+"_"+vid, dynamicInfo.CommentDisabled, common.LIKE_DISABLED+"_"+vid, dynamicInfo.LikeDisabled, common.SHARE_DISABLED+"_"+vid, dynamicInfo.ShareDisabled)
				_, err = db.HMSet(common.VIDEO_DYNAMIC, fields)
				if err != nil {
					glog.Error("hmset video dynamic info failed. key:%s, fields:%v, err:%v", common.VIDEO_DYNAMIC, fields, err)
					continue
				}
			}
			//close active. 防止for循环很多，defer只在函数退出才执行，占用过多连接
			rows.Close()
		}
	}
	glog.Debug("end init video dynamic info")
}

//InitUserAct only call once
func InitUserAct() {
	go initUserActVideo()
	go initUserActComment()
	go initUserActTopic()
	go initUserActQuestion()
}

func initUserActVideo() {
	glog.Debug("start to initUserActVideo")
	userVideoTableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_favor_video_"]
	for i := 0; i < int(userVideoTableNum); i++ {
		userVideoTableName := fmt.Sprintf("user_favor_video_%d", i)
		sqlString := fmt.Sprintf("select uid, vid from %s", userVideoTableName)
		rows, err := db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("query %s failed. err:%v", userVideoTableName, err)
			continue
		}
		defer rows.Close()

		uservidMap := make(map[string][]string, 0)
		for rows.Next() {
			var (
				uid uint64
				vid uint64
			)
			rows.Scan(&uid, &vid)
			uservidMap[strconv.FormatUint(uid, 10)] = append(uservidMap[strconv.FormatUint(uid, 10)], strconv.FormatUint(vid, 10))
		}

		for uid, vidArr := range uservidMap {
			var uservidArr []string
			for _, vid := range vidArr {
				vidKey := fmt.Sprintf("%s_vfavor", vid)
				uservidArr = append(uservidArr, vidKey, "1")
			}
			uservidInter := common.TransStrArrToInterface(uservidArr)
			key := fmt.Sprintf("user_act_%s", uid)
			ret, err := db.HMSet(key, uservidInter)
			if ret != "ok" && err != nil {
				glog.Error("insert UserActQuestion info into redis failed. err:%v", err)
			}
		}
	}
	glog.Debug("end to initUserActVideo")
}

func initUserActComment() {
	glog.Debug("start to initUserActComment")
	userCommentTableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_favor_comment_"]
	for i := 0; i < int(userCommentTableNum); i++ {
		userCommentTableName := fmt.Sprintf("user_favor_comment_%d", i)
		sqlString := fmt.Sprintf("select uid, cid from %s", userCommentTableName)
		rows, err := db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("query %s failed. err:%v", userCommentTableName, err)
			continue
		}
		defer rows.Close()

		usercidMap := make(map[string][]string, 0)
		for rows.Next() {
			var (
				uid uint64
				cid uint64
			)
			rows.Scan(&uid, &cid)
			usercidMap[strconv.FormatUint(uid, 10)] = append(usercidMap[strconv.FormatUint(uid, 10)], strconv.FormatUint(cid, 10))
		}

		for uid, cidArr := range usercidMap {
			var usercidArr []string
			for _, cid := range cidArr {
				cidKey := fmt.Sprintf("%s_cfavor", cid)
				usercidArr = append(usercidArr, cidKey, "1")
			}
			usercidInter := common.TransStrArrToInterface(usercidArr)
			key := fmt.Sprintf("user_act_%s", uid)
			ret, err := db.HMSet(key, usercidInter)
			if ret != "ok" && err != nil {
				glog.Error("insert UserActQuestion info into redis failed. err:%v", err)
			}
		}
	}
	glog.Debug("end to initUserActComment")
}

func initUserActTopic() {
	glog.Debug("start to initUserActTopic")
	userTopicTableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_follow_topic_"]
	for i := 0; i < int(userTopicTableNum); i++ {
		userTopicTableName := fmt.Sprintf("user_follow_topic_%d", i)
		sqlString := fmt.Sprintf("select uid, topic_id from %s", userTopicTableName)
		rows, err := db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("query %s failed. err:%v", userTopicTableName, err)
			continue
		}
		defer rows.Close()

		usertopicidMap := make(map[string][]string, 0)
		for rows.Next() {
			var (
				uid     uint64
				topicid uint64
			)
			rows.Scan(&uid, &topicid)
			usertopicidMap[strconv.FormatUint(uid, 10)] = append(usertopicidMap[strconv.FormatUint(uid, 10)], strconv.FormatUint(topicid, 10))
		}

		for uid, topiccidArr := range usertopicidMap {
			var userTopicidArr []string
			for _, topicid := range topiccidArr {
				topicidKey := fmt.Sprintf("%s_subscribe", topicid)
				userTopicidArr = append(userTopicidArr, topicidKey, "1")
			}
			userTopicidInter := common.TransStrArrToInterface(userTopicidArr)
			key := fmt.Sprintf("user_act_%s", uid)
			ret, err := db.HMSet(key, userTopicidInter)
			if ret != "ok" && err != nil {
				glog.Error("insert UserActQuestion info into redis failed. err:%v", err)
			}
		}
	}
	glog.Debug("end to initUserActTopic")
}

type questionOptionID struct {
	QuestionID string
	OptionID   string
}

func initUserActQuestion() {
	glog.Debug("start to initUserActQuestion")
	userQuesTableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_question_"]
	for i := 0; i < int(userQuesTableNum); i++ {
		userQuesTableName := fmt.Sprintf("user_question_%d", i)
		sqlString := fmt.Sprintf("select uid, question_id, option_id from %s", userQuesTableName)
		rows, err := db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("query %s failed. err:%v", userQuesTableName, err)
			continue
		}
		defer rows.Close()

		userQuesMap := make(map[string][]*questionOptionID, 0)
		for rows.Next() {
			var (
				uid        uint64
				questionid uint64
				optionid   uint64
			)
			rows.Scan(&uid, &questionid, &optionid)
			questionOption := &questionOptionID{
				QuestionID: strconv.FormatUint(questionid, 10),
				OptionID:   strconv.FormatUint(optionid, 10),
			}
			userQuesMap[strconv.FormatUint(uid, 10)] = append(userQuesMap[strconv.FormatUint(uid, 10)], questionOption)
		}

		for uid, questionOptionArr := range userQuesMap {
			var questionidArr []string
			for _, questionOption := range questionOptionArr {
				questionidArr = append(questionidArr, questionOption.QuestionID, questionOption.OptionID)
			}
			questionidInter := common.TransStrArrToInterface(questionidArr)
			key := fmt.Sprintf("user_act_%s", uid)
			ret, err := db.HMSet(key, questionidInter)
			if ret != "ok" && err != nil {
				glog.Error("insert UserActQuestion info into redis failed. err:%v", err)
			}
		}
	}
	glog.Debug("end to initUserActQuestion")
}

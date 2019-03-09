package topicserver

import (
	"common"
	"db"
	"fmt"
	"strconv"
	"time"
	pb "twirprpc"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
)

func (s *Server) cronTask() {
	timer := time.NewTimer(time.Second * time.Duration(1))

	for {
		if s.cronTaskDisable {
			glog.Debug("exit config server crontab")
			break
		}
		select {
		case <-timer.C:
			timer.Reset(time.Minute * time.Duration(30))
			s.UpdateTopicWithScore()
			s.UpdateTopicVideo()
			s.UpdateTopicWithID()
		}
	}
}

// UpdateTopicWithScore load topic data to redis regularly
func (s *Server) UpdateTopicWithScore() {
	sqlString := fmt.Sprintf("select topic_id, disable, weight from topic")
	rows, err := db.Query(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("UpdateTopicWithScore query topic err:%v", err)
		return
	}
	defer rows.Close()

	loop := 0
	var topicPackage []interface{}
	for rows.Next() {
		topicInfo := TopicInfo{}
		rows.Scan(&topicInfo.TopicID, &topicInfo.Disable, &topicInfo.Weight)
		if topicInfo.Disable == 1 {
			exist, _ := db.ZExists(common.TOPICSORTSET, strconv.FormatUint(topicInfo.TopicID, 10))
			if exist == true {
				db.ZDelete(common.TOPICSORTSET, strconv.FormatUint(topicInfo.TopicID, 10))
			}
			continue
		}

		topicPackage = append(topicPackage, topicInfo.Weight, strconv.FormatUint(topicInfo.TopicID, 10))
		loop++
		if loop == 200 {
			loop = 0
			_, err := db.ZAddMulti(common.TOPICSORTSET, topicPackage)
			if err != nil {
				glog.Debug("Misaligned field in hash table")
			}
			topicPackage = make([]interface{}, 0)
		}
	}

	if len(topicPackage) != 0 {
		_, err = db.ZAddMulti(common.TOPICSORTSET, topicPackage)
		if err != nil {
			glog.Debug("Misaligned field in hash table")
		}
	}
	glog.Debug("UpdateTopicWithScore success")

	return
}

// UpdateTopicWithID load topic data to redis regularly
func (s *Server) UpdateTopicWithID() {
	rows, err := db.Query(common.BUDAODB, "select id, topic_id, name, pic, disable, weight, description, user_num, fake_user_num, video_num from topic")
	if err != nil {
		glog.Error("UpdateTopicWithID query topic failed. err:%v", err)
		return
	}
	defer rows.Close()

	loop := 0
	var topicPackage []interface{}
	for rows.Next() {
		topicInfo := TopicInfo{}
		rows.Scan(&topicInfo.AutoID, &topicInfo.TopicID, &topicInfo.Name, &topicInfo.Picture, &topicInfo.Disable, &topicInfo.Weight, &topicInfo.Desc, &topicInfo.UserNUM, &topicInfo.FakeUserNUM, &topicInfo.VideoNUM)
		if topicInfo.Disable == 1 {
			exist, _ := db.HExists(common.TOPICHASH, strconv.FormatUint(topicInfo.TopicID, 10))
			if exist == true {
				db.HDelete(common.TOPICHASH, strconv.FormatUint(topicInfo.TopicID, 10))
			}
			continue
		}

		topicItem := &pb.TopicItem{
			TopicId:        strconv.FormatUint(topicInfo.TopicID, 10),
			Name:           topicInfo.Name,
			IconUrl:        topicInfo.Picture,
			Desc:           topicInfo.Desc,
			SubscribeCount: topicInfo.UserNUM + topicInfo.FakeUserNUM,
		}
		data, err := proto.Marshal(topicItem)
		if err != nil {
			glog.Debug("Serialization videoItem failed")
			continue
		}

		topicPackage = append(topicPackage, strconv.FormatUint(topicInfo.TopicID, 10), data)
		loop++
		if loop == 200 {
			loop = 0
			ret, err := db.HMSet(common.TOPICHASH, topicPackage)
			if ret != "ok" && err != nil {
				glog.Error("UpdateTopicWithID hmset failed. err:%v", err)
			}
			topicPackage = topicPackage[:0]
		}
	}

	if len(topicPackage) != 0 {
		ret, err := db.HMSet(common.TOPICHASH, topicPackage)
		if ret != "ok" && err != nil {
			glog.Error("UpdateTopicWithID hmset failed. err:%v", err)
		}
	}
	glog.Debug("UpdateTopicWithID success")

	return
}

// UpdateTopicVideo load topic-video to redis regularly
func (s *Server) UpdateTopicVideo() {
	rows, err := db.Query(common.BUDAODB, "select topic_id from topic")
	if err != nil {
		glog.Error("UpdateTopicVideo query topic table failed. err:%v", err)
		return
	}
	defer rows.Close()

	var topicID uint64
	for rows.Next() {
		rows.Scan(&topicID)
		tableName, _ := db.GetTableName("topic_video_", topicID)
		sqlString := fmt.Sprintf("select vid, disable, create_time from %s where topic_id = %d", tableName, topicID)
		tempRows, err := db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("UpdateTopicVideo query topic_video_ table failed. err:%v", err)
			continue
		}
		defer tempRows.Close()

		var vid uint64
		var disable uint32
		var vTime time.Time

		var vidArr []interface{}
		loop := 0
		key := fmt.Sprintf("%s%d", common.TOPICWITHVIDEOID, topicID)
		for tempRows.Next() {
			tempRows.Scan(&vid, &disable, &vTime)
			if disable == 1 {
				exist, _ := db.ZExists(key, strconv.FormatUint(vid, 10))
				if exist == true {
					db.ZDelete(key, strconv.FormatUint(vid, 10))
				}
				continue
			}

			loc, _ := time.LoadLocation("Local")
			theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", vTime.Format("2006-01-02 15:04:05"), loc)
			timeStamp := uint64(theTime.Unix())
			vidArr = append(vidArr, timeStamp, strconv.FormatUint(vid, 10))

			loop++
			if loop == 200 {
				loop = 0
				_, err = db.ZAddMulti(key, vidArr)
				if err != nil {
					glog.Error("UpdateTopicVideo add field faild into sort set. err:%v", err)
				}
				vidArr = make([]interface{}, 0)
			}
		}

		if len(vidArr) != 0 {
			_, err = db.ZAddMulti(key, vidArr)
			if err != nil {
				glog.Error("UpdateTopicVideo add field faild into sort set. err:%v", err)
			}
		}
	}
	glog.Debug("UpdateTopicVideo success")

	return
}

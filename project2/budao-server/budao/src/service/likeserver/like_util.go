package likeserver

import (
	"common"
	"db"
	"fmt"

	"github.com/sumaig/glog"
)

// UpdateVideoFavorNum update
func UpdateVideoFavorNum(vid uint64, increment int) (favorNum int, err error) {
	//1.update mysql
	videoTableName, _ := db.GetTableName("video_", vid)
	sqlString := fmt.Sprintf("update %s set favor_num = favor_num+%d where vid = %d", videoTableName, increment, vid)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("like video update video table faild. sqlString:%s, err:%v", sqlString, err)
		return
	}

	//2.update redis
	field := fmt.Sprintf("%s_%d", common.LIKE_COUNT, vid)
	if favorNum, err = db.HIncrBy(common.VIDEO_DYNAMIC, field, increment); err != nil {
		glog.Error("hincrby video like num failed. err:%s", err)
		return
	}

	return
}

package common

import (
	"time"
	pb "twirprpc"
)

//TODO add func use for chech userid, commentid, videoid, topicid
func CheckIDValidity(id uint64) (valid bool) {
	valid = true
	return
}

func GetLocalUnix(t time.Time) uint64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"), loc)
	return uint64(theTime.Unix())
}

func TransStrArrToInterface(strArr []string) []interface{} {
	l := len(strArr)
	retArr := make([]interface{}, l)
	for i := 0; i < l; i++ {
		retArr[i] = strArr[i]
	}
	return retArr
}

/**
 * 获取请求响应的初始状态
 */
func GetInitStatus() (status *pb.Status) {
	return &pb.Status{
		Code:       pb.Status_SERVER_ERR,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}
}

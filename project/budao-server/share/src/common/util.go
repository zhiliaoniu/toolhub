package common

import (
	"time"

	pb "twirprpc"
)

// GetInitStatus 获取请求响应的初始状态
func GetInitStatus() (status *pb.Status) {
	return &pb.Status{
		Code:       pb.Status_SERVER_ERR,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}
}

// TransStrArrToInterface string to interface{} array
func TransStrArrToInterface(strArr []string) []interface{} {
	l := len(strArr)
	retArr := make([]interface{}, l)
	for i := 0; i < l; i++ {
		retArr[i] = strArr[i]
	}
	return retArr
}

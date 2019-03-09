package common

import (
	"time"
	pb "twirprpc"
)

const (
	BUDAODB = "budao"
)

// GetInitStatus 获取请求响应的初始状态
func GetInitStatus() (status *pb.Status) {
	return &pb.Status{
		Code:       pb.Status_SERVER_ERR,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}
}

package util

import (
	"github.com/sumaig/glog"
	"github.com/zhenjl/cityhash"

	"common"
	pb "twirprpc"
)

func GetAbtestItems(header *pb.Header) (abtestItems []*pb.AbtestItem, err error) {
	abtestItems = make([]*pb.AbtestItem, 0)
	abtestConfig := common.GetABTestConfig()
	if abtestConfig == nil {
		glog.Debug("no abtest")
		return
	}

	//1. calc cityhash
	deviceID := header.GetDeviceInfo().GetDeviceId()
	buf := []byte(deviceID)
	seed := uint64(307976497148328517)
	hash := cityhash.CityHash64WithSeed(buf, uint32(len(buf)), seed)
	glog.Debug("cityhash:%d", hash)

	//2. choose abtest
	for _, questionTest := range abtestConfig.QuestionOpenTestMap {
		percent := questionTest.Percent
		glog.Debug("percent:%d", percent)
		if hash%100 < 0 {
			continue
		}
		abtestItem := &pb.AbtestItem{
			Name: "questionOpenTest",
			Id:   "1",
		}
		abtestItems = append(abtestItems, abtestItem)
	}

	return
}

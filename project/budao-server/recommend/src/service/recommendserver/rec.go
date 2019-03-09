package recommendserver

import (
	"common"
	"db"
	"fmt"
	"strings"
	pb "twirprpc"

	"github.com/sumaig/glog"
)

type profile struct {
	devid string
	os    int // 1 ios, 2 android
	hw    string
}

type userRequest struct {
	count     int
	channelId string
}

func play(req *pb.GetRecommendVideoRequest) (ret []index) {
	//1.get user expose vid
	exposeVidMap := make(map[string]bool)
	deviceID := req.GetHeader().GetDeviceInfo().GetDeviceId()
	exposeKey := fmt.Sprintf("%s%s", RECOMMENT_USER_EXPOSE_PREFIX, deviceID)
	exposevidString, err := db.GetString(exposeKey)
	if err != nil {
		if !strings.Contains(err.Error(), common.REDIS_RET_NIL) {
			glog.Error("get user exposekey failed. exposeKey:%s, err:%v", exposeKey, err)
			return
		}
	}
	err = nil
	if exposevidString != "" {
		exposevidArr := strings.Split(exposevidString, ",")
		for _, vid := range exposevidArr {
			exposeVidMap[vid] = true
		}
	}

	//2.get first recommend by condidate
	r := &userRequest{int(req.Count), "timeline"}

	//abtest
	abtestItems := req.GetAbtestItems()
	if abtestItems != nil {
		for _, abtestItem := range abtestItems {
			if abtestItem.Name == "questionOpenTest" {
				ret = commonTrigger(QUESTION, r.count, exposeVidMap)
				if len(ret) == 0 {
					glog.Error("question trigger no new data")
				} else {
					glog.Debug("get %v from question", ret)
					return
				}
			}
		}
	}

	//candidate = hotTrigger(r.count)
	if req.ChannelId == "1" {
		ret = newTrigger(r.count, exposeVidMap)
		if len(ret) == 0 {
			glog.Error("new trigger no new data")
			return
		}
		return
	} else {
		candidate := GetMultiTriger(r.count, exposeVidMap)
		if len(candidate) == 0 {
			glog.Error("multi trigger no new data")
			return
		}

		//3.get final recommend
		ret = rerank(candidate, r.count)
		return
	}

	return
}

const (
	MULTI_COUNT = 200
)

func GetMultiTriger(count int, exposeVidMap map[string]bool) (indexArrMap map[string][]index) {
	indexArrMap = make(map[string][]index, 0)
	//indexArrMap[HOT] = hotTrigger(count*MULTI_COUNT, exposeVidMap)
	//indexArrMap[NEW] = newTrigger(count*MULTI_COUNT, exposeVidMap)
	//indexArrMap[SOURCE_PV] = sourcePvTrigger(count*MULTI_COUNT, exposeVidMap)
	source := []int{1, 4, 5, 6}
	for _, sourceNum := range source {
		sourceKey := fmt.Sprintf("%s%d", SOURCE_PREFIX, sourceNum)
		indexArrMap[sourceKey] = sourceTrigger(sourceNum, count*MULTI_COUNT, exposeVidMap)
		glog.Debug("key:%s, len:%d", sourceKey, len(indexArrMap[sourceKey]))
	}
	indexArrMap[TOPIC] = topicTrigger(count*MULTI_COUNT, exposeVidMap)
	return
}

func indexTrigger(key string, count int) (ret []index) {
	if i, ok := getVideoIndex(key); ok {
		if len(i) < count {
			ret = i
		} else {
			ret = i[:count]
		}
	}
	return
}

func commonTrigger(indexType string, count int, exposeVidMap map[string]bool) (ret []index) {
	glog.Debug("begin %s exposeVidMap.len:%d", indexType, len(exposeVidMap))
	videoIndex, ok := getVideoIndex(indexType)
	if !ok {
		glog.Error("no %s trigger", indexType)
		return
	}
	glog.Debug("%s trigger len:%d", indexType, len(videoIndex))
	num := 1
	for _, index := range videoIndex {
		if _, ok := exposeVidMap[index.id]; !ok {
			ret = append(ret, index)
			num++
			if num > count {
				break
			}
		}
	}
	glog.Debug("end new exposeVidMap.len:%d, ret.len:%d", len(exposeVidMap), len(ret))
	return
}

func newTrigger(count int, exposeVidMap map[string]bool) (ret []index) {
	glog.Debug("begin new exposeVidMap.len:%d", len(exposeVidMap))
	videoIndex, ok := getVideoIndex(NEW)
	if !ok {
		glog.Error("no new trigger")
		return
	}
	//glog.Debug("")
	glog.Debug("new trigger len:%d", len(videoIndex))
	num := 1
	for _, index := range videoIndex {
		//glog.Debug("index:%v", index)
		if _, ok := exposeVidMap[index.id]; !ok {
			//			exposeVidMap[index.id] = true
			ret = append(ret, index)
			num++
			if num > count {
				break
			}
		}
	}
	glog.Debug("end new exposeVidMap.len:%d, ret.len:%d", len(exposeVidMap), len(ret))
	return
}

func sourcePvTrigger(count int, exposeVidMap map[string]bool) (ret []index) {
	glog.Debug("begin sourcePv exposeVidMap.len:%d", len(exposeVidMap))
	videoIndex, ok := getVideoIndex(SOURCE_PV)
	if !ok {
		return
	}
	num := 1
	for _, index := range videoIndex {
		if _, ok := exposeVidMap[index.id]; !ok {
			//			exposeVidMap[index.id] = true
			ret = append(ret, index)
			num++
			if num > count {
				break
			}
		}
	}
	glog.Debug("end sourcePv exposeVidMap.len:%d", len(exposeVidMap))
	return
}

func sourceTrigger(source int, count int, exposeVidMap map[string]bool) (ret []index) {
	sourceKey := fmt.Sprintf("%s%d", SOURCE_PREFIX, source)
	glog.Debug("begin %s exposeVidMap.len:%d", sourceKey, len(exposeVidMap))
	videoIndex, ok := getVideoIndex(sourceKey)
	if !ok {
		return
	}
	num := 1
	for _, index := range videoIndex {
		if _, ok := exposeVidMap[index.id]; !ok {
			//			exposeVidMap[index.id] = true
			ret = append(ret, index)
			num++
			if num > count {
				break
			}
		}
	}
	glog.Debug("end %s exposeVidMap.len:%d", sourceKey, len(exposeVidMap))
	return
}

func topicTrigger(count int, exposeVidMap map[string]bool) (ret []index) {
	glog.Debug("begin topic exposeVidMap.len:%d", len(exposeVidMap))
	videoIndex, ok := getVideoIndex(TOPIC)
	if !ok {
		return
	}
	num := 1
	for _, index := range videoIndex {
		if _, ok := exposeVidMap[index.id]; !ok {
			//			exposeVidMap[index.id] = true
			ret = append(ret, index)
			num++
			if num > count {
				break
			}
		}
	}
	glog.Debug("end topic exposeVidMap.len:%d", len(exposeVidMap))
	shuffle(ret)
	return
}

func hotTrigger(count int, exposeVidMap map[string]bool) (ret []index) {
	glog.Debug("begin hot exposeVidMap.len:%d", len(exposeVidMap))
	videoIndex, ok := getVideoIndex(HOT)
	if !ok {
		return
	}
	num := 1
	for _, index := range videoIndex {
		if _, ok := exposeVidMap[index.id]; !ok {
			//			exposeVidMap[index.id] = true
			ret = append(ret, index)
			num++
			if num > count {
				break
			}
		}
	}
	glog.Debug("end hot exposeVidMap.len:%d", len(exposeVidMap))
	return
}

package timelineserver

import (
	"base"
	"context"
	"db"
	"encoding/json"
	"errors"
	"fmt"
	"service/topicserver"
	"service/util"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
	"github.com/twitchtv/twirp"
	"github.com/zhenjl/cityhash"

	"common"
	pb "twirprpc"
)

const (
	testPerformance = false
)

// Server identify for timeline RPC
type Server struct {
	maxRetVideoNum      int
	retPerTopicViewNum  int
	exposeVideoMaxLen   int
	cronTaskDisable     bool
	cronTaskInternalSec int
	recommendClient     *util.RecommendClient
	logStater           base.LogStater
}

// GetServer return server of gettimeline service
func GetServer() *Server {
	server := &Server{}
	server.initServer()
	return server
}

func (s *Server) initServer() {
	timelineConf := common.GetConfig().TimelineConf
	s.cronTaskDisable = timelineConf.CronTaskDisable
	s.cronTaskInternalSec = timelineConf.CronTaskInternalSec
	s.maxRetVideoNum = timelineConf.MaxRetVideoNum
	s.retPerTopicViewNum = timelineConf.RetPerTopicViewNum
	s.exposeVideoMaxLen = timelineConf.ExposeVideoMaxLen
	s.recommendClient = util.GetRecommendClient()

	s.logStater = base.GetLogStater()

	//init ios online audit vids
	InitIOSOnlineAuditVids()

	//// init
	//go InitVideoDynamicInfo()

	//// init user act
	//go InitUserAct()

	// cron update recommend list of video
	go UpdateVideoRecommendList()

	// cron update listitem by vid
	//go s.CronTask()
}

// Close server
func (s *Server) Close() {
	s.cronTaskDisable = true
}

// GetTimeLine function handle request for timeline.
func (s *Server) GetTimeLine(ctx context.Context, req *pb.GetTimeLineRequest) (resp *pb.GetTimeLineResponse, err error) {
	//1.check req
	beginGetTimeLine := time.Now()
	clientIp, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIp:%v, req:%v", clientIp, req)
	resp = &pb.GetTimeLineResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			s.logStater.IncInternalError("GetTimeLine_ERR", "twirp.status_codes.GetTimeLine.500", 1, 1)
			err = nil
		}
		if len(resp.ListItems) != s.maxRetVideoNum {
			s.logStater.IncInternalError("GetTimeLine_FewVideo", "twirp.status_codes.GetTimeLine.500", 1, 1)
		}
		glog.Debug("resp:%v", resp)
	}()

	//begin test cityhash
	buf, _ := proto.Marshal(req)
	seed := uint64(307976497148328517)
	glog.Debug("cityhash:%d", cityhash.CityHash64WithSeed(buf, uint32(len(buf)), seed))

	deviceID := req.GetHeader().GetDeviceInfo().GetDeviceId()
	if deviceID == "" {
		glog.Error("no deviceID. req:[%v]", req)
		err = errors.New("not set deviceID")
		resp.Status.Code = pb.Status_OK
		return
	}

	//2.get unexpose vids
	vidArr := make([]string, 0)
	beginGetNotExposureVideoList := time.Now()
	exposeKey := fmt.Sprintf("%s%s", common.USEREXPOSURELISTPREFIX, deviceID)
	// get abtest items
	abtestItems, _ := util.GetAbtestItems(req.GetHeader())
	var recommendItems map[string]*pb.RecommendVideoItem
	if req.ChannelId != "4" {
		vidArr, recommendItems, err = GetNotExposureVideoList(req, exposeKey, s.maxRetVideoNum, abtestItems)
		if err != nil {
			glog.Error("get devicid not exposure set failed. exposeKey:%s, err:%v", exposeKey, err)
			return
		}
		if len(vidArr) == 0 {
			resp.Status.Code = pb.Status_OK
			glog.Debug("have no new video. vid is empty")
			err = errors.New("have no new video")
			return
		}
	} else {
		vidArr = append(vidArr, "48656960621889")
		vidArr = append(vidArr, "65628174145044")
		vidArr = append(vidArr, "42897929549039")
		vidArr = append(vidArr, "16225975805997")
		vidArr = append(vidArr, "15247696455879")
		vidArr = append(vidArr, "11007383725071")
		vidArr = append(vidArr, "10978572654864")
		vidArr = append(vidArr, "17024093359294")
		vidArr = append(vidArr, "22742701644900")
	}
	if testPerformance {
		glog.Error("===============GetNotExposureVideoList cost:%v", time.Now().Sub(beginGetNotExposureVideoList))
	}

	//collect exposure vids
	PrintTimelineLog(vidArr, clientIp, req)

	//3.get listItems
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	beginGetListItems := time.Now()
	items, err := util.GetListItemsNotReal(vidArr, userID)
	if testPerformance {
		glog.Error("===============GetListItems cost:%v", time.Now().Sub(beginGetListItems))
	}
	if err != nil {
		glog.Error("get listItems failed. err:%v", err)
		return
	}
	if len(items) == 0 {
		resp.Status.Code = pb.Status_OK
		glog.Debug("have no new video. listItems is empty")
		err = errors.New("have no new video")
		return
	}

	//4.recompose listItems by abtest
	util.RecomposeListItems(items, abtestItems, recommendItems, req.ChannelId)

	resp.ListItems = items
	resp.HasMore = true
	resp.ClearCache = false
	resp.Tips = fmt.Sprintf("为你推荐%d条视频", len(items))
	resp.Status.Code = pb.Status_OK
	if testPerformance {
		glog.Error("===============GetTimeLine cost:%v", time.Now().Sub(beginGetTimeLine))
	}

	return
}

// GetSubscribedTimeLine function handle request for timeline.
func (s *Server) GetSubscribedTimeLine(ctx context.Context, req *pb.GetSubscribedTimeLineRequest) (resp *pb.GetSubscribedTimeLineResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.GetSubscribedTimeLineResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	userId := req.GetHeader().GetUserId()
	if userId == "" {
		glog.Error("bad request. no set userId. req:[%v]", req)
		resp.Status.Code = pb.Status_USER_NOTLOGIN
		err = errors.New("not set userId")
		return
	}

	deviceID := req.GetHeader().GetDeviceInfo().GetDeviceId()
	if deviceID == "" {
		glog.Error("bad request .no deviceID. req:[%v]", req)
		resp.Status.Code = pb.Status_BAD_REQUEST
		err = errors.New("not set deviceID")
		return
	}
	glog.Debug("userId:%s, deviceID:%s", userId, deviceID)

	//refVideoId := req.GetRefVideoId()
	//refVideoTime := req.GetRefVideoTime()
	userIdNum, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)

	//2.get user subscribe topicid
	topicIds, err := topicserver.GetUserSubscribedTopicId(userId)
	if err != nil {
		glog.Error("get user %s subscribed topicid failed. deviceID:%s, err:%v", userId, deviceID, err)
		return
	}
	resp.TopicNum = int32(len(topicIds))
	if len(topicIds) == 0 {
		glog.Debug("user not subscribe any topic")
		resp.Status.Code = pb.Status_OK
		err = errors.New("user not subscribe any topic")
		return
	}
	glog.Debug("user %s subscribe topic:%v", userId, topicIds)
	//topicIdOffsetMap := make(map[string]int, 0)

	//3.get offset of topic, if direction is not Direction_ALL, or offset = 0
	direction := req.GetDirection()
	topicIdMap, err := util.GetTopicVideoOffsetByDirction(userIdNum, topicIds, direction)
	if err != nil {
		glog.Error("get topic video offset failed. err:%v", err)
		return
	}
	glog.Debug("user %s, direction:%d, topic offset:%v", userId, direction, topicIdMap)

	//4.get vid of topic;updateMap use to update offset
	vids, updateMap, hasMore, err := util.GetTopicVidByOffset(userIdNum, topicIdMap, topicIds, direction, s.retPerTopicViewNum)
	if err != nil {
		glog.Error("get topic vid failed. userId:%s, topicIdMap:%v, direction:%d, err:%v", userId, topicIdMap, direction, err)
		return
	}
	resp.HasMore = hasMore
	if len(vids) == 0 {
		glog.Debug("topic has no more video")
		resp.Status.Code = pb.Status_OK
		err = errors.New("topic has no more video")
		return
	}
	glog.Debug("user %s view vids:%v", userId, vids)

	//5.get listItems
	listItems, err := util.GetListItemsNotReal(vids, userIdNum)
	if err != nil {
		glog.Error("get listItems failed. err:%v", err)
		return
	}

	//6.update user topic video offset
	if err = util.UpdateTopicVideoOffsetByDirction(userIdNum, updateMap, direction); err != nil {
		glog.Error("update topic video offset failed. err:%v", err)
	}
	glog.Debug("end request. userId:%s, deviceID:%s", userId, deviceID)

	resp.ListItems = listItems
	resp.ClearCache = false
	resp.Tips = "success"
	resp.Status.Code = pb.Status_OK
	return
}

// GetNotExposureVideoList exposure elimination is repeated
func GetNotExposureVideoList(req *pb.GetTimeLineRequest, exposeKey string, num int, abtestItems []*pb.AbtestItem) (vidArr []string, recommendItemMap map[string]*pb.RecommendVideoItem, err error) {
	//1.get recommend
	recommendList := make([]string, 0)
	if req.GetHeader().GetDeviceInfo().GetDeviceType() == pb.DeviceInfo_IOS &&
		//1.1 get recommend from ios audit
		common.IsVersionBigger(req.GetHeader().GetDeviceInfo().GetAppVersion(), common.GetConfig().IOSAuditConf.AppVersion) {
		recommendList = GIOSAuditVids
	} else {
		//1.2 get recommend from recommendServer
		recommendItems, err := util.GRecommendClient.GetRecommendVideo(req.GetHeader(), req.ChannelId, uint32(num), abtestItems)
		if err != nil || len(recommendItems) == 0 {
			glog.Error("get recommed vid failed. err:%v", err)
		} else {
			recommendItemMap = make(map[string]*pb.RecommendVideoItem, 0)
			for _, recommendItem := range recommendItems {
				recommendList = append(recommendList, recommendItem.VideoId)
				recommendItemMap[recommendItem.VideoId] = recommendItem
			}
		}
		if len(recommendList) == 0 {
			//1.3 get recommend from default
			recommendList = GetRecommendList()
		}
	}

	glog.Debug("final recommend vid:%v, device_id:%s", recommendList, req.GetHeader().GetDeviceInfo().GetDeviceId())
	//2.get user expose vid
	var exposevidArr []string
	exposevidString, err := db.GetString(exposeKey)
	if err != nil {
		if !strings.Contains(err.Error(), common.REDIS_RET_NIL) {
			glog.Error("get user exposekey failed. exposeKey:%s, device_id:%s, err:%v", exposeKey, req.GetHeader().GetDeviceInfo().GetDeviceId(), err)
			return
		}
	}
	err = nil
	if exposevidString != "" {
		exposevidArr = strings.Split(exposevidString, ",")
	}

	//3.get unexpose vid
	vidArr = util.SliceDiffWithNum(recommendList, exposevidArr, num)
	if len(vidArr) == 0 {
		glog.Error("recommendlist have no new video. device_id:%s", req.GetHeader().GetDeviceInfo().GetDeviceId())
		return
	}
	glog.Debug("get unexpose vid:%v", vidArr)

	//4.save exposevid
	go func() {
		exposevidArr = append(exposevidArr, vidArr...)
		exposevidArrLen := len(exposevidArr)
		if exposevidArrLen >= common.GetConfig().TimelineConf.ExposeVideoMaxLen {
			exposevidArr = exposevidArr[exposevidArrLen-common.GetConfig().TimelineConf.ExposeVideoMaxLen:]
		}

		// add vids to expose zset
		value := strings.Join(exposevidArr, ",")
		err = db.SetString(exposeKey, value)
		if err != nil {
			glog.Error("update user:%s failed. err:%v", exposeKey, err)
		}
	}()

	return
}

func PrintTimelineLog(vidArr []string, clientIp string, req *pb.GetTimeLineRequest) {
	var vidStr string
	for _, vid := range vidArr {
		vidStr = vidStr + vid + ","
	}
	vidStr = strings.TrimRight(vidStr, ",")
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), loc)
	userIDStr := req.GetHeader().GetUserId()
	devInfo := req.GetHeader().GetDeviceInfo()
	dev, _ := json.Marshal(devInfo)
	envInfo := req.GetHeader().GetEnvironmentInfo()
	env, _ := json.Marshal(envInfo)

	logStr := fmt.Sprintf("budao-server|%s|exposure|%s|%s|%s|%s|%s", theTime, userIDStr, string(dev), string(env), vidStr, clientIp)
	glog.Info(logStr)
}

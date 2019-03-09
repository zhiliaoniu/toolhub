package topicserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sumaig/glog"

	"common"
	"db"
	"service/util"
	pb "twirprpc"
)

// TopicInfo storage topic info from mysql
type TopicInfo struct {
	AutoID      uint64
	TopicID     uint64
	Name        string
	Picture     string
	Disable     uint32
	Weight      uint32
	Desc        string
	UserNUM     uint32
	FakeUserNUM uint32
	VideoNUM    uint64
}

// Server identify for TopicService RPC
type Server struct {
	cronTaskDisable  bool //定时任务是否可用
	cronTaskInternal int  //定时任务执行间隔

	maxReturnVideoNum   int
	maxReturnTopicNum   int
	maxExplicitVideoNum int
}

// GetServer get server
func GetServer() *Server {
	server := &Server{}
	server.initServer()
	//go server.cronTask()

	return server
}

func (s *Server) initServer() {
	s.cronTaskDisable = false
	s.cronTaskInternal = 10
	s.maxReturnVideoNum = 10
	s.maxReturnTopicNum = 10
	s.maxExplicitVideoNum = 2

	//init ios online audit topicids
	InitIOSOnlineAuditTopicIds()

	//do init
	//go InitUserSubscribedTopic()
	//go InitTopicDynamicInfo()
}

// Close server
func (s *Server) Close() {
	s.cronTaskDisable = true
}

// GetTopicVideoList get video list under the topic
func (s *Server) GetTopicVideoList(context context.Context, req *pb.GetTopicVideoListRequest) (resp *pb.GetTopicVideoListResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.GetTopicVideoListResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	userId, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	topicId, _ := strconv.ParseUint(req.GetTopicId(), 10, 64)
	videoID := req.GetLastVideoId()
	if topicId == 0 {
		glog.Error("not set topicId. req:%v", req)
		resp.Status.Code = pb.Status_BAD_REQUEST
		err = errors.New("not set topicId")
		return
	}

	//2.get vids
	key := fmt.Sprintf("%s%d", common.TOPICWITHVIDEOID, topicId)
	vids, err := GetIdsByLastId(key, videoID, s.maxReturnVideoNum)
	if err != nil {
		glog.Error("get vids of topicid:%d failed", topicId)
		return
	}
	if len(vids) == 0 {
		glog.Error("topic:%d has no video", topicId)
		err = errors.New("topic has no video")
		return
	}
	if len(vids) > s.maxReturnVideoNum {
		vids = vids[:s.maxReturnVideoNum]
		resp.HasMore = true
	}

	//3.getlistItem
	listItems, err := util.GetListItemsNotReal(vids, userId)
	if err != nil {
		glog.Error("get listItems failed. err:%v", err)
		return
	}
	if len(listItems) == 0 {
		glog.Error("user:%d get empty listItem.vids:%s ", userId, vids)
		return
	}

	resp.Status.Code = pb.Status_OK
	for _, listItem := range listItems {
		if listItem.TopicItem != nil {
			if listItem.TopicItem.TopicId == req.GetTopicId() {
				resp.TopicItem = listItems[0].TopicItem
				break
			}
		} else {
			//glog.Error("this_is_bug topic video not have topicinfo.listItem:%v, topicid:%s", listItem, req.GetTopicId())
		}
	}
	if resp.TopicItem == nil {
		glog.Error("this_is_bug topic video not have topicinfo.topicid:%s", req.GetTopicId())
	}

	for _, listItem := range listItems {
		listItem.TopicItem = nil
	}
	resp.ListItems = listItems
	return
}

//2.GetTopicList get topic items
func (s *Server) GetTopicList(ctx context.Context, req *pb.GetTopicListRequest) (resp *pb.GetTopicListResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.GetTopicListResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	userId, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	lastTopicId := req.GetLastTopicId()

	if lastTopicId == "" {
		bannerItems, err := GetBannerItems()
		if err != nil {
			glog.Error("get bannerItems failed. err:%v", err)
		} else {
			resp.BannerItems = bannerItems
		}
	}

	//2.get topic id
	var topicIds []string
	if req.GetHeader().GetDeviceInfo().GetDeviceType() == pb.DeviceInfo_IOS &&
		common.IsVersionBigger(req.GetHeader().GetDeviceInfo().GetAppVersion(), common.GetConfig().IOSAuditConf.AppVersion) {
		topicIds = GIOSAuditTopicIds
		resp.HasMore = false
	} else {
		topicIds, err = GetIdsByLastId(common.TOPICSORTSET, lastTopicId, int(s.maxReturnTopicNum))
		if err != nil {
			glog.Error("get topicids failed. err:%v", err)
			return
		}
		if len(topicIds) == 0 {
			glog.Error("has no more topic")
			err = errors.New("has no more topic")
			resp.HasMore = false
			resp.Status.Code = pb.Status_OK
			return
		}
		if len(topicIds) > s.maxReturnTopicNum {
			topicIds = topicIds[:s.maxReturnTopicNum]
			resp.HasMore = true
		}
	}
	glog.Debug("========len:%d, topicIds:%v", len(topicIds), topicIds)

	//3.get topicItem
	topicIdMap, err := GetTopicItemsWithExplicit(userId, topicIds, s.maxExplicitVideoNum)
	if err != nil {
		return
	}
	glog.Debug("========len(topicIdMap):%d", len(topicIdMap))

	//遍历有序的数组，保证返回的数据也有序，遍历map无序
	for _, topicId := range topicIds {
		if topicItem, ok := topicIdMap[topicId]; ok {
			resp.TopicItems = append(resp.TopicItems, topicItem)
		}
	}

	resp.Status.Code = pb.Status_OK
	return
}

// GetSubscribedTopicList get subscribed topiclist
func (s *Server) GetSubscribedTopicList(ctx context.Context, req *pb.GetSubscribedTopicListRequest) (resp *pb.GetSubscribedTopicListResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.GetSubscribedTopicListResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()
	userId, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	if userId == 0 {
		resp.Status.Code = pb.Status_USER_NOTLOGIN
		err = errors.New("user not login")
		return
	}
	//lastTime := req.GetSubscribeTime()
	lastTopicId := req.GetLastTopicId()

	//2.get topicid
	key := fmt.Sprintf("%s%d", common.USER_SUBSCRIBE_TOPIC_WITH_TIME, userId)
	topicIds, err := GetIdsByLastId(key, lastTopicId, s.maxReturnTopicNum)
	if err != nil {
		glog.Error("get user:%d subscribe topicid list failed. err:%v", userId, err)
		return
	}
	if len(topicIds) == 0 {
		glog.Error("topicIds is empty. key:%s userId:%d", key, userId)
		resp.Status.Code = pb.Status_OK
		return
	}
	resp.HasMore = false
	if len(topicIds) == s.maxReturnTopicNum {
		resp.HasMore = true
	}
	glog.Debug("userId:%d, lastTopicId:%s, topicIds:%v", userId, lastTopicId, topicIds)

	//3.get topicItem
	topicIdMap, err := GetTopicItemsWithExplicit(userId, topicIds, s.maxExplicitVideoNum)
	if err != nil {
		return
	}

	//遍历有序的数组，保证返回的数据也有序，遍历map无序
	for _, topicId := range topicIds {
		if topicItem, ok := topicIdMap[topicId]; ok {
			resp.TopicItems = append(resp.TopicItems, topicItem)
		}
	}
	resp.Status.Code = pb.Status_OK
	return
}

// SubscribeTopic subscribe or unsubscribe topic
func (s *Server) SubscribeTopic(ctx context.Context, req *pb.SubscribeTopicRequest) (resp *pb.SubscribeTopicResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.SubscribeTopicResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	userId, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	token := req.GetHeader().GetToken()

	var valid bool
	if valid, err = util.CheckIDValid(userId, token); !valid {
		resp.Status.Code = pb.Status_USER_NOTLOGIN
		glog.Error("userid:%d, token:%s check not valid. err:%v", userId, token, err)
		return
	}
	topicId, _ := strconv.ParseUint(req.GetTopicId(), 10, 64)
	topicIdStr := req.GetTopicId()
	subscribeAction := req.GetAction()

	//2.update user act
	userTableName, _ := db.GetTableName("user_", userId)
	followTopicTN, _ := db.GetTableName("user_follow_topic_", userId)
	var increment int
	var sqlString string
	if subscribeAction == pb.SubscribeAction_SUBSCRIBE {
		//2.1 insert user act to mysql
		increment = 1
		sqlString = fmt.Sprintf("insert into %s (uid, topic_id) values (%d, %d)", followTopicTN, userId, topicId)
		if _, err = db.Exec(common.BUDAODB, sqlString); err != nil {
			if strings.Contains(err.Error(), common.DUPLICATE_ENTRY) {
				glog.Debug("user:%d has subscribe topic:%d", userId, topicId)
				err = errors.New("user has subscribed this topic")
				resp.Status.Code = pb.Status_OK
				return
			}
			glog.Error("subscribe topic faild. sqlString:%s, err:%v", sqlString, err)
			return
		}

		//2.2 add user subscribed topic zset and view record hash
		key := fmt.Sprintf("%s%d", common.USER_SUBSCRIBE_TOPIC_WITH_TIME, userId)
		now := time.Now().Unix()
		if _, err = db.ZAddMulti(key, []interface{}{now, topicIdStr}); err != nil {
			glog.Error("zadd multi failed. userId:%d, topicId:%d, err:%v", userId, topicId, err)
		}

		key = fmt.Sprintf("%s%d", common.USER_VIEW_TOPIC, (userId % 1000))
		var fields []interface{}
		field1 := fmt.Sprintf("%d_%s_up", userId, topicId)
		field2 := fmt.Sprintf("%d_%s_down", userId, topicId)
		fields = append(fields, field1, "", field2, "")
		if _, err = db.HMSet(key, fields); err != nil {
			glog.Error("hmset failed. userId:%d, topicId:%d, err:%v", userId, topicId, err)
		}
	} else {
		//2.1 delete user act from mysql
		increment = -1
		sqlString = fmt.Sprintf("delete from %s where uid=%d and topic_id=%d", followTopicTN, userId, topicId)
		var result sql.Result
		result, err = db.Exec(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("delete user subscribe faild. sqlString:%s, err:%v", sqlString, err)
		}
		//检查是否真正删除了一条
		var ok bool
		if ok, err = db.CheckDeleteRowSuccess(result); err != nil || !ok {
			glog.Debug("user:%s not subscribe topic:%s", userId, topicId)
			err = errors.New("user not subscribe this topic")
			resp.Status.Code = pb.Status_OK
			return
		}

		//2.2delete user subscribed topic from zset and hash
		key := fmt.Sprintf("%s%d", common.USER_SUBSCRIBE_TOPIC_WITH_TIME, userId)
		_, err = db.ZDelete(key, topicIdStr)
		if err != nil {
			glog.Error("zrem failed. userId:%d, topicId:%d, err:%v", userId, topicId, err)
		}

		key = fmt.Sprintf("%s%d", common.USER_VIEW_TOPIC, (userId % 1000))
		var fields []interface{}
		field1 := fmt.Sprintf("%d_%s_up", userId, topicId)
		field2 := fmt.Sprintf("%d_%s_down", userId, topicId)
		fields = append(fields, field1, field2)
		if _, err = db.HMDelete(key, fields); err != nil {
			glog.Error("hdel failed. userId:%d, topicId:%d, err:%v", userId, topicId, err)
		}
	}

	//3.update user subscrebe topic num in mysql
	sqlString = fmt.Sprintf("update %s set follow_topic_num=follow_topic_num+%d where uid=%d", userTableName, increment, userId)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("update user subscribe topic num faild. sqlString:%s, err:%v", sqlString, err)
		return
	}

	//4.update topic subscribe num
	//4.1 update mysql
	sqlString = fmt.Sprintf("update topic set user_num=user_num+%d where topic_id = %d", increment, topicId)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("update topic subscribe user num faild. sqlString:%s, err:%v", sqlString, err)
		return
	}
	//4.2 update user act in redis
	key := fmt.Sprintf("%s%d", common.USER_ACT_PREFIX, userId)
	field := fmt.Sprintf("%d%s", topicId, common.TOPIC_SUBSCRIBE_SUFFIX)
	count, err := db.HIncrBy(key, field, increment)
	if err != nil {
		glog.Error("update user subscribe topic num in redis failed")
	}

	resp.Count = uint32(count)
	resp.Status.Code = pb.Status_OK
	return
}

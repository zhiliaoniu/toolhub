package util

import (
	"common"
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/sumaig/glog"

	pb "twirprpc"
)

var (
	GRecommendClient    *RecommendClient
	onceRecommendClient sync.Once
)

func GetRecommendClient() *RecommendClient {
	onceRecommendClient.Do(func() {
		GRecommendClient = &RecommendClient{}
		GRecommendClient.initRecommendClient()
	})
	return GRecommendClient
}

type RecommendClient struct {
	reqTimeoutMs int    //请求超时毫秒
	addr         string //推荐系统地址
	//recommendClient *pb.RecommendService
}

func (s *RecommendClient) initRecommendClient() {
	recommendClientConf := common.GetConfig().RecommendClientConf
	s.reqTimeoutMs = recommendClientConf.ReqTimeoutMs
	s.addr = recommendClientConf.Addr //"http://localhost:8001" //":8001"
	//s.recommendClient = pb.NewRecommendServiceProtobufClient(addr, &http.Client{})
}

func (s *RecommendClient) GetRecommendVideo(header *pb.Header, channelId string, count uint32, abtestItems []*pb.AbtestItem) (recommendVideoItems []*pb.RecommendVideoItem, err error) {
	glog.Debug("header:%v, channelId:%s", header, channelId)
	//1.get req num of session
	deviceID := header.GetDeviceInfo().GetDeviceId()
	var sessionReqNum uint32
	if sessionReqNum, err = GetUserSessionReqNum(deviceID); err != nil {
		glog.Error("get header:%v sessionReqNum failed. err:%v", header, err)
		sessionReqNum = 0
	}

	//2.prepare req
	req := &pb.GetRecommendVideoRequest{
		Header:        header,
		ChannelId:     channelId,
		Count:         count,
		SessionReqNum: sessionReqNum,
		AbtestItems:   abtestItems,
	}

	//send req
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.reqTimeoutMs)*time.Millisecond)
	defer cancel()
	ch := make(chan bool)
	var resp *pb.GetRecommendVideoResponse
	go func() {
		client := pb.NewRecommendServiceProtobufClient(s.addr, &http.Client{})
		resp, err = client.GetRecommendVideo(ctx, req)
		ch <- true
	}()

	//wait resp
	select {
	case <-ch:
		break
	case <-ctx.Done():
		glog.Error("req timeout. req:%v, timeoutMS:%d", req, s.reqTimeoutMs)
		return nil, ctx.Err()
	}
	if err != nil {
		glog.Error("reqeust recommend video failed. req:%v, resp:%v, err:%v", req, resp, err)
		return
	}

	//3.parse resp
	glog.Debug("request recommend video resp:%v", resp)
	if resp.Status.Code != pb.Status_OK {
		glog.Error("resp status not ok. resp:%v", resp)
		err = errors.New("resp status not ok")
		return
	}
	recommendVideoItems = resp.RecommendVideoItems

	return
}

//TODO
func GetUserSessionReqNum(deviceId string) (sessionReqNum uint32, err error) {
	return 1, nil
}

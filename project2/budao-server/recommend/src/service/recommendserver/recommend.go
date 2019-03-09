package recommendserver

import (
	"common"
	"context"
	"errors"
	_ "net/http/pprof"
	"strconv"
	"time"

	pb "twirprpc"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
)

const (
	RECOMMENT_USER_EXPOSE_PREFIX = "exposure_list_deviceid_"
)

type Server struct {
}

func GetServer() *Server {
	server := &Server{}
	server.initServer()
	return server
}

func (s *Server) initServer() {
	//recommendConf := common.GetConfig().RecommendConf

	InitIndex()
}

func (s *Server) Close() {
}

func (s *Server) GetRecommendVideo(ctx context.Context, req *pb.GetRecommendVideoRequest) (resp *pb.GetRecommendVideoResponse, err error) {
	glog.Debug("req:%v", req)
	resp = &pb.GetRecommendVideoResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	//1.get recommend vid TODO
	start := time.Now()
	indexArr := play(req)
	glog.Debug("play cost:%v", time.Since(start))
	if len(indexArr) == 0 {
		err = errors.New("no data")
		return
	}
	glog.Debug("recommendList.len:%d", len(indexArr))

	//2.compose resp
	tarceID := strconv.FormatUint(uint64(time.Now().UnixNano()), 10)
	payload := &pb.Payload{
		AbtestItems: req.GetAbtestItems(),
	}
	payloadStr, err := proto.Marshal(payload)
	if err != nil {
		glog.Error("proto marshal payload failed. err:%v", err)
	}
	for _, index := range indexArr {
		recommendVideoItem := &pb.RecommendVideoItem{
			VideoId: index.id,
			Payload: string(payloadStr),
			TraceId: tarceID,
		}
		resp.RecommendVideoItems = append(resp.RecommendVideoItems, recommendVideoItem)
	}
	resp.Status.Code = pb.Status_OK

	return
}

func (s *Server) ReloadIndex(ctx context.Context, req *pb.ReloadIndexRequest) (resp *pb.ReloadIndexResponse, err error) {
	glog.Debug("req:%v", req)
	resp = &pb.ReloadIndexResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	indexPath := req.GetIndexPath()
	if err = DoReloadIndex(indexPath); err != nil {
		return
	}

	resp.Status.Code = pb.Status_OK
	return
}

package reportserver

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/robfig/cron"
	"github.com/sumaig/glog"
	twirp "github.com/twitchtv/twirp"

	"common"
	pb "twirprpc"
)

// Server identify for ReportService RPC
type Server struct {
	c *cron.Cron
}

// GetServer return server of report service
func GetServer() *Server {
	server := &Server{}
	server.initServer()

	return server
}

func (s *Server) initServer() {
	go s.cronTask()
}

// Close server
func (s *Server) Close() {
	s.c.Stop()
}

// ReportStatisData client report statistic data
func (s *Server) ReportStatisData(ctx context.Context, req *pb.ReportStatisDataRequest) (resp *pb.ReportStatisDataResponse, err error) {
	resp = &pb.ReportStatisDataResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	var (
		flag           int32
		duration       uint32
		playDuration   float32
		targetDuration float32
		bannerID       string
		videoBase      *pb.VideoBase
	)
	switch req.Type.(type) {
	case *pb.ReportStatisDataRequest_Click:
		flag = 1
		videoBase = req.GetClick().GetBase()
	case *pb.ReportStatisDataRequest_Progress:
		flag = 2
		videoBase = req.GetProgress().GetBase()
		duration = req.GetProgress().GetDuration()
		playDuration = req.GetProgress().GetPlayDuration()
		if playDuration != 0 {
			targetDuration = playDuration
			break
		}
		if duration != 0 {
			targetDuration = float32(duration)
		}
	case *pb.ReportStatisDataRequest_Exposure:
		flag = 3
		videoBase = req.GetExposure().GetBase()
	case *pb.ReportStatisDataRequest_BannerClick:
		flag = 4
		bannerID = req.GetBannerClick().GetBannerId()
	}
	glog.Debug("flag:%d", flag)

	var (
		vid, channelID, traceID, payload string
		trace                            uint64
	)
	switch req.Type.(type) {
	case *pb.ReportStatisDataRequest_Click, *pb.ReportStatisDataRequest_Progress, *pb.ReportStatisDataRequest_Exposure:
		vid = videoBase.GetVideoId()
		channelID = videoBase.GetChannelId()
		traceID = videoBase.GetTraceId()
		payload = videoBase.GetPayload()
		trace, _ = strconv.ParseUint(traceID, 10, 64)
		glog.Debug("flag:%d, vid:%s, traceid:%s, channelid:%s", flag, vid, traceID, channelID)

		if vid == "" || traceID == "" || channelID == "" {
			glog.Error("request lack para")
			resp.Status.Code = pb.Status_BAD_REQUEST
			return
		}
		if payload == "" {
			payload = "客户端透传信息"
		}
	case *pb.ReportStatisDataRequest_BannerClick:
		glog.Debug("flag:%d, bannerid:%s", flag, bannerID)
		if bannerID == "" {
			glog.Error("request lack para")
			resp.Status.Code = pb.Status_BAD_REQUEST
			return
		}
	}

	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), loc)

	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	devInfo := req.GetHeader().GetDeviceInfo()
	dev, _ := json.Marshal(devInfo)
	envInfo := req.GetHeader().GetEnvironmentInfo()
	env, _ := json.Marshal(envInfo)

	clientIP, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIP:%v", clientIP)

	// record log
	var logStr string
	switch req.Type.(type) {
	case *pb.ReportStatisDataRequest_Click, *pb.ReportStatisDataRequest_Exposure:
		logStr = fmt.Sprintf("statistic-server|%s|%d|%d|%s|%s|%s|%d|%s|%s|%s", theTime, flag, userID, string(dev), string(env), vid, trace, channelID, payload, clientIP)
	case *pb.ReportStatisDataRequest_Progress:
		logStr = fmt.Sprintf("statistic-server|%s|%d|%d|%s|%s|%s|%d|%s|%s|%f|%s", theTime, flag, userID, string(dev), string(env), vid, trace, channelID, payload, targetDuration, clientIP)
	case *pb.ReportStatisDataRequest_BannerClick:
		logStr = fmt.Sprintf("statistic-server|%s|%d|%d|%s|%s|%s|%s", theTime, flag, userID, string(dev), string(env), bannerID, clientIP)
	}
	glog.Info(logStr)

	resp.Status.Code = pb.Status_OK

	return
}

// ReportBatchData client batch report data
func (s *Server) ReportBatchData(ctx context.Context, req *pb.ReportBatchDataRequest) (resp *pb.ReportBatchDataResponse, err error) {
	resp = &pb.ReportBatchDataResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), loc)
	userID := req.GetHeader().GetUserId()
	devInfo := req.GetHeader().GetDeviceInfo()
	dev, _ := json.Marshal(devInfo)
	envInfo := req.GetHeader().GetEnvironmentInfo()
	env, _ := json.Marshal(envInfo)

	clientIP, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIP:%v", clientIP)

	if len(req.Click) != 0 {
		for _, click := range req.GetClick() {
			vid := click.GetBase().GetVideoId()
			channelID := click.GetBase().GetChannelId()
			traceID := click.GetBase().GetTraceId()
			payload := click.GetBase().GetPayload()

			if vid == "" || traceID == "" || channelID == "" {
				glog.Error("request lack para")
				continue
			}

			if payload == "" {
				payload = "客户端透传信息"
			}

			flag := 1
			logStr := fmt.Sprintf("statistic-server|%s|%d|%s|%s|%s|%s|%s|%s|%s|%s", theTime, flag, userID, string(dev), string(env), vid, traceID, channelID, payload, clientIP)
			glog.Info(logStr)
		}
	}

	if len(req.Progress) != 0 {
		for _, process := range req.GetProgress() {
			vid := process.GetBase().GetVideoId()
			channelID := process.GetBase().GetChannelId()
			traceID := process.GetBase().GetTraceId()
			payload := process.GetBase().GetPayload()

			if vid == "" || traceID == "" || channelID == "" {
				glog.Error("request lack para")
				continue
			}

			duration := process.GetDuration()
			playDuration := process.GetPlayDuration()
			var targetDuration float32
			if duration != 0 {
				targetDuration = float32(duration)
			}
			if playDuration != 0 {
				targetDuration = playDuration
			}

			if payload == "" {
				payload = "客户端透传信息"
			}

			flag := 2
			logStr := fmt.Sprintf("statistic-server|%s|%d|%s|%s|%s|%s|%s|%s|%s|%f|%s", theTime, flag, userID, string(dev), string(env), vid, traceID, channelID, payload, targetDuration, clientIP)
			glog.Info(logStr)
		}
	}

	if len(req.Exposure) != 0 {
		for _, exposure := range req.GetExposure() {
			vid := exposure.GetBase().GetVideoId()
			channelID := exposure.GetBase().GetChannelId()
			traceID := exposure.GetBase().GetTraceId()
			payload := exposure.GetBase().GetPayload()

			if vid == "" || traceID == "" || channelID == "" {
				glog.Error("request lack para")
				continue
			}

			if payload == "" {
				payload = "客户端透传信息"
			}

			flag := 3
			logStr := fmt.Sprintf("statistic-server|%s|%d|%s|%s|%s|%s|%s|%s|%s|%s", theTime, flag, userID, string(dev), string(env), vid, traceID, channelID, payload, clientIP)
			glog.Info(logStr)
		}
	}

	resp.Status.Code = pb.Status_OK

	return
}

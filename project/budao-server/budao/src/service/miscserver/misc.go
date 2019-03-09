package miscserver

import (
	"common"
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"service/util"
	"strconv"
	"strings"
	"time"

	"github.com/sumaig/glog"

	"db"
	pb "twirprpc"
)

// Server identify for Miscservice RPC
type Server struct{}

// ReportVideo user report video
func (s *Server) ReportVideo(ctx context.Context, req *pb.ReportVideoRequest) (resp *pb.ReportVideoResponse, err error) {
	glog.Debug("req:%v", req)
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	deviceId := req.GetHeader().GetDeviceInfo().GetDeviceId()
	videoID := req.GetVideoId()
	reason := req.GetReason()

	status := &pb.Status{
		Code:       pb.Status_OK,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}

	sqlString := fmt.Sprintf("insert into user_report_video (uid, device_id, vid, reason) values (%d, '%s', %s, '%s')", userID, deviceId, videoID, reason)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("report video insert user_report_video table faild. sqlString:%s, err:%v", sqlString, err)
		status.Code = pb.Status_SERVER_ERR
		status.Message = err.Error()
	}

	return &pb.ReportVideoResponse{
		Status: status,
	}, nil
}

// ReportComment user report comment
func (s *Server) ReportComment(ctx context.Context, req *pb.ReportCommentRequest) (resp *pb.ReportCommentResponse, err error) {
	glog.Debug("req:%v", req)

	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	deviceID := req.GetHeader().GetDeviceInfo().GetDeviceId()
	commentID := req.GetCommentId()
	reason := req.GetReason()

	status := &pb.Status{
		Code:       pb.Status_OK,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}

	sqlString := fmt.Sprintf("insert into user_report_comment (uid, device_id, cid, reason) values (%d, '%s', %s, '%s')", userID, deviceID, commentID, reason)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("report comment insert user_report_comment table faild. sqlString:%s, err:%v", sqlString, err)
		status.Code = pb.Status_SERVER_ERR
		status.Message = err.Error()
	}

	return &pb.ReportCommentResponse{
		Status: status,
	}, nil
}

// UserFeedback user feed info
func (s *Server) UserFeedback(ctx context.Context, req *pb.UserFeedbackRequest) (resp *pb.UserFeedbackResponse, err error) {
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	deviceId := req.GetHeader().GetDeviceInfo().GetDeviceId()
	feedback := req.GetFeedback()
	contact := req.GetContact()

	status := &pb.Status{
		Code:       pb.Status_OK,
		Message:    "success",
		ServerTime: uint64(time.Now().Unix()),
	}

	sqlString := fmt.Sprintf("insert into user_feedback(uid, device_id, contact, feedback) values (%d, '%s', '%s', '%s')", userID, deviceId, contact, feedback)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("insert user feedback faild. sqlString:%s, err:%v", sqlString, err)
		status.Code = pb.Status_SERVER_ERR
	}

	return &pb.UserFeedbackResponse{
		Status: status,
	}, nil
}

//客户端推送手机的设备id和补刀app的token到服务端，用于以后服务端给客户端推送数据。设备id作为唯一主键
func (s *Server) PostDeviceToken(ctx context.Context, req *pb.PostDeviceTokenRequest) (resp *pb.PostDeviceTokenResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.PostDeviceTokenResponse{
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
	deviceId := req.GetHeader().GetDeviceInfo().GetDeviceId()
	deviceToken := req.GetDeviceToken()
	osType := req.GetHeader().GetDeviceInfo().GetDeviceType()
	if deviceId == "" || deviceToken == "" {
		glog.Error("not set deviceID or deviceToken. req:[%v]", req)
		err = errors.New("not set deviceID or deviceToken")
		resp.Status.Code = pb.Status_BAD_REQUEST
		return
	}

	//2.choose table by crc32(deviceId)
	tableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["user_token_"]
	crc := crc32.ChecksumIEEE([]byte(deviceId))
	tableIndex := crc % uint32(tableNum)
	tableName := fmt.Sprintf("%s%d", "user_token_", tableIndex)

	//3.先插入，通过插入判断是否重复
	sqlString := fmt.Sprintf("insert into %s (uid, device_id, device_token, os_type) values (%d, '%s', '%s', %d)", tableName, userId, deviceId, deviceToken, osType)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		if strings.Contains(err.Error(), common.DUPLICATE_ENTRY) {
			//4.如果已经存在，更新token
			sqlString := fmt.Sprintf("update %s set uid=%d, device_token='%s', os_type=%d where device_id='%s'", tableName, userId, deviceToken, osType, deviceId)
			_, err = db.Exec(common.BUDAODB, sqlString)
			if err != nil {
				glog.Error("update device_id token failed. sqlString:%s, err:%v", sqlString, err)
				return
			}
			resp.Status.Code = pb.Status_OK
			return
		}
		glog.Error("insert %s table faild. sqlString:%s, err:%v", tableName, sqlString, err)
		return
	}
	return
}

func (s *Server) GetRemoteConfig(ctx context.Context, req *pb.GetRemoteConfigRequest) (resp *pb.GetRemoteConfigResponse, err error) {
	glog.Debug("req:%v", req)
	resp = &pb.GetRemoteConfigResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	resp.Status.Code = pb.Status_OK
	return
}

//function:push info to client
func (s *Server) GetVideoInfo(ctx context.Context, req *pb.GetVideoInfoRequest) (resp *pb.GetVideoInfoResponse, err error) {
	glog.Debug("req:%v", req)
	resp = &pb.GetVideoInfoResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	//userId := req.GetHeader().GetUserId()
	userIdNum, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	videoId := req.GetVideoId()

	//1.get listItems
	listItems, err := util.GetListItemsNotReal([]string{videoId}, userIdNum)
	if err != nil {
		glog.Error("get listItems failed. err:%v", err)
		return
	}
	if len(listItems) == 0 {
		glog.Error("not find video")
		err = errors.New("not find video")
		return
	}

	resp.ListItem = listItems[0]
	resp.Status.Code = pb.Status_OK
	return
}

func (s *Server) ReloadConf(ctx context.Context, req *pb.ReloadConfRequest) (resp *pb.ReloadConfResponse, err error) {
	glog.Debug("req:%v", req)
	resp = &pb.ReloadConfResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	//1.reload conf
	confPath := req.GetConfPath()
	glog.Debug("confPath:%s", confPath)
	if err = common.ParseABTestConfig(confPath); err != nil {
		glog.Error("reload conf failed. confPath:%s, err:%v", confPath, err)
		return
	}

	resp.Status.Code = pb.Status_OK
	return
}

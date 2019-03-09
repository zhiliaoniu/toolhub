package likeserver

import (
	"context"
	"database/sql"
	"fmt"
	"service/commentserver"
	"service/util"
	"strconv"
	"strings"

	"github.com/sumaig/glog"

	"common"
	"db"
	pb "twirprpc"
)

// Server identify for LikeService RPC
type Server struct{}

// GetServer return server of like service
func GetServer() *Server {
	server := &Server{}

	return server
}

// LikeVideo like video
func (s *Server) LikeVideo(ctx context.Context, req *pb.LikeVideoRequest) (resp *pb.LikeVideoResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.LikeVideoResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	videoID, _ := strconv.ParseUint(req.GetVideoId(), 10, 64)
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	likeAction := req.GetAction()
	token := req.GetHeader().GetToken()

	var valid bool
	if valid, err = util.CheckIDValid(userID, token); !valid {
		resp.Status.Code = pb.Status_USER_NOTLOGIN
		glog.Error("userid:%d, token:%s check not valid. err:%v", userID, token, err)
		return
	}

	userTableName, _ := db.GetTableName("user_", userID)
	usrFavVideoName, _ := db.GetTableName("user_favor_video_", userID)

	//2.update user op record
	var sqlString string
	increment := 1
	if likeAction == pb.LikeAction_LIKE {
		//先插入，通过插入判断是否重复点赞
		//0:无操作 1:赞 2:踩
		sqlString = fmt.Sprintf("insert into %s (uid, vid, result) values (%d, %d, %d)", usrFavVideoName, userID, videoID, 1)
		_, err = db.Exec(common.BUDAODB, sqlString)
		if err != nil {
			if strings.Contains(err.Error(), common.DUPLICATE_ENTRY) {
				resp.Status.Code = pb.Status_OK
				return
			}
			glog.Error("like video faild. sqlString:%s, err:%v", sqlString, err)
			return
		}
	} else {
		increment = -1
		sqlString = fmt.Sprintf("delete from %s where uid=%d and vid=%d", usrFavVideoName, userID, videoID)
		var result sql.Result
		result, err = db.Exec(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("unlike video faild. sqlString:%s, err:%v", sqlString, err)
			return
		}

		//检查是否真正删除了一条
		var ok bool
		if ok, err = db.CheckDeleteRowSuccess(result); err != nil || !ok {
			return
		}
	}

	//3 update video favor num and user favor num
	favorNum, err := UpdateVideoFavorNum(videoID, increment)
	if err != nil {
		glog.Error("update video favor num  failed. err:%v", err)
		return
	}

	sqlString = fmt.Sprintf("update %s set video_favor_num=video_favor_num+%d where uid = %d", userTableName, increment, userID)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("like video update user table faild. err:%v", err)
		return
	}

	// update user_act_[uid] hash
	userActKey := fmt.Sprintf("%s%d", common.USER_ACT_PREFIX, userID)
	userActField := fmt.Sprintf("%d%s", videoID, common.VID_FAVOR_SUFFIX)
	if likeAction == pb.LikeAction_LIKE {
		_, err = db.HSet(userActKey, userActField, 1)
		if err != nil {
			glog.Error("insert user act likevideo hash failed. err:%v", err)
			return
		}
	} else {
		_, err = db.HDelete(userActKey, userActField)
		if err != nil {
			glog.Error("delete user act likevideo hash failed. err:%v", err)
			return
		}
	}

	// update vid_dynamic
	_, err = db.SAdd(common.VID_DYNAMIC, req.GetVideoId())
	if err != nil {
		glog.Error("insert vid_dynamic set failed. err:%v", err)
		return
	}

	resp.Status.Code = pb.Status_OK
	resp.Count = uint32(favorNum)

	return
}

// LikeComment like comment
func (s *Server) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (resp *pb.LikeCommentResponse, err error) {
	//1.parse req
	glog.Debug("req:%v", req)
	resp = &pb.LikeCommentResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	commentID, _ := strconv.ParseUint(req.GetCommentId(), 10, 64)
	likeAction := req.GetAction()
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	token := req.GetHeader().GetToken()
	var valid bool
	if valid, err = util.CheckIDValid(userID, token); !valid {
		resp.Status.Code = pb.Status_USER_NOTLOGIN
		glog.Error("userid:%d, token:%s check not valid. err:%v", userID, token, err)
		return
	}

	userTableName, _ := db.GetTableName("user_", userID)
	usrFavCommentName, _ := db.GetTableName("user_favor_comment_", userID)

	//2.update user op record
	var sqlString string
	increment := 1
	if likeAction == pb.LikeAction_LIKE {
		//先插入，通过插入判断是否重复点赞
		sqlString = fmt.Sprintf("insert into %s (uid, cid, result) values (%d, %d, %d)", usrFavCommentName, userID, commentID, 1)
		_, err = db.Exec(common.BUDAODB, sqlString)
		if err != nil {
			if strings.Contains(err.Error(), common.DUPLICATE_ENTRY) {
				resp.Status.Code = pb.Status_OK
				return
			}
			glog.Error("like comment insert user_favor_comment table faild. err:%v", err)
			return
		}
	} else {
		increment = -1
		sqlString = fmt.Sprintf("delete from %s where uid=%d and cid=%d", usrFavCommentName, userID, commentID)
		var result sql.Result
		if result, err = db.Exec(common.BUDAODB, sqlString); err != nil {
			glog.Error("unlike comment delete user_favor_comment table faild. sqlString:%s, err:%v", sqlString, err)
			return
		}
		//检查是否真正删除了一条
		var ok bool
		if ok, err = db.CheckDeleteRowSuccess(result); err != nil || !ok {
			return
		}
	}

	//3 update comment favor num
	vid, _ := commentserver.GetCommentVid(req.GetCommentId())
	favorNum, err := commentserver.UpdateCommentFavorNum(commentID, vid, increment)
	if err != nil {
		glog.Error("update comment favor num  failed. err:%v", err)
		return
	}

	sqlString = fmt.Sprintf("update %s set comment_favor_num=comment_favor_num+%d where uid = %d", userTableName, increment, userID)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("like comment update user table faild. err:%v", err)
		return
	}

	// update user_act_[uid] hash
	userActKey := fmt.Sprintf("%s%d", common.USER_ACT_PREFIX, userID)
	userActField := fmt.Sprintf("%d%s", commentID, common.CID_FAVOR_SUFFIX)
	if likeAction == pb.LikeAction_LIKE {
		_, err = db.HSet(userActKey, userActField, 1)
		if err != nil {
			glog.Error("insert user act likecomment hash failed. err:%v", err)
			return
		}
	} else {
		_, err = db.HDelete(userActKey, userActField)
		if err != nil {
			glog.Error("delete user act likecomment hash failed. err:%v", err)
			return
		}
	}

	// update vid_dynamic
	_, err = db.SAdd(common.VID_DYNAMIC, strconv.FormatUint(vid, 10))
	if err != nil {
		glog.Error("insert vid_dynamic set failed. err:%v", err)
		return
	}

	resp.Status.Code = pb.Status_OK
	resp.Count = uint32(favorNum)

	return
}

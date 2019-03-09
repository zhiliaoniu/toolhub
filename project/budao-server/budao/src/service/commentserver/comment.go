package commentserver

import (
	"base"
	"common"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
	"github.com/twitchtv/twirp"

	"db"
	"service/transfer"
	"service/util"
	pb "twirprpc"
)

func GetServer() *Server {
	server := &Server{}
	server.initServer()
	return server
}

// Server identify for comment RPC
type Server struct {
	idGenerator *transfer.IdGenerator
	logStater   base.LogStater

	cronTaskDisable     bool //定时任务是否可用
	cronTaskInternalSec int  //评论相关定时任务执行间隔
	returnHotCommentNum int  //第一次返回多少条热评论
	returnNewCommentNum int  //每次返回多少条新评论
	returnReplyNum      int  //每次返回多少条回复
	explicitReplyNum    int  //第一次取评论外显几条回复
}

func (s *Server) initServer() {
	s.idGenerator = transfer.GetIdGenerator()
	s.logStater = base.GetLogStater()

	commentConf := common.GetConfig().CommentConf
	s.cronTaskDisable = commentConf.CronTaskDisable
	s.cronTaskInternalSec = commentConf.CronTaskInternalSec
	s.returnHotCommentNum = commentConf.ReturnHotCommentNum
	s.returnNewCommentNum = commentConf.ReturnCommentNum
	s.returnReplyNum = commentConf.ReturnReplyNum
	s.explicitReplyNum = commentConf.ExplicitReplyNum

	//go s.cronTask()

	//执行更新任务
	//go InitCommentInfoByTraverseMysql()
}

func (s *Server) Close() {
	s.cronTaskDisable = true
}

//获取视频评论
func (s *Server) GetVideoCommentList(ctx context.Context, req *pb.GetVideoCommentListRequest) (resp *pb.GetVideoCommentListResponse, err error) {
	clientIp, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIp:%s, req:%+v", clientIp, req)

	videoID, _ := strconv.ParseUint(req.GetVideoId(), 10, 64)
	//lastCommentTime := req.GetLastCommentTime()
	lastCommentId := req.GetLastCommentId()
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)

	resp = &pb.GetVideoCommentListResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	//1.get new comment item
	vcnewKey := fmt.Sprintf("%s%d", util.VCNEW_PREFIX, videoID)
	commentIds, err := GetCommentIdsByLastCommentId(vcnewKey, lastCommentId, s.returnNewCommentNum)
	if err != nil {
		glog.Error("get comment id from redis failed. err:%v", err)
		return
	}
	if len(commentIds) == 0 {
		glog.Debug("video %d has no comment.", videoID)
		resp.Status.Code = pb.Status_OK
		return
	}

	glog.Debug("commentIds:%v, len:%d", commentIds, len(commentIds))
	if len(commentIds) > s.returnNewCommentNum {
		commentIds = append(commentIds[:s.returnNewCommentNum])
		resp.HasMore = true
	}
	glog.Debug("after del. commentIds:%v, len:%d", commentIds, len(commentIds))
	resp.CommentItems, err = GetCommentItemsByCommentIds(userID, commentIds, videoID, int64(s.explicitReplyNum))
	if err != nil {
		glog.Error("get new comment item failed. err:%v", err)
		return
	}

	//第一次来，返回热评论
	for lastCommentId == "" {
		vchotKey := fmt.Sprintf("%s%d", util.VCHOT_PREFIX, videoID)
		hotCommentIds, err := db.ZRevRange(vchotKey, 0, int64(s.returnHotCommentNum-1))
		if err != nil {
			glog.Debug("get hot comment id from redis failed. err:%v", err)
			break
		}
		if len(hotCommentIds) == 0 {
			glog.Debug("video %d has no hot comment.", videoID)
			break
		}
		resp.HotCommentItems, err = GetCommentItemsByCommentIds(userID, hotCommentIds, videoID, int64(s.explicitReplyNum))
		if err != nil {
			glog.Error("get hot comment item failed. err:%v", err)
			break
		}
		break
	}

	resp.Status.Code = pb.Status_OK
	return
}

//获取评论回复
func (s *Server) GetCommentReplyList(ctx context.Context, req *pb.GetCommentReplyListRequest) (resp *pb.GetCommentReplyListResponse, err error) {
	clientIp, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIp:%s, req:%+v", clientIp, req)
	commentID, _ := strconv.ParseUint(req.GetCommentId(), 10, 64)
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	//lastReplyTime := req.GetLastReplyTime()
	lastReplyId := req.GetLastReplyId()

	resp = &pb.GetCommentReplyListResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	//1.get reply of comment
	crnewKey := fmt.Sprintf("%s%d", util.CRNEW_PREFIX, commentID)
	replyIds, err := GetCommentIdsByLastCommentId(crnewKey, lastReplyId, s.returnReplyNum)
	if err != nil {
		glog.Error("get reply id from redis failed. err:%v", err)
		return
	}
	if len(replyIds) == 0 {
		glog.Debug("has no reply more. cid:%d", commentID)
		resp.Status.Code = pb.Status_OK
		return
	}
	if len(replyIds) > s.returnReplyNum {
		replyIds = append(replyIds[:s.returnReplyNum])
		resp.HasMore = true
	}

	resp.ReplyItems, err = GetMultiReplyFromRedis(replyIds)
	if err != nil {
		glog.Error("get reply item from redis failed. err:%v", err)
		return
	}

	cids := make([]string, 0)
	//2.get comment
	commentItems, err := util.GetMultiCommentFromRedis([]string{strconv.FormatUint(commentID, 10)})
	if err != nil || len(commentItems) == 0 {
		glog.Error("not find comment item. commentid:%d err:%v", commentID, err)
	} else {
		cids = append(cids, req.GetCommentId())
		resp.CommentItem = commentItems[0]
	}

	//3.从mysql中读取当前用户是否对这些评论点过赞
	for _, replyItem := range resp.ReplyItems {
		cids = append(cids, replyItem.ReplyId)
	}
	cidMap, err := util.GetCommentFavorCids(userID, cids)
	if err != nil {
		return
	}
	//4.获取评论的动态信息.
	videoId, _ := GetCommentVid(req.GetCommentId())
	dynamicMap, err := util.GetCommentDynamicInfoWithVideoId(cids, videoId)
	if err != nil {
		return
	}
	for _, replyItem := range resp.ReplyItems {
		if _, ok := cidMap[replyItem.ReplyId]; ok {
			replyItem.Liked = true
		}
		if item, ok := dynamicMap[replyItem.ReplyId]; ok {
			replyItem.LikeCount = item.FavorNum
		}
	}
	if _, ok := cidMap[req.GetCommentId()]; ok {
		resp.CommentItem.Liked = true
	}
	if item, ok := dynamicMap[req.GetCommentId()]; ok {
		resp.CommentItem.ReplyCount = item.ReplyNum
		resp.CommentItem.LikeCount = item.FavorNum
	}

	resp.Status.Code = pb.Status_OK
	return
}

//评论视频
func (s *Server) CommentVideo(ctx context.Context, req *pb.CommentVideoRequest) (resp *pb.CommentVideoResponse, err error) {
	clientIp, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIp:%s, req:%+v", clientIp, req)
	//1.parse req
	videoID, _ := strconv.ParseUint(req.GetVideoId(), 10, 64)
	content := req.GetContent()
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	token := req.GetHeader().GetToken()

	resp = &pb.CommentVideoResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	var valid bool
	if valid, err = util.CheckIDValid(userID, token); !valid {
		resp.Status.Code = pb.Status_USER_NOTLOGIN
		glog.Error("userid:%d, token:%s check not valid. err:%v", userID, token, err)
		return
	}

	contentString, err := db.MysqlEscapeString(content)
	if err != nil {
		resp.Status.Code = pb.Status_BAD_REQUEST
		glog.Error("EscapeStrin failed. content:%s err:%v", content, err)
		return
	}

	//2.insert comment to mysql
	commentID, autoIncreID, err := s.idGenerator.GetItemId("comment_")
	if err != nil {
		glog.Error("generate commentid err:", err)
		return
	}
	userItem, err := getUserItem(userID)
	commentTableName, _ := db.GetTableName("comment_", commentID)

	exeSql := fmt.Sprintf("update %s set cid=%d, vid=%d, from_uid=%d, from_name='%s', from_photo='%s', content='%s' where id=%d", commentTableName, commentID, videoID, userID, userItem.Name, userItem.PhotoUrl, contentString, autoIncreID)
	_, err = db.Exec(common.BUDAODB, exeSql)
	if err != nil {
		glog.Error("insert comment failed. exeSql:%s, err:%v", exeSql, err)
		return
	}

	//3.audit comment
	if err = AuditComment(commentID, contentString); err != nil {
		glog.Error("audit comment failed. err:%v", err)
		return
	}

	//4.update video comment num in mysql and in redis
	if err = UpdateVideoCommentNum(videoID); err != nil {
		return
	}

	resp.CommentItem = &pb.CommentItem{
		CommentId:   strconv.FormatUint(commentID, 10),
		VideoId:     req.GetVideoId(),
		UserItem:    userItem,
		Content:     content,
		CommentTime: uint64(time.Now().Unix()),
	}

	//5.add comment to redis
	commentItem, _ := proto.Marshal(resp.CommentItem)
	if err = AddVideoCommentToRedis(videoID, commentID, string(commentItem)); err != nil {
		glog.Error("add comment to redis failed. err:%v", err)
		return
	}

	//6.update user_0 comment_num
	if err = UpdateUserCommentNum(userID, 1); err != nil {
		return
	}

	//7.update vid_dynamic
	_, err = db.SAdd(common.VID_DYNAMIC, req.GetVideoId())
	if err != nil {
		glog.Error("insert vid_dynamic set failed. err:%v", err)
		return
	}

	//TODO
	//record TRACE log

	resp.Status.Code = pb.Status_OK
	return
}

//回复(评论/回复)
func (s *Server) ReplyComment(ctx context.Context, req *pb.ReplyCommentRequest) (resp *pb.ReplyCommentResponse, err error) {
	clientIp, _ := twirp.RequestIp(ctx)
	glog.Debug("clientIp:%s, req:%+v", clientIp, req)
	//1.parse req
	content := req.GetContent()
	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	token := req.GetHeader().GetToken()

	resp = &pb.ReplyCommentResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	var valid bool
	if valid, err = util.CheckIDValid(userID, token); !valid {
		resp.Status.Code = pb.Status_USER_NOTLOGIN
		glog.Error("userid:%d, token:%s check not valid. err:%v", userID, token, err)
		return
	}

	contentString, err := db.MysqlEscapeString(content)
	if err != nil {
		glog.Error("EscapeStrin failed. content:%s err:%v", content, err)
		resp.Status.Code = pb.Status_BAD_REQUEST
		return
	}

	//2.get parentCommentId, videoId
	var commentID uint64
	var parentCommentID, videoID uint64
	isReplyComment := true
	switch req.Type.(type) {
	case *pb.ReplyCommentRequest_CommentId:
		commentID, _ = strconv.ParseUint(req.GetCommentId(), 10, 64)
	case *pb.ReplyCommentRequest_ReplyId:
		commentID, _ = strconv.ParseUint(req.GetReplyId(), 10, 64)
		isReplyComment = false
	}
	tableName, err := db.GetTableName("comment_", commentID)
	querySql := fmt.Sprintf("select vid, parentcomid from %s where cid=%d", tableName, commentID)
	var rows *sql.Rows
	rows, err = db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&videoID, &parentCommentID)
		if err != nil {
			glog.Error("scan failed. err:%v", err)
			return
		}
	}
	if isReplyComment {
		parentCommentID = commentID
	}

	//3. insert comment
	newCommentID, autoIncreID, err := s.idGenerator.GetItemId("comment_")
	if err != nil {
		glog.Error("generate commentid err:", err)
		return
	}
	userItem, err := getUserItem(userID)
	tableName, _ = db.GetTableName("comment_", newCommentID)
	sqlString := fmt.Sprintf("update %s set cid=%d, vid=%d, from_uid=%d, from_name='%s', from_photo='%s', to_comment_id=%d, parentcomid=%d, content='%s' where id=%d", tableName, newCommentID, videoID, userID, userItem.Name, userItem.PhotoUrl, commentID, parentCommentID, contentString, autoIncreID)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("insert comment err:", err)
		return
	}

	//4.审核评论
	if err = AuditComment(newCommentID, contentString); err != nil {
		glog.Error("audit comment failed. err:%v", err)
		return
	}

	//5.1 update parentcomment reply num in mysql and redis
	_, err = UpdateCommentReplyNum(parentCommentID, videoID)
	if err != nil {
		glog.Error("update comment reply num failed. err:%v", err)
	}
	//5.2 update vid comment num in mysql
	if err = UpdateVideoCommentNum(videoID); err != nil {
		return
	}
	//5.3 update user comment num in mysql
	if err = UpdateUserCommentNum(userID, 1); err != nil {
		return
	}

	//6.add reply comment to redis
	resp.ReplyItem = &pb.ReplyItem{
		ReplyId:   strconv.FormatUint(newCommentID, 10),
		UserItem:  userItem,
		Content:   contentString,
		ReplyTime: uint64(time.Now().Unix()),
	}
	if !isReplyComment {
		replyItems, err := GetMultiReplyFromRedis([]string{strconv.FormatUint(commentID, 10)})
		if err == nil {
			resp.ReplyItem.Target = &pb.ReplyItem_ReplyItem{
				ReplyItem: replyItems[0],
			}
		} else {
			glog.Error("not find reply item. replyid:%d err:%v", commentID, err)
		}
	}

	replyItemStr, _ := proto.Marshal(resp.ReplyItem)
	if err = AddCommentReplyToRedis(parentCommentID, newCommentID, string(replyItemStr)); err != nil {
		glog.Error("add comment reply to redis failed. err:%v", err)
		return
	}

	//7.update vid_dynamic
	_, err = db.SAdd(common.VID_DYNAMIC, strconv.FormatUint(videoID, 10))
	if err != nil {
		glog.Error("insert vid_dynamic set failed. err:%v", err)
		return
	}

	resp.Status.Code = pb.Status_OK

	//TODO
	//record TRACE log

	return
}

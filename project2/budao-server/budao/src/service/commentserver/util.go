package commentserver

import (
	"common"
	"context"
	"db"
	"errors"
	"fmt"
	"service/transfer"
	"service/util"
	"strconv"
	"strings"
	"time"
	pb "twirprpc"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
)

const (
	USER_ITEM_PREFIX = "user_item_"
	USER_ITEM_EXPIRE = 24 * 60 * 60
)

// UserInfo define the user information from mysql
type UserInfo struct {
	ID              int64
	Name            string
	Photo           string
	VideoFavorNUM   int
	VideoShareNUM   int
	ContentNUM      int
	ContentFavorNUM int
	FollowTopicNUM  int
}

// GetUserInfoFromMySQLByUID function get user information.
func GetUserInfoFromMySQLByUID(userID uint64) (userInfo *UserInfo, err error) {
	tablename, err := db.GetTableName("user_", userID)
	sqlString := fmt.Sprintf("select uid, name, photo, video_favor_num, video_share_num, comment_num, comment_favor_num, follow_topic_num from %s where uid = %v", tablename, userID)
	rows, err := db.Query(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("query %s failed. sqlString:%s, err:%v", tablename, sqlString, err)
		return nil, err
	}
	defer rows.Close()

	user := UserInfo{}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Photo, &user.VideoFavorNUM, &user.VideoShareNUM, &user.ContentNUM, &user.ContentFavorNUM, &user.FollowTopicNUM)
	}

	return &UserInfo{
		ID:              user.ID,
		Name:            user.Name,
		Photo:           user.Photo,
		VideoFavorNUM:   user.VideoFavorNUM,
		VideoShareNUM:   user.VideoShareNUM,
		ContentNUM:      user.ContentNUM,
		ContentFavorNUM: user.ContentFavorNUM,
		FollowTopicNUM:  user.FollowTopicNUM,
	}, nil
}

func getUserItem(userID uint64) (userItem *pb.UserItem, err error) {
	//1.get from redis. hash, key:global_user_item field:userid value:pb
	userItemKey := fmt.Sprintf("%s%d", USER_ITEM_PREFIX, userID)
	userItem = &pb.UserItem{}
	value, err := db.GetString(userItemKey)
	if err != nil {
		if strings.Contains(err.Error(), "redigo: nil returned") {
			glog.Debug("user:%d not in redis.", userID)
		} else {
			glog.Error("get string failed. err:%v", err)
			return
		}
	}

	if value != "" {
		if err = proto.Unmarshal([]byte(value), userItem); err == nil {
			return
		}
		glog.Error("unmarshal failed. err:%v", err)
	}

	//2.not in redis, get from mysql
	item, err := GetUserInfoFromMySQLByUID(userID)
	if err != nil {
		glog.Error("query err:%v", err)
		return
	}

	userItem = &pb.UserItem{
		UserId:   strconv.FormatInt(item.ID, 10),
		Name:     item.Name,
		PhotoUrl: item.Photo,
	}

	//3.add to redis
	userItemStr, err := proto.Marshal(userItem)
	if err != nil {
		glog.Error("proto marshal failed. err:%v", err)
		return
	}
	if err = db.SetString(userItemKey, string(userItemStr)); err != nil {
		glog.Error("set useritem failed. err:%v", err)
		return
	}
	_, _ = db.KExpire(userItemKey, USER_ITEM_EXPIRE)
	return
}

func AuditComment(commentId uint64, content string) (err error) {
	req := &pb.AuditContentRequest{
		ContentId:   commentId,
		Content:     content,
		ContentType: "comment",
	}
	resp, err := transfer.GetTransfer().AuditContent(context.Background(), req)
	if err != nil {
		glog.Error("audit comment failed. req:%v,resp:%v,err:%v", req, resp, err)
		return
	}
	if resp.Result != pb.EPostResult_E_PostResult_OK {
		glog.Error("post audit content failed.req:%v,resp:%v", req, resp)
		err = errors.New("post audit content failed.")
		return
	}
	if resp.AuditResult != pb.EAuditResult_E_AuditResult_Pass {
		err = errors.New("audit content not pass.")
		glog.Error("audit content failed.req:%v,resp:%v", req, resp)
		return
	}
	glog.Debug("post comment success. req:%v,resp:%v", req, resp)
	return
}

func AddCommentReplyToRedis(commentId, replyId uint64, replyItem string) (err error) {
	//添加顺序集合
	zsetKey := fmt.Sprintf("%s%d", util.CRNEW_PREFIX, commentId)
	timeCids := make([]interface{}, 0)
	timeCids = append(timeCids, time.Now().Unix(), replyId)
	if r, err := db.ZAddMulti(zsetKey, timeCids); err != nil {
		glog.Error("zadd multi failed. err:%v, r:%v", err, r)
	}
	idItems := make([]interface{}, 0)
	idItems = append(idItems, replyId, replyItem)
	if r, err := db.HMSet(util.COMMENT_ITEM_ALL_KEY, idItems); err != nil && r != "ok" {
		glog.Error("hmset failed. err:%v, r:%v", err, r)
	}
	return
}

func AddVideoCommentToRedis(videoId, commentId uint64, commentItem string) (err error) {
	//1.add video comment new
	vcNewKey := fmt.Sprintf("%s%d", util.VCNEW_PREFIX, videoId)
	newCids := []interface{}{time.Now().Unix(), commentId}
	//newCids = append(newCids, time.Now().Unix(), commentId)
	if r, err := db.ZAddMulti(vcNewKey, newCids); err != nil {
		glog.Error("zadd multi failed. err:%v, r:%v", err, r)
	}
	//2.add video comment hot
	vcHotKey := fmt.Sprintf("%s%d", util.VCHOT_PREFIX, videoId)
	hotCids := []interface{}{0, commentId}
	//hotCids = append(timeCids, 0, commentId)
	if r, err := db.ZAddMulti(vcHotKey, hotCids); err != nil {
		glog.Error("zadd multi failed. err:%v, r:%v", err, r)
	}
	//3.add comment item
	idItems := []interface{}{commentId, commentItem}
	//idItems = append(idItems, commentId, commentItem)
	if r, err := db.HMSet(util.COMMENT_ITEM_ALL_KEY, idItems); err != nil && r != "ok" {
		glog.Error("hmset failed. err:%v, r:%v", err, r)
	}
	return
}

func UpdateVideoCommentNum(videoId uint64) (err error) {
	//1. update mysql
	videoTableName, err := db.GetTableName("video_", videoId)
	exeSql := fmt.Sprintf("update %s set comment_num=comment_num+1 where vid=%d", videoTableName, videoId)
	_, err = db.Exec(common.BUDAODB, exeSql)
	if err != nil {
		glog.Error("update video comment num failed. exeSql:%s. err:%v", exeSql, err)
		return
	}

	//2. update redis
	field := fmt.Sprintf("%s_%d", common.COMMENT_COUNT, videoId)
	if _, err = db.HIncrBy(common.VIDEO_DYNAMIC, field, 1); err != nil {
		glog.Error("hincrby video comment num failed. err:%s", err)
		return
	}
	return
}

func GetMultiReplyFromRedis(replyIds []string) (replyItems []*pb.ReplyItem, err error) {
	iarr := make([]interface{}, len(replyIds))
	for i, id := range replyIds {
		iarr[i] = interface{}(id)
	}
	strArr, err := db.HMGet(util.COMMENT_ITEM_ALL_KEY, iarr)
	if err != nil {
		return
	}
	replyItems = make([]*pb.ReplyItem, 0)
	for _, str := range strArr {
		replyItem := &pb.ReplyItem{}
		err = proto.Unmarshal([]byte(str), replyItem)
		if err != nil {
			glog.Error("proto unmarshal failed. str:%s, err:%v", str, err)
			continue
		}
		replyItems = append(replyItems, replyItem)
	}
	return
}

func GetCommentIdsByLastCommentId(key, lastCommentId string, returnNum int) (commentIds []string, err error) {
	//1.get lastCommentId rank in zset
	offset := 0
	if lastCommentId != "" {
		offset, err = db.ZRevRank(key, lastCommentId)
		if err != nil {
			return
		}
		offset += 1
	}
	commentIds, err = db.ZRevRange(key, int64(offset), int64(offset+returnNum))

	return
}

func GetCommentItemsByCommentIds(userID uint64, commentIds []string, videoID uint64, explicitReplyNum int64) (commentItems []*pb.CommentItem, err error) {
	commentItems = make([]*pb.CommentItem, 0)
	if len(commentIds) == 0 {
		glog.Error("userId:%d get commentItems failed. commentIds is empty")
		return
	}
	//1.get commentItem from redis
	commentItems, err = util.GetMultiCommentFromRedis(commentIds)
	if err != nil {
		glog.Error("get comment item from redis failed. err:%v", err)
		return
	}
	//fill comment dynamic info
	cids := make([]string, 0)
	num := 0
	for _, commentItem := range commentItems {
		cids = append(cids, commentItem.CommentId)
		num++
	}

	dynamicMap, err := util.GetCommentDynamicInfoWithVideoId(cids, videoID)
	if err != nil {
		return
	}
	for _, commentItem := range commentItems {
		if item, ok := dynamicMap[commentItem.CommentId]; ok {
			commentItem.ReplyCount = item.ReplyNum
			commentItem.LikeCount = item.FavorNum
		}
	}

	//2.get reply of comment
	for _, commentItem := range commentItems {
		if commentItem.ReplyCount == 0 {
			continue
		}
		crnewKey := fmt.Sprintf("%s%s", util.CRNEW_PREFIX, commentItem.CommentId)
		replyIds, err := db.ZRevRange(crnewKey, 0, explicitReplyNum-1)
		if err != nil {
			glog.Error("get reply id from redis failed. err:%v", err)
			continue
		}
		if len(replyIds) == 0 {
			glog.Debug("no reply real. cid:%s", commentItem.CommentId)
			continue
		}
		replyItems, err := GetMultiReplyFromRedis(replyIds)
		if err != nil {
			glog.Error("get reply item from redis failed. err:%v", err)
			continue
		}
		commentItem.ReplyItems = append(commentItem.ReplyItems, replyItems...)
		for _, replyItem := range commentItem.ReplyItems {
			cids = append(cids, replyItem.ReplyId)
		}
	}

	//fill reply dynamic info
	if len(cids) != num {
		dynamicMap, err = util.GetCommentDynamicInfoWithVideoId(cids[num:], videoID)
		if err != nil {
			return
		}
		for _, commentItem := range commentItems {
			if commentItem.ReplyCount == 0 {
				continue
			}
			for _, replyItem := range commentItem.ReplyItems {
				if item, ok := dynamicMap[replyItem.ReplyId]; ok {
					replyItem.LikeCount = item.FavorNum
				}
			}
		}
	}

	//3.从mysql中读取当前用户是否对这些评论点过赞
	for {
		var cidMap map[string]bool
		cidMap, err = util.GetCommentFavorCids(userID, cids)
		if err != nil {
			return
		}
		if len(cidMap) == 0 {
			break
		}
		for _, commentItem := range commentItems {
			if _, ok := cidMap[commentItem.CommentId]; ok {
				commentItem.Liked = true
			}
			for _, replyItem := range commentItem.ReplyItems {
				if _, ok := cidMap[replyItem.ReplyId]; ok {
					replyItem.Liked = true
				}
			}
		}
		break
	}

	return
}

func GetCommentVid(commentIdStr string) (vidNum uint64, err error) {
	commentID, _ := strconv.ParseUint(commentIdStr, 10, 64)
	tableName, _ := db.GetTableName("comment_", commentID)
	querySql := fmt.Sprintf("select vid from %s where cid=%s", tableName, commentIdStr)
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
		return
	}
	defer rows.Close()
	var vid string
	for rows.Next() {
		err = rows.Scan(&vid)
		if err != nil {
			glog.Error("scan failed. err:%v", err)
			return
		}
	}
	vidNum, _ = strconv.ParseUint(vid, 10, 64)

	return
}

func UpdateCommentReplyNum(cid, vid uint64) (ret int, err error) {
	//update mysql
	tableName, _ := db.GetTableName("comment_", cid)
	exeSql := fmt.Sprintf("update %s set reply_num=reply_num+1 where cid=%d", tableName, cid)
	if _, err = db.Exec(common.BUDAODB, exeSql); err != nil {
		glog.Error("update comment reply num failed. exeSql:%s. err:%v", exeSql, err)
		return
	}

	//update redis
	key := fmt.Sprintf("%s%d", util.COMMENT_DYNAMIC_PREFIX, vid%common.COMMENT_DYNAMIC_KEY_NUM)
	field := fmt.Sprintf("%s%d", util.REPLY_NUM_PREFIX, cid)
	ret, err = db.HIncrBy(key, field, 1)
	if err != nil {
		glog.Error("hincrby failed. err:%v", err)
	}
	return
}

func UpdateCommentFavorNum(cid, vid uint64, increment int) (ret int, err error) {
	//update mysql
	tableName, _ := db.GetTableName("comment_", cid)
	exeSql := fmt.Sprintf("update %s set favor_num=favor_num+%d where cid=%d", tableName, increment, cid)
	if _, err = db.Exec(common.BUDAODB, exeSql); err != nil {
		glog.Error("update comment num favor failed. exeSql:%s. err:%v", exeSql, err)
		return
	}

	//update redis
	key := fmt.Sprintf("%s%d", util.COMMENT_DYNAMIC_PREFIX, vid%common.COMMENT_DYNAMIC_KEY_NUM)
	field := fmt.Sprintf("%s%d", util.FAVOR_NUM_PREFIX, cid)
	ret, err = db.HIncrBy(key, field, increment)
	if err != nil {
		glog.Error("hincrby failed. err:%v", err)
	}
	return
}

func GetCommentDynamicInfo(commentIdStr string) (dynamicInfo *util.CommentDynamicInfo, err error) {
	vid, _ := GetCommentVid(commentIdStr)
	dynamicInfoMap, _ := util.GetCommentDynamicInfoWithVideoId([]string{commentIdStr}, vid)
	dynamicInfo = dynamicInfoMap[commentIdStr]
	return
}

func UpdateCommentRank(cid string, vid uint64) (err error) {
	//1 read dynamic info
	dynamicMap, err := util.GetCommentDynamicInfoWithVideoId([]string{cid}, vid)
	if err != nil {
		glog.Error("get dynamic info failed. cid:%s, vid:%d, err:%v", cid, vid, err)
		return
	}
	if len(dynamicMap) == 0 {
		glog.Error("not find dynamic info. cid:%s, vid:%d, err:%v", cid, vid, err)
		return
	}
	//2 update score
	dynamic := dynamicMap[cid]
	newScore := dynamic.Weight + dynamic.FavorNum + dynamic.ReplyNum
	vcHotKey := fmt.Sprintf("%s%d", util.VCHOT_PREFIX, vid)
	hotCids := make([]interface{}, 0)
	hotCids = append(hotCids, newScore, cid)
	if _, err = db.ZAddMulti(vcHotKey, hotCids); err != nil {
		glog.Error("zadd multi failed. vcHotKey:%s, hotCids:%v, err:%v", vcHotKey, hotCids, err)
		return
	}
	return
}

func UpdateUserCommentNum(userId uint64, increment int) (err error) {
	tableName, err := db.GetTableName("user_", userId)
	exeSql := fmt.Sprintf("update %s set comment_num=comment_num+%d where uid=%d", tableName, increment, userId)
	_, err = db.Exec(common.BUDAODB, exeSql)
	if err != nil {
		glog.Error("update user comment num failed. exeSql:%s. err:%v", exeSql, err)
		return
	}
	return
}

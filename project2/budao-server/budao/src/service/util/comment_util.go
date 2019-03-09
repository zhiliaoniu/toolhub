package util

import (
	"common"
	"db"
	"errors"
	"fmt"
	"strings"
	pb "twirprpc"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
)

const (
	USER_OP_COMMENT_FAVOR = iota
	USER_OP_COMMENT_SHARED
	USER_OP_COMMENT_REPLY

	COMMENT_ITEM_ALL_KEY   = "comment_item_all"
	VCNEW_PREFIX           = "vcnew_"
	VCHOT_PREFIX           = "vchot_"
	CRNEW_PREFIX           = "crnew_"
	COMMENT_DYNAMIC_PREFIX = "comment_dynamic_"
	FAVOR_NUM_PREFIX       = "favor_num_"
	REPLY_NUM_PREFIX       = "reply_num_"
	WEIGHT_PREFIX          = "weight_"
)

type CommentDynamicInfo struct {
	IsReply  bool
	FavorNum uint32
	ReplyNum uint32
	Weight   uint32
}

func GetMultiCommentFromRedis(commentIds []string) (commentItems []*pb.CommentItem, err error) {
	iarr := make([]interface{}, len(commentIds))
	for i, id := range commentIds {
		iarr[i] = interface{}(id)
	}
	strArr, err := db.HMGet(COMMENT_ITEM_ALL_KEY, iarr)
	if err != nil {
		glog.Error("hmget failed. key:%s, fields:%v", COMMENT_ITEM_ALL_KEY, iarr)
		return
	}
	commentItems = make([]*pb.CommentItem, 0)
	for _, str := range strArr {
		commentItem := &pb.CommentItem{}
		err = proto.Unmarshal([]byte(str), commentItem)
		if err != nil {
			continue
		}
		commentItems = append(commentItems, commentItem)
	}

	return
}

func GetCommentFavorCids(userID uint64, commentIds []string) (cidMap map[string]bool, err error) {
	if len(commentIds) == 0 {
		err = errors.New("empty cids")
		return
	}
	cidsStr := strings.Join(commentIds, ",")
	tableName, err := db.GetTableName("user_favor_comment_", userID)
	querySql := fmt.Sprintf("select cid from %s where uid=%d and result=1 and cid in(%s)", tableName, userID, cidsStr)
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed.querySql:%s, err:%v", querySql, err)
		return
	}
	defer rows.Close()
	var cid string
	cidMap = make(map[string]bool, 0)
	for rows.Next() {
		err = rows.Scan(&cid)
		if err != nil {
			glog.Error("scan failed. err:", err)
			return
		}
		cidMap[cid] = true
	}

	return
}
func GetCommentDynamicInfoWithVideoId(commentIds []string, videoId uint64) (dynamicMap map[string]*CommentDynamicInfo, err error) {
	key := fmt.Sprintf("%s%d", COMMENT_DYNAMIC_PREFIX, videoId%common.COMMENT_DYNAMIC_KEY_NUM)
	fields := make([]interface{}, 0)
	var field string
	for _, cid := range commentIds {
		field = fmt.Sprintf("%s%s", FAVOR_NUM_PREFIX, cid)
		fields = append(fields, field)
		field = fmt.Sprintf("%s%s", REPLY_NUM_PREFIX, cid)
		fields = append(fields, field)
	}
	dynamicArr, err := db.HMGetInt(key, fields)
	if err != nil {
		glog.Error("hmget failed. err:%v", err)
		return
	}
	glog.Debug("key:%s, fields:%v, dynamicArr:%v", key, fields, dynamicArr)

	dynamicMap = make(map[string]*CommentDynamicInfo, 0)
	for i := 0; i < len(dynamicArr); i += 2 {
		dynamicMap[commentIds[i/2]] = &CommentDynamicInfo{
			FavorNum: uint32(dynamicArr[i]),
			ReplyNum: uint32(dynamicArr[i+1]),
		}
	}
	return
}

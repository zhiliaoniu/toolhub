package util

import (
	"common"
	"db"
	"fmt"
	"strconv"
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

func GetVidsHotComment(vids []string) (vidCommentMap map[string]*pb.CommentItem, err error) {
	vidCommentMap = make(map[string]*pb.CommentItem, 0)
	for _, vid := range vids {
		//1.get hot comment id
		vchotKey := fmt.Sprintf("%s%s", VCHOT_PREFIX, vid)
		cids, err := db.ZRevRange(vchotKey, 0, 0)
		if err != nil {
			glog.Debug("get vid:%s hot comment id from redis failed. err:%v", vid, err)
			continue
		}
		if len(cids) == 0 {
			glog.Debug("video %s has no hot comment.", vid)
			continue
		}

		//2.get hot comment
		commentItems, err := GetMultiCommentFromRedis(cids)
		if err != nil {
			glog.Error("get vid:%s comment item cid:%v from redis failed.err:%v", vid, cids, err)
			continue
		}
		if len(commentItems) == 0 {
			glog.Error("video %s get hot comment cids:%v failed.", vid, cids)
			continue
		}

		//3.fill comment dynamic info
		vidNum, _ := strconv.ParseUint(vid, 10, 64)
		dynamicMap, err := GetCommentDynamicInfoWithVideoId(cids, vidNum)
		if err == nil {
			for _, commentItem := range commentItems {
				if item, ok := dynamicMap[commentItem.CommentId]; ok {
					commentItem.ReplyCount = item.ReplyNum
					commentItem.LikeCount = item.FavorNum
				}
			}
		}
		vidCommentMap[vid] = commentItems[0]
	}

	return
}

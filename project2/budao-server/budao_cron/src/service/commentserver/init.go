package commentserver

import (
	"db"
	"fmt"
	"service/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"

	"common"
	pb "twirprpc"
)

var UpdateCommentWG sync.WaitGroup

func InitCommentInfoByTraverseMysql() {
	defer common.WG.Done()
	if isNeedInit() {
		UpdateCommentInfoByTraverseMysql()
	}
}

func UpdateCommentInfoByTraverseMysql() {
	glog.Debug("==================start update comment by traverse mysql")
	tableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc[common.VIDEO_TABLE_PREFIX]
	var num uint64
	deleteVids := make([]string, 0)
	for num = 0; num < tableNum; num++ {
		tableName := fmt.Sprintf("%s%d", common.VIDEO_TABLE_PREFIX, num)
		count := 0
		for {
			glog.Debug("read_vids from %d to %d", count, count+common.READ_MYSQL_MAX_ROWS)
			querySql := fmt.Sprintf("select vid, state, op_state from %s limit %d, %d", tableName, count, common.READ_MYSQL_MAX_ROWS)
			count += common.READ_MYSQL_MAX_ROWS
			rows, err := db.Query(common.BUDAODB, querySql)
			if err != nil {
				glog.Error("read mysql failed. querySql:%s, err:%v", querySql, err)
				break
			}
			defer rows.Close()

			var vid string
			var state, opState int
			vids := make([]string, 0)
			for rows.Next() {
				err = rows.Scan(&vid, &state, &opState)
				if err != nil {
					glog.Error("scan err:%v", err)
					break
				}
				if state != 2 || opState != 0 {
					deleteVids = append(deleteVids, vid)
					continue
				}
				vids = append(vids, vid)
			}
			go DeleteVideoAllCommentFromRedis(deleteVids)

			if len(vids) == 0 {
				break
			}
			//glog.Debug("vids:%v", vids)
			commentMap, replyMap, dynamicMap, err := GetCommentAndReplyByVid(vids)
			if err != nil {
				glog.Error("get comment and reply by vid failed. err:%v", err)
				return
			}
			UpdateCommentWG.Add(1)
			go UpdateCommentMultiInfo(vids, commentMap, replyMap, dynamicMap)
		}
	}
	UpdateCommentWG.Wait()
	glog.Debug("==================end update comment by traverse mysql")
}

func UpdateCommentMultiInfo(vids []string, commentMap map[string]map[string]*pb.CommentItem, replyMap map[string]map[string]*pb.ReplyItem, dynamicMap map[string]map[string]*util.CommentDynamicInfo) {
	defer UpdateCommentWG.Done()
	SaveCommentToRedis(commentMap)

	go SaveReplyToRedis(replyMap)

	SaveDynamicToRedis(dynamicMap)

	UpdateCommentHotRank(dynamicMap)

	go InitVidHotComment(vids)

}

func SaveDynamicToRedis(dynamicMap map[string]map[string]*util.CommentDynamicInfo) {
	for vid, vDynamicItem := range dynamicMap {
		//glog.Debug("vid:%s", vid)
		vidNum, _ := strconv.ParseUint(vid, 10, 64)
		key := fmt.Sprintf("%s%d", util.COMMENT_DYNAMIC_PREFIX, vidNum%common.COMMENT_DYNAMIC_KEY_NUM)
		fields := make([]interface{}, 0)
		for cid, item := range vDynamicItem {
			field := fmt.Sprintf("%s%s", util.FAVOR_NUM_PREFIX, cid)
			fields = append(fields, field, item.FavorNum)
			field = fmt.Sprintf("%s%s", util.WEIGHT_PREFIX, cid)
			fields = append(fields, field, item.Weight)
			if item.IsReply == false {
				field = fmt.Sprintf("%s%s", util.REPLY_NUM_PREFIX, cid)
				fields = append(fields, field, item.ReplyNum)
			}
		}
		if _, err := db.HMSet(key, fields); err != nil {
			glog.Error("hmset failed. err:%v", err)
			continue
		}
	}
}

func UpdateCommentHotRank(dynamicMap map[string]map[string]*util.CommentDynamicInfo) {
	for vid, vDynamicItem := range dynamicMap {
		hotCids := make([]interface{}, 0)
		for cid, item := range vDynamicItem {
			if item.IsReply == false {
				score := item.Weight + item.FavorNum + item.ReplyNum
				hotCids = append(hotCids, score, cid)
			}
		}
		if len(hotCids) == 0 {
			continue
		}
		vcHotKey := fmt.Sprintf("%s%s", util.VCHOT_PREFIX, vid)
		if r, err := db.ZAddMulti(vcHotKey, hotCids); err != nil {
			glog.Error("zadd multi failed. vcHotKey:%s, hotCids:%v, err:%v, r:%v", vcHotKey, hotCids, err, r)
			continue
		}
	}
}

func GetCommentAndReplyByVid(vids []string) (commentMap map[string]map[string]*pb.CommentItem, replyMap map[string]map[string]*pb.ReplyItem, dynamicMap map[string]map[string]*util.CommentDynamicInfo, err error) {
	//1.compose query sql
	commentTableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["comment_"]

	vidstr := strings.Join(vids, ",")
	commentMap = make(map[string]map[string]*pb.CommentItem, 0)          //<vid, <cid, c>>
	replyMap = make(map[string]map[string]*pb.ReplyItem, 0)              //<cid, <rid, r>>
	dynamicMap = make(map[string]map[string]*util.CommentDynamicInfo, 0) //<vid, <cid, CommentDynamicInfo>>
	ridMap := make(map[string]*pb.ReplyItem, 0)                          //<rid, r>
	rridMap := make(map[string]string, 0)                                //<rid, rid>

	//2.traverse comment table
	var num uint64 = 0
	for {
		if num >= commentTableNum {
			break
		}
		tableName := fmt.Sprintf("%s%d", "comment_", num)
		num++
		querySql := fmt.Sprintf("select cid, vid, from_uid, from_name, from_photo, to_comment_id, parentcomid, content, favor_num, fake_favor_num, reply_num, weight, state, create_time from %s where vid in (%s)", tableName, vidstr)
		rows, err := db.Query(common.BUDAODB, querySql)
		if err != nil {
			glog.Error("read mysql failed. querySql:%s err:%v", querySql, err)
			continue
		}
		defer rows.Close()

		for rows.Next() {
			item := &pb.CommentItem{}
			user := &pb.UserItem{}
			var ToCommentId, ParentComId uint64
			var CreateTime time.Time
			var LikeCount, FakeLikeCount, ReplyCount, Weight, State uint32
			err = rows.Scan(&item.CommentId, &item.VideoId, &user.UserId, &user.Name, &user.PhotoUrl, &ToCommentId, &ParentComId, &item.Content, &LikeCount, &FakeLikeCount, &ReplyCount, &Weight, &State, &CreateTime)
			if err != nil {
				glog.Error("query mysql failed. err:%v", err)
				break
			}
			if State != 2 {
				DeleteVideoOneCommentFromRedis(item.CommentId, item.VideoId, ParentComId)
				continue
			}
			item.UserItem = user
			item.CommentTime = common.GetLocalUnix(CreateTime)
			//glog.Debug("query mysql comment item:%v", item)
			//glog.Debug("CommentId:%s VideoId:%s ParentComId:%d ToCommentId:%d", item.CommentId, item.VideoId, ParentComId, ToCommentId)

			//save dynamic info
			vDynamicItem, ok := dynamicMap[item.VideoId]
			if !ok {
				vDynamicItem = make(map[string]*util.CommentDynamicInfo, 0)
				dynamicMap[item.VideoId] = vDynamicItem
			}
			vDynamicItem[item.CommentId] = &util.CommentDynamicInfo{true, LikeCount + FakeLikeCount, ReplyCount, Weight}
			if ParentComId == 0 {
				vDynamicItem[item.CommentId].IsReply = false
			}

			//save statis info
			if ParentComId == 0 {
				vCommentItem, ok := commentMap[item.VideoId]
				if !ok {
					vCommentItem = make(map[string]*pb.CommentItem, 0)
					commentMap[item.VideoId] = vCommentItem
				}
				vCommentItem[item.CommentId] = item
			} else {
				//compose reply
				replyItem := &pb.ReplyItem{
					ReplyId:   item.CommentId,
					UserItem:  item.UserItem,
					Content:   item.Content,
					ReplyTime: item.CommentTime,
				}
				ridMap[replyItem.ReplyId] = replyItem
				parentComIdStr := strconv.FormatUint(ParentComId, 10)
				cReplyItem, ok := replyMap[parentComIdStr]
				if !ok {
					cReplyItem = make(map[string]*pb.ReplyItem, 0)
					replyMap[parentComIdStr] = cReplyItem
				}
				cReplyItem[replyItem.ReplyId] = replyItem

				//reply to reply
				if ParentComId != ToCommentId {
					ToCommentIdStr := strconv.FormatUint(ToCommentId, 10)
					rridMap[replyItem.ReplyId] = ToCommentIdStr
				}
			}
		}
	}
	//compose reply's reply
	for rid, torId := range rridMap {
		rItem, _ := ridMap[rid]
		torItem, ok := ridMap[torId]
		if !ok {
			glog.Error("not find toreplyItem. replyId:%s toReplyId:%s", rid, torId)
			continue
		}
		rItem.Target = &pb.ReplyItem_ReplyItem{
			ReplyItem: torItem,
		}
	}

	return
}

func SaveReplyToRedis(replyMap map[string]map[string]*pb.ReplyItem) {
	//glog.Debug("SaveReplyToRedis---------------replyMap:%+v", replyMap)
	replyNum := 0
	for cid, cReplyItems := range replyMap {
		timeRids := make([]interface{}, 0)
		idItems := make([]interface{}, 0)
		//compose redis cmd
		//glog.Debug("-----------cid:%s-------------", cid)
		for rid, item := range cReplyItems {
			timeRids = append(timeRids, item.ReplyTime, rid)
			itemstr, err := proto.Marshal(item)
			if err != nil {
				glog.Error("err:%v", err)
				continue
			}
			idItems = append(idItems, item.ReplyId, itemstr)
			//glog.Debug("rid:%s", rid)
		}
		replyNum += len(cReplyItems)
		zsetKey := fmt.Sprintf("%s%s", util.CRNEW_PREFIX, cid)
		if r, err := db.ZAddMulti(zsetKey, timeRids); err != nil {
			glog.Error("zadd multi failed. err:%v, r:%v", err, r)
			continue
		}
		//glog.Debug("timeRids:%v timeRids.len:%d r:%d", timeRids, len(timeRids)/2, r)
		if r, err := db.HMSet(util.COMMENT_ITEM_ALL_KEY, idItems); err != nil && r != "ok" {
			glog.Error("hmset failed. err:%v, r:%v", err, r)
			continue
		}
		//glog.Debug("idItems.len:%d ", len(idItems)/2)
	}
	//glog.Debug("replyNum:%d", replyNum)
}

func SaveCommentToRedis(commentMap map[string]map[string]*pb.CommentItem) {
	//glog.Debug("SaveCommentToRedis---------------commentMap:%+v", commentMap)
	commentNum := 0
	for vid, vCommentItems := range commentMap {
		timeCids := make([]interface{}, 0)
		idItems := make([]interface{}, 0)
		//compose redis cmd
		//glog.Debug("-----------vid:%s-------------", vid)
		for cid, item := range vCommentItems {
			timeCids = append(timeCids, item.CommentTime, cid)
			itemstr, err := proto.Marshal(item)
			if err != nil {
				glog.Error("err:%v", err)
				continue
			}
			idItems = append(idItems, item.CommentId, itemstr)
			//glog.Debug("cid:%s", cid)
		}
		commentNum += len(vCommentItems)
		vcNewKey := fmt.Sprintf("%s%s", util.VCNEW_PREFIX, vid)
		if r, err := db.ZAddMulti(vcNewKey, timeCids); err != nil {
			glog.Error("zadd multi failed. vcNewKey:%s, timeCids:%v, err:%v, r:%v", vcNewKey, timeCids, err, r)
			continue
		}
		//glog.Debug("timeCids:%v timeCids.len:%d r:%d", timeCids, len(timeCids)/2, r)
		if r, err := db.HMSet(util.COMMENT_ITEM_ALL_KEY, idItems); err != nil && r != "ok" {
			glog.Error("hmset failed. err:%v, r:%v", err, r)
			continue
		}
		//glog.Debug("idItems.len:%d ", len(idItems)/2)
	}
	//glog.Debug("commentNum:%d", commentNum)
}

func isNeedInit() bool {
	//1.assume video_0 has need data
	querySql := fmt.Sprintf("select vid from video_0 where comment_num > 0 order by create_time asc limit 1")
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. err:", err)
	}
	defer rows.Close()

	var vid uint64
	for rows.Next() {
		err = rows.Scan(&vid)
		if err != nil {
			glog.Error("query mysql failed. err:%v", err)
			return true
		}
	}

	//2.judge whether key exist in redis
	key := fmt.Sprintf("%s%d", util.VCNEW_PREFIX, vid)
	var value int
	if value, err = db.IsKeyExist(key); err != nil {
		glog.Error("query redis failed. err:%v", err)
		return true
	}
	if value == 1 {
		glog.Debug("vid:%d comment %s has been in redis", vid, key)
		return false
	}

	glog.Debug("vid:%d comment %s not in redis", vid, key)
	return true
}

//
func InitVidHotComment(vids []string) {
	//1 读取最新热评信息
	vidHotCommentMap, err := util.GetVidsHotComment(vids)
	if err != nil {
		glog.Error("get video hot comment failed. vids:%v err:%v", vids, err)
		return
	}
	if len(vidHotCommentMap) == 0 {
		glog.Debug("initvids:%s has no hot comment")
		return
	}

	//2.save video hot comment
	hotCommentArr := make([]interface{}, 0)
	for vid, commentItem := range vidHotCommentMap {
		artictlItem := &pb.ArticleItem{
			CommentItems: []*pb.CommentItem{commentItem},
		}
		artictlItemStr, err := proto.Marshal(artictlItem)
		if err != nil {
			glog.Error("vid:%s marshal hot comment failed. err:%v", vid, err)
			continue
		}
		hotCommentArr = append(hotCommentArr, vid, artictlItemStr)
	}
	if len(hotCommentArr) == 0 {
		glog.Error("vids:%s has no hot comment")
		return
	}
	_, err = db.HMSet(common.VIDEO_HOT_COMMENT, hotCommentArr)
	if err != nil {
		glog.Error("hmset save vid hot comment failed. err:%v", err)
		return
	}
	glog.Debug("save vid:%s hot comment success. hotvid.len:%d", vids, len(hotCommentArr)/2)

	return
}

func DeleteVideoAllCommentFromRedis(vids []string) {
	if len(vids) == 0 {
		return
	}

	//1.delete video hot comment in articleItem
	vidsInter := common.TransStrArrToInterface(vids)
	if _, err := db.HMDelete(common.VIDEO_HOT_COMMENT, vidsInter); err != nil {
		glog.Error("hmdelete failed. err:%v", err)
		return
	}

	//2.delete video comment rank
	keys := make([]interface{}, 0)
	var key1, key2 string
	for _, vid := range vids {
		key1 = fmt.Sprintf("%s%s", util.VCHOT_PREFIX, vid)
		key2 = fmt.Sprintf("%s%s", util.VCNEW_PREFIX, vid)
		keys = append(keys, key1, key2)
	}
	if _, err := db.DelMultiKey(keys); err != nil {
		glog.Error("del failed. keys:%v, err:%v", keys, err)
		return
	}
}

func DeleteVideoOneCommentFromRedis(commentId, videoId string, parentComId uint64) {
	exist, _ := db.HExists(util.COMMENT_ITEM_ALL_KEY, commentId)
	if exist == false {
		return
	}

	glog.Debug("delete commentId:%s from redis", commentId)
	//1.delete comment item
	db.HDelete(util.COMMENT_ITEM_ALL_KEY, commentId)

	//2.delete rank
	if parentComId != 0 {
		//delete reply
		key := fmt.Sprintf("%s%d", util.CRNEW_PREFIX, parentComId)
		_, err := db.ZDelete(key, commentId)
		if err != nil {
			glog.Error("zdelete failed. err:%v", err)
		}
	} else {
		//delte comment hot rank
		key := fmt.Sprintf("%s%s", util.VCHOT_PREFIX, videoId)
		_, err := db.ZDelete(key, commentId)
		if err != nil {
			glog.Error("zdelete failed. err:%v", err)
		}
		//delte comment new rank
		key = fmt.Sprintf("%s%s", util.VCNEW_PREFIX, videoId)
		_, err = db.ZDelete(key, commentId)
		if err != nil {
			glog.Error("zdelete failed. err:%v", err)
		}
	}
}

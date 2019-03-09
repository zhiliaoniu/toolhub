package util

import (
	"common"
	"db"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	pb "twirprpc"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"
)

type VideoDynamicInfo struct {
	LikeCount       uint32
	CommentCount    uint32
	ViewCount       uint32
	ShareCount      uint32
	CommentDisabled bool
	LikeDisabled    bool
	ShareDisabled   bool
}

func GetTopicVideoOffsetByDirction(userId uint64, topicIds []string, direction pb.Direction) (topicIdMap map[string]string, err error) {
	topicIdMap = make(map[string]string, 0)
	key := fmt.Sprintf("%s%d", common.USER_VIEW_TOPIC, (userId % 1000))
	if direction == pb.Direction_ALL {
		for _, topicId := range topicIds {
			topicIdMap[topicId] = ""
		}
		return
	}

	fields := make([]interface{}, 0)
	var offsets []string
	var suffix string
	if direction == pb.Direction_TOP {
		suffix = "up"
	} else if direction == pb.Direction_BOTTOM {
		suffix = "down"
	}
	for _, topicId := range topicIds {
		field := fmt.Sprintf("%d_%s_%s", userId, topicId, suffix)
		fields = append(fields, field)
	}
	if offsets, err = db.HMGet(key, fields); err != nil {
		return
	}
	for i := 0; i < len(topicIds); i++ {
		//浏览过程中，关注新话题，偏移量中没有，从零开始取
		//if offsets[i] == "" {
		//	continue
		//}
		topicIdMap[topicIds[i]] = offsets[i]
	}

	return
}

func GetTopicVidByOffset(userId uint64, topicIdMap map[string]string, topicIds []string, direction pb.Direction, retPerTopicViewNum int) (vids []string, updateMap map[string]string, hasMore bool, err error) {
	vids = make([]string, 0)
	updateMap = make(map[string]string, 0)
	readyMap := make(map[string][]string, 0)
	var totalNum int = 0
	hasMore = false
	if direction == pb.Direction_TOP {
		for _, topicId := range topicIds {
			vid, ok := topicIdMap[topicId]
			if !ok {
				continue
			}
			if vid == "" {
				continue
			}
			key := fmt.Sprintf("%s%s", common.TOPICWITHVIDEOID, topicId)
			var offset int
			if offset, err = db.ZRevRank(key, vid); err != nil {
				glog.Error("get topic vid offset failed. err:%v", err)
				continue
			}
			if offset == 0 {
				glog.Debug("topic:%s has no newer video. first is vid:%s", topicId, vid)
				continue
			}
			offset -= 1
			var topicVids []string
			if offset > 10 {
				topicVids, err = db.ZRevRange(key, int64(offset-retPerTopicViewNum+1), int64(offset))
			} else {
				topicVids, err = db.ZRevRange(key, 0, int64(offset))
			}
			if err != nil {
				glog.Error("get topic vid failed. err:%v", err)
				continue
			}
			l := len(topicVids)
			totalNum += l
			if l == 0 {
				glog.Debug("topicid:%s has no video", topicId)
				continue
			}
			glog.Debug("topicid:%s, len:%d, topicVids:%v", topicId, l, topicVids)
			readyMap[topicId] = topicVids
		}
	} else if direction == pb.Direction_BOTTOM {
		for _, topicId := range topicIds {
			vid, ok := topicIdMap[topicId]
			if !ok {
				continue
			}
			key := fmt.Sprintf("%s%s", common.TOPICWITHVIDEOID, topicId)
			offset := 0
			if vid != "" {
				if offset, err = db.ZRevRank(key, vid); err != nil {
					glog.Error("get topic vid offset failed. key:%s, vid:%s, err:%v", key, vid, err)
					continue
				}
				offset += 1
			}
			topicVids, err := db.ZRevRange(key, int64(offset), int64(offset+retPerTopicViewNum-1))
			if err != nil {
				glog.Error("get topic vid failed. err:%v", err)
				continue
			}
			l := len(topicVids)
			totalNum += l
			if l == 0 {
				glog.Debug("topicid:%s has no video", topicId)
				continue
			}
			glog.Debug("topicid:%s, len:%d, topicVids:%v", topicId, l, topicVids)
			readyMap[topicId] = topicVids
		}
	} else if direction == pb.Direction_ALL {
		//topMap := make(map[string]string, 0)
		for _, topicId := range topicIds {
			key := fmt.Sprintf("%s%s", common.TOPICWITHVIDEOID, topicId)
			topicVids, err := db.ZRevRange(key, 0, int64(retPerTopicViewNum-1))
			if err != nil {
				glog.Error("get topic vid failed. err:%v", err)
				continue
			}
			l := len(topicVids)
			if l == 0 {
				glog.Debug("topicid:%s has no video", topicId)
				continue
			}
			totalNum += l
			glog.Debug("topicid:%s, len:%d, topicVids:%v", topicId, l, topicVids)
			readyMap[topicId] = topicVids
		}
	}

	//only ret retPerTopicViewNum
	curNum := 0
	topMap := make(map[string]string, 0)
LOOP:
	for index := 0; index < retPerTopicViewNum; index++ {
		for _, topicId := range topicIds {
			topicVids, ok := readyMap[topicId]
			if !ok {
				continue
			}
			l := len(topicVids)
			if l <= index {
				continue
			}
			if direction == pb.Direction_TOP {
				vids = append(vids, topicVids[l-1-index])
				updateMap[topicId] = topicVids[l-1-index]
				curNum++
			} else if direction == pb.Direction_BOTTOM {
				vids = append(vids, topicVids[index])
				updateMap[topicId] = topicVids[index]
				curNum++
			} else if direction == pb.Direction_ALL {
				vids = append(vids, topicVids[index])
				updateMap[topicId] = topicVids[index]
				topMap[topicId] = topicVids[0]
				curNum++
			}
			//如果按关注话题的顺序返回视频，第一轮，尽可能所有关注话题都选择一个，
			//这样返回的视频可能多余十个，
			//但是避免只返回前十个话题的视频
			if index >= 1 && retPerTopicViewNum <= curNum {
				break LOOP
			}
		}
	}
	if curNum <= totalNum {
		hasMore = true
	}

	//update 用户关注上部偏移量
	if direction == pb.Direction_ALL && len(topMap) != 0 {
		UpdateTopicVideoOffsetByDirction(userId, topMap, pb.Direction_TOP)
	}

	return
}

func UpdateTopicVideoOffsetByDirction(userId uint64, updateMap map[string]string, direction pb.Direction) (err error) {
	glog.Debug("userId:%d, direction:%d, updateMap:%v", userId, direction, updateMap)
	key := fmt.Sprintf("%s%d", common.USER_VIEW_TOPIC, (userId % 1000))
	fields := make([]interface{}, 0)
	suffix := "up"
	if direction != pb.Direction_TOP {
		suffix = "down"
	}
	for topicId, vid := range updateMap {
		field := fmt.Sprintf("%d_%s_%s", userId, topicId, suffix)
		fields = append(fields, field, vid)
	}
	_, err = db.HMSet(key, fields)

	return
}

func GetVideoDynamicInfo(vids []string) (dynamicInfoMap map[string]*VideoDynamicInfo, err error) {
	dynamicInfoMap = make(map[string]*VideoDynamicInfo, 0)
	fields := make([]interface{}, 0)
	for _, vid := range vids {
		fields = append(fields, common.LIKE_COUNT+"_"+vid, common.COMMENT_COUNT+"_"+vid, common.VIEW_COUNT+"_"+vid, common.SHARE_COUNT+"_"+vid, common.COMMENT_DISABLED+"_"+vid, common.LIKE_DISABLED+"_"+vid, common.SHARE_DISABLED+"_"+vid)
	}
	intArr, err := db.HMGetInt(common.VIDEO_DYNAMIC, fields)
	if err != nil {
		glog.Error("hmget video dynamic info failed. key:%s, fields:%v, err:%v", common.VIDEO_DYNAMIC, fields, err)
		return
	}
	for index, vid := range vids {
		dynamicInfo := &VideoDynamicInfo{
			LikeCount:       uint32(intArr[index*7+0]),
			CommentCount:    uint32(intArr[index*7+1]),
			ViewCount:       uint32(intArr[index*7+2]),
			ShareCount:      uint32(intArr[index*7+3]),
			CommentDisabled: intArr[index*7+4] == 1,
			LikeDisabled:    intArr[index*7+5] == 1,
			ShareDisabled:   intArr[index*7+6] == 1,
		}
		if dynamicInfo.LikeCount != 0 ||
			dynamicInfo.CommentCount != 0 ||
			dynamicInfo.ViewCount != 0 ||
			dynamicInfo.ShareCount != 0 ||
			dynamicInfo.CommentDisabled != false ||
			dynamicInfo.LikeDisabled != false ||
			dynamicInfo.ShareDisabled != false {
			dynamicInfoMap[vid] = dynamicInfo
		}
	}
	return
}

func GetUserLikeVideosState(uid uint64, vids []string) (vidMap map[string]bool, err error) {
	vidsStr := strings.Join(vids, ",")
	tableName, err := db.GetTableName("user_favor_video_", uid)
	querySql := fmt.Sprintf("select vid from %s where uid=%d and vid in(%s) and result=1", tableName, uid, vidsStr)
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. err:", err)
		return
	}
	defer rows.Close()
	var vid string
	vidMap = make(map[string]bool, 0)
	for rows.Next() {
		err = rows.Scan(&vid)
		if err != nil {
			glog.Error("scan failed. err:", err)
			return
		}
		vidMap[vid] = true
	}

	return
}

//TODO put into questionserver
type UserQuestion struct {
	QuestionId string
	OptionId   string
	Result     bool
}

func GetVideoQuestionDynamicInfo(questionIds []string) (questionMap map[string]map[string]int, err error) {
	questionMap = make(map[string]map[string]int, 0)
	if len(questionIds) == 0 {
		glog.Debug("questionIds is empty")
		return
	}
	//1.read dynamic info
	questionIdInter := common.TransStrArrToInterface(questionIds)
	chooseNums, err := db.HMGetInt(common.QUESTION_DYNAMIC, questionIdInter)
	if err != nil {
		glog.Error("get question dynamic info failed. questionIds:%v, err:%v", questionIds, err)
		return
	}

	//2.parse question option choosenum
	for index, num := range chooseNums {
		composeId := questionIds[index]
		ids := strings.Split(composeId, "_")
		questionId, optionId := ids[0], ids[1]

		question, ok := questionMap[questionId]
		if !ok {
			question = make(map[string]int, 0)
			questionMap[questionId] = question
		}
		question[optionId] = num
	}
	return
}

func GetUserAnswerVideosQuestionState(uid uint64, questionIds []string) (questionIdMap map[string]*UserQuestion, err error) {
	questionIdMap = make(map[string]*UserQuestion, 0)
	if len(questionIds) == 0 {
		glog.Debug("uid %d choosed video have no question", uid)
		return
	}
	questionIdsStr := strings.Join(questionIds, ",")
	tableName, err := db.GetTableName("user_question_", uid)
	querySql := fmt.Sprintf("select question_id, option_id, result from %s where uid=%d and question_id in(%s)", tableName, uid, questionIdsStr)
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. querySql:%s err:%v", querySql, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		userQuestion := &UserQuestion{}
		var result int
		if err = rows.Scan(&userQuestion.QuestionId, &userQuestion.OptionId, &result); err != nil {
			glog.Error("scan failed. err:", err)
			break
		}
		if result == 1 {
			userQuestion.Result = true
		} else {
			userQuestion.Result = false
		}
		questionIdMap[userQuestion.QuestionId] = userQuestion
	}

	return
}

func GetQuestionOptionState(uid uint64, questionIds []string) (optionIdMap map[string]uint32, err error) {
	optionIdMap = make(map[string]uint32, 0)
	if len(questionIds) == 0 {
		glog.Debug("uid %d choosed video have no question", uid)
		return
	}
	questionIdsStr := strings.Join(questionIds, ",")
	querySql := fmt.Sprintf("select option_id, answer_num from question_option where question_id in(%s)", questionIdsStr)
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. err:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var optionId string
		var answerNum uint32
		if err = rows.Scan(&optionId, &answerNum); err != nil {
			glog.Error("scan failed. err:", err)
			break
		}
		optionIdMap[optionId] = answerNum
	}

	return
}

//快速取视频的热评，一次读取redis
func GetVidsHotCommentQuikly(vids []string, hotCommentNum int) (vidCommentsMap map[string][]*pb.CommentItem, err error) {
	vidCommentsMap = make(map[string][]*pb.CommentItem, 0)
	//1.get hotcomment from article by vid
	vidsInter := common.TransStrArrToInterface(vids)
	articleItems, err := db.HMGet(common.VIDEO_HOT_COMMENT, vidsInter)
	if err != nil {
		glog.Error("get vid hot comment failed. err:%v", err)
		return
	}
	for index, articleItemStr := range articleItems {
		vid := vids[index]
		if articleItemStr == "" {
			glog.Debug("vid:%s has no hot comment", vids[index])
			continue
		}
		articleItem := &pb.ArticleItem{}
		err = proto.Unmarshal([]byte(articleItemStr), articleItem)
		if err != nil {
			glog.Error("unmarshal vid:%s hot comment failed. err:%v", vids[index], err)
			continue
		}
		//2.parse articleItem
		if articleItem.CommentItems == nil || len(articleItem.CommentItems) == 0 {
			glog.Debug("vid:%s article has no hot comment. err:%v", vids[index], err)
			continue
		}
		maxLen := hotCommentNum
		if maxLen > len(articleItem.CommentItems) {
			maxLen = len(articleItem.CommentItems)
		}
		glog.Debug("vid:%s has hot comment. first cid:%s", vids[index], articleItem.CommentItems[0].CommentId)
		vidCommentsMap[vid] = articleItem.CommentItems[:maxLen]
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

func GetVidsHotCommentWithUid(uid uint64, vids []string) (vidCommentMap map[string]*pb.CommentItem, err error) {
	vidCommentMap = make(map[string]*pb.CommentItem, 0)
	cidArr := make([]string, 0)
	for _, vid := range vids {
		//1.get hot comment id
		vchotKey := fmt.Sprintf("%s%s", VCHOT_PREFIX, vid)
		cids, err := db.ZRevRange(vchotKey, 0, 0)
		if err != nil {
			glog.Debug("get hot comment id from redis failed. err:%v", err)
			continue
		}
		if len(cids) == 0 {
			glog.Debug("video %s has no hot comment.", vid)
			continue
		}

		//2.get hot comment
		commentItems, err := GetMultiCommentFromRedis(cids)
		if err != nil {
			glog.Error("get comment item from redis failed. err:%v", err)
			continue
		}
		if len(commentItems) == 0 {
			glog.Error("video %s get hot comment failed.", vid)
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
		cidArr = append(cidArr, cids[0])
	}

	//4.从mysql中读取当前用户是否对这些评论点过赞
	if len(cidArr) == 0 {
		glog.Debug("vids:%v have no hot comment", vids)
		return
	}
	cidMap, err := GetCommentFavorCids(uid, cidArr)
	if err != nil {
		return
	}
	for _, commentItem := range vidCommentMap {
		if _, ok := cidMap[commentItem.CommentId]; ok {
			commentItem.Liked = true
		}
	}

	return
}

func GetTopicDynamicInfo(topicIds []string) (topicMap map[string]int, err error) {
	topicMap = make(map[string]int, 0)
	if len(topicIds) == 0 {
		glog.Debug("topicIds is empty")
		return
	}
	topicIdInter := common.TransStrArrToInterface(topicIds)
	subscribeNums, err := db.HMGetInt(common.TOPIC_DYNAMIC, topicIdInter)
	if err != nil {
		glog.Error("get topic dynamic info failed. topicIds:%v, err:%v", topicIds, err)
		return
	}
	for index, num := range subscribeNums {
		topicMap[topicIds[index]] = num
	}

	return
}

func GetUserSubscribedTopicInfo(uid uint64, topicIds []string) (topicMap map[string]*pb.TopicItem, err error) {
	if len(topicIds) == 0 {
		glog.Debug("topicIds is empty")
		return
	}
	topicMap = make(map[string]*pb.TopicItem, 0)
	//1. get static topicItem
	topicIdInter := common.TransStrArrToInterface(topicIds)
	topicHashItems, err := db.HMGet(common.TOPICHASH, topicIdInter)
	if err != nil {
		glog.Error("get topic info failed. err:%v", err)
		return
	}
	for index, topicHashItem := range topicHashItems {
		if topicHashItem == "" {
			glog.Error("topicId:%d have no topicItem in redis", topicIds[index])
			continue
		}
		topicItem := &pb.TopicItem{}
		err = proto.Unmarshal([]byte(topicHashItem), topicItem)
		if err != nil {
			glog.Error("unmarshal failed. err:%v", err)
			continue
		}
		topicMap[topicItem.TopicId] = topicItem
	}

	//2. get dynamic topicInfo. such as, subscribed_user_num
	if len(topicMap) == 0 {
		glog.Error("topicid:[%v] not in redis")
		return
	}
	topicIdsStr := strings.Join(topicIds, ",")
	querySql := fmt.Sprintf("select topic_id, user_num from topic where topic_id in(%s)", topicIdsStr)
	rows, err := db.Query(common.BUDAODB, querySql)
	if err != nil {
		glog.Error("read mysql failed. err:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var topicId string
		var userNum uint32
		if err = rows.Scan(&topicId, &userNum); err != nil {
			glog.Error("scan failed. err:", err)
			break
		}
		if topicItem, ok := topicMap[topicId]; ok {
			topicItem.SubscribeCount = userNum
		}
	}

	//3. get user topic subscribed state
	if uid == 0 {
		return
	}

	key := fmt.Sprintf("%s%d", common.USER_SUBSCRIBE_TOPIC_WITH_TIME, uid)
	topicIdArr, err := db.Zrange(key, 0, -1)
	if err != nil {
		glog.Error("get user subscribed topicid failed. err:%v", err)
		return
	}
	if len(topicIdArr) == 0 {
		glog.Debug("user %d not subscribed topic", uid)
		return
	}
	for _, topicId := range topicIdArr {
		if topic, ok := topicMap[topicId]; ok {
			topic.Subscribed = true
		}
	}

	return
}

// SliceDiff returns diff slice of slice1 - slice2.
func SliceDiff(slice1, slice2 []string) (diffslice []string) {
	exposeMap := make(map[string]string, 0)
	for _, v := range slice2 {
		exposeMap[v] = v
	}

	i := 0
	for _, vid := range slice1 {
		if _, ok := exposeMap[vid]; !ok {
			diffslice = append(diffslice, vid)
			i++
		}
		if i == 10 {
			break
		}
	}

	return
}

func SliceDiffWithNum(slice1, slice2 []string, num int) (diffslice []string) {
	exposeMap := make(map[string]string, 0)
	for _, v := range slice2 {
		exposeMap[v] = v
	}

	i := 0
	for _, vid := range slice1 {
		if _, ok := exposeMap[vid]; !ok {
			diffslice = append(diffslice, vid)
			i++
		}
		if i == num {
			break
		}
	}

	return
}

// GetListItems get listitems by vid
func GetListItemsReal(vidArr []string, userID uint64) (listItems []*pb.ListItem, err error) {
	//1.get video info
	vidArrInter := common.TransStrArrToInterface(vidArr)
	hashItems, _ := db.HMGet(common.FULLVIDEOHASHKEYREDIS, vidArrInter)
	listItemMap := make(map[string]*pb.ListItem, 0) //<vid, *listItem>
	var topicIds []string
	var vids []string
	var questionIds []string
	vidTopicIDMap := make(map[string]string, 0)
	for index, hashItem := range hashItems {
		if hashItem == "" {
			glog.Error("vid:%s videoItem is empty in redis", vidArr[index])
			continue
		}
		listItem := &pb.ListItem{}
		if err = proto.Unmarshal([]byte(hashItem), listItem); err != nil {
			glog.Error("Unmarshal pb message failed. err:%v", err)
			continue
		}
		vid := listItem.VideoItem.VideoId
		if listItem.TopicItem != nil {
			topicId := listItem.TopicItem.TopicId
			topicIds = append(topicIds, topicId)
			vidTopicIDMap[vid] = topicId
		}

		vids = append(vids, vid)
		if listItem.QuestionItem != nil {
			questionIds = append(questionIds, listItem.QuestionItem.QuestionId)
		}
		listItemMap[vid] = listItem
	}
	glog.Debug("###################, len(listItemMap):%d", len(listItemMap))
	if len(listItemMap) == 0 {
		glog.Error("no adapt video")
		err = errors.New("no adapt video")
		return
	}
	glog.Debug("listItemMap:%v", listItemMap)

	//3.1 get video dynamic info
	vidDynamicInfoMap, err := GetVideoDynamicInfo(vids)
	if err != nil {
		glog.Error("get dynamic info failed. vids:%v, err:%v", vids, err)
	}

	//4.get topic info
	topicMap, err := GetUserSubscribedTopicInfo(userID, topicIds)
	if err != nil {
		glog.Error("get topic info failed. err:%v", err)
		return
	}

	//5.get comment info
	vidCommentMap, err := GetVidsHotCommentWithUid(userID, vids)
	if err != nil {
		glog.Error("get video hot comment failed. err:%v", err)
		return
	}

	//6.get user answer question info; get option state
	questionIDMap, err := GetUserAnswerVideosQuestionState(userID, questionIds)
	if err != nil {
		glog.Error("get user video answer question state failed. err:%v", err)
		return
	}
	optionIDMap, err := GetQuestionOptionState(userID, questionIds)
	if err != nil {
		glog.Error("get user video answer question state failed. err:%v", err)
		return
	}

	//7.get user like video info by mysql
	vidMap, err := GetUserLikeVideosState(userID, vids)
	if err != nil {
		glog.Error("get user video like state failed. err:%v", err)
		return
	}

	//8.fill listItem
	recommendItem := &pb.RecommendItem{}
	recommendItem.Reason = "小编精选"
	for vid, listItem := range listItemMap {
		//8.1 fill question and option
		questionItem := listItem.QuestionItem
		if len(questionIDMap) != 0 && questionItem != nil {
			if userQuestion, ok := questionIDMap[questionItem.QuestionId]; ok {
				questionItem.ChooseOptionId = userQuestion.OptionId
			}
			var answerCount uint32 = 0
			for _, optionItem := range questionItem.Options {
				if chooseCount, ok := optionIDMap[optionItem.OptionId]; ok {
					optionItem.ChooseCount = chooseCount
					answerCount += chooseCount
				}
			}
			questionItem.AnsweredCount = answerCount
		}
		//8.2 fill user like video
		if _, ok := vidMap[vid]; ok {
			listItem.ArticleItem.Liked = true
		}
		//8.3 fill hot comment
		if commentItem, ok := vidCommentMap[vid]; ok {
			listItem.ArticleItem.CommentItems = []*pb.CommentItem{commentItem}
		}
		//8.4 fill recommend
		listItem.RecommendItem = recommendItem
		//8.5 fill topic info
		if len(topicMap) != 0 {
			if topic, ok := topicMap[vidTopicIDMap[vid]]; ok {
				listItem.TopicItem = topic
			}
		}
		//8.6 fill video dynamic info
		if dynamicInfo, ok := vidDynamicInfoMap[vid]; ok {
			if dynamicInfo.LikeCount != 0 ||
				dynamicInfo.CommentCount != 0 ||
				dynamicInfo.ViewCount != 0 ||
				dynamicInfo.ShareCount != 0 {
				listItem.ArticleItem = &pb.ArticleItem{
					LikeCount:    dynamicInfo.LikeCount,
					CommentCount: dynamicInfo.CommentCount,
					ViewCount:    dynamicInfo.ViewCount,
					ShareCount:   dynamicInfo.ShareCount,
				}
			}

			if dynamicInfo.CommentDisabled == true || dynamicInfo.LikeDisabled == true || dynamicInfo.ShareDisabled == true {
				listItem.SwitchItem = &pb.SwitchItem{
					CommentDisabled: dynamicInfo.CommentDisabled,
					LikeDisabled:    dynamicInfo.LikeDisabled,
					ShareDisabled:   dynamicInfo.ShareDisabled,
				}
			}
		}
	}
	listItems = make([]*pb.ListItem, 0)
	for _, listItem := range listItemMap {
		listItems = append(listItems, listItem)
	}
	glog.Debug("listItems numbers:%d", len(listItems))

	return
}

// GetListItems get listitems by vid
func GetListItemsNotReal(vids []string, userId uint64) (listItems []*pb.ListItem, err error) {
	glog.Debug("vids:%v", vids)
	userActMap := make(map[string]int64, 0)
	var wg sync.WaitGroup
	//judge user login
	if userId != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//1.get user act
			key := fmt.Sprintf("%s%d", common.USER_ACT_PREFIX, userId)
			userActMap, err = db.HGetAll(key)
			if err != nil {
				glog.Error("hgetall user:%d act failed. err:%v", userId, err)
			}
		}()
	}

	//2.get listitem
	wg.Add(1)
	go func() {
		defer wg.Done()
		vidsInter := common.TransStrArrToInterface(vids)
		hashItems, err := db.HMGet(common.FULLVIDEOHASHKEYREDIS, vidsInter)
		if err != nil {
			glog.Error("get listitem failed. err:%v", err)
			return
		}
		for index, hashItem := range hashItems {
			if hashItem == "" {
				glog.Error("vid:%s videoItem is empty in redis", vids[index])
				continue
			}
			listItem := &pb.ListItem{}
			if err := proto.Unmarshal([]byte(hashItem), listItem); err != nil {
				glog.Error("Unmarshal pb message failed. vid:%s, err:%v", vids[index], err)
				continue
			}

			listItems = append(listItems, listItem)
		}
	}()

	wg.Wait()

	if userId == 0 {
		return
	}

	//3.fill user act to listitem
	for _, listItem := range listItems {
		vid := listItem.VideoItem.VideoId

		if listItem.ArticleItem != nil {
			//3.1 favor video
			vidFavor := vid + common.VID_FAVOR_SUFFIX
			if _, ok := userActMap[vidFavor]; ok {
				listItem.ArticleItem.Liked = true
			}
			if listItem.ArticleItem.CommentItems != nil {
				//3.2 favor comment
				cid := listItem.ArticleItem.CommentItems[0].CommentId
				cidFavor := cid + common.CID_FAVOR_SUFFIX
				if _, ok := userActMap[cidFavor]; ok {
					listItem.ArticleItem.CommentItems[0].Liked = true
				}
			}
		}

		//3.3 subscribe topic
		if listItem.TopicItem != nil {
			topicKey := listItem.TopicItem.TopicId + common.TOPIC_SUBSCRIBE_SUFFIX
			if num, ok := userActMap[topicKey]; ok && num != 0 {
				listItem.TopicItem.Subscribed = true
			}
		} else {
			glog.Debug("vid:%s has no topic", listItem.VideoItem.VideoId)
		}

		//3.4 answer question
		if listItem.QuestionItem != nil {
			questionId := listItem.QuestionItem.QuestionId
			if optionId, ok := userActMap[questionId]; ok {
				listItem.QuestionItem.ChooseOptionId = strconv.FormatInt(optionId, 10)
			}
		}
	}
	return
}

// VideoStaticInfo video static infomation
type VideoStaticInfo struct {
	VID        uint64
	State      uint32
	CoverURL   string
	Title      string
	Width      uint32
	Height     uint32
	Duration   uint32
	CreateTime time.Time
}

// QuestionInfo define the question information from mysql
type QuestionInfo struct {
	ID         uint64
	VideoID    uint64
	State      uint32
	HardLever  uint32
	Score      uint32
	Content    string
	OptionType uint32
}

// OptionInfo question option
type OptionInfo struct {
	QuestionID uint64
	OptionID   uint64
	Content    string
	IsAnswer   uint32
	AnswerNum  uint32
}

// SliceRemoveDuplicates remove the duplicate element of slice
func SliceRemoveDuplicates(src []string) (dst []string) {
	srcLen := len(src)
	for i := 0; i < srcLen; i++ {
		if (i > 0 && src[i-1] == src[i]) || len(src[i]) == 0 {
			continue
		}
		dst = append(dst, src[i])
	}
	return
}

// GetQuestionItemByVideoID get question information by video id
func GetQuestionItemByVideoID(videoID uint64) (questionItem *pb.QuestionItem, err error) {
	sqlString := fmt.Sprintf("select id, vid, state, hard_level, score, content, option_style from question where vid = %d", videoID)
	rows, err := db.Query(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("query question failed. err:%v", err)
		return nil, err
	}
	defer rows.Close()

	question := QuestionInfo{}
	for rows.Next() {
		rows.Scan(&question.ID, &question.VideoID, &question.State, &question.HardLever, &question.Score, &question.Content, &question.OptionType)
	}

	if question.ID == 0 || question.VideoID == 0 {
		return nil, err
	}
	var style pb.QuestionItem_Style
	if question.OptionType == 0 {
		style = pb.QuestionItem_SINGLE
	} else {
		style = pb.QuestionItem_DOUBLE
	}
	optionItems, err := GetOptionItemsByQuestionID(question.ID)
	if len(optionItems) == 0 || err != nil {
		return nil, err
	}

	return &pb.QuestionItem{
		QuestionId: strconv.FormatUint(question.ID, 10),
		VideoId:    strconv.FormatUint(question.VideoID, 10),
		HardLevel:  question.HardLever,
		Score:      question.Score,
		Content:    question.Content,
		Options:    optionItems,
		Style:      style,
	}, nil
}

// GetOptionItemsByQuestionID get options of question
func GetOptionItemsByQuestionID(questionID uint64) (optionItems []*pb.OptionItem, err error) {
	sqlString := fmt.Sprintf("select question_id, option_id, option_content, is_answer, answer_num, fake_answer_num from question_option where question_id=%d", questionID)
	rows, err := db.Query(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("query question_option failed. err:%v", err)
		return nil, err
	}
	defer rows.Close()

	optionItems = make([]*pb.OptionItem, 0)
	for rows.Next() {
		option := OptionInfo{}
		right := false
		var FakeAnswerNum uint32
		rows.Scan(&option.QuestionID, &option.OptionID, &option.Content, &option.IsAnswer, &option.AnswerNum, &FakeAnswerNum)
		if option.IsAnswer == 1 {
			right = true
		}

		opionItem := &pb.OptionItem{
			OptionId:    strconv.FormatUint(option.OptionID, 10),
			QuestionId:  strconv.FormatUint(option.QuestionID, 10),
			Content:     option.Content,
			Right:       right,
			ChooseCount: option.AnswerNum + FakeAnswerNum,
		}
		optionItems = append(optionItems, opionItem)
	}

	return optionItems, err
}

func RecomposeListItemsByAbtest(listItems []*pb.ListItem, abtestItems []*pb.AbtestItem) {
	if abtestItems == nil {
		return
	}
	for _, listItem := range listItems {
		for _, abtestItem := range abtestItems {
			if abtestItem.Name != "questionOpenTest" {
				listItem.QuestionItem = nil
			}
		}
	}
}

func RecomposeListItems(listItems []*pb.ListItem, abtestItems []*pb.AbtestItem, recommendItems map[string]*pb.RecommendVideoItem, channelId string) {
	//1.change by abtest
	RecomposeListItemsByAbtest(listItems, abtestItems)

	//2.fill channelid
	for _, listItem := range listItems {
		if listItem.ArticleItem != nil {
			listItem.ArticleItem.ChannelId = channelId
		}
	}

	//3.fill recommendItem
	tarceID := strconv.FormatUint(uint64(time.Now().UnixNano()), 10)
	if recommendItems == nil {
		for _, listItem := range listItems {
			if listItem.RecommendItem == nil {
				continue
			}
			listItem.RecommendItem.TraceId = tarceID
		}
	} else {
		for _, listItem := range listItems {
			if listItem.RecommendItem == nil {
				continue
			}
			if recommendItem, ok := recommendItems[listItem.VideoItem.VideoId]; ok {
				listItem.RecommendItem.TraceId = recommendItem.TraceId
				listItem.RecommendItem.Payload = recommendItem.Payload
			} else {
				listItem.RecommendItem.TraceId = tarceID
			}
		}
	}
}

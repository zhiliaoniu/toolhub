package util

import (
	"common"
	"db"
	"fmt"
	"strconv"
	"strings"
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

//快速取视频的热评，一次读取redis
func GetVidsHotCommentQuickly(vids []string, hotCommentNum int) (vidCommentsMap map[string][]*pb.CommentItem, err error) {
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
		//glog.Debug("vid:%s has hot comment. first cid:%s", vids[index], articleItem.CommentItems[0].CommentId)
		vidCommentsMap[vid] = articleItem.CommentItems[:maxLen]
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

// OptionInfo question option
type OptionInfo struct {
	QuestionID uint64
	OptionID   uint64
	Content    string
	IsAnswer   uint32
	AnswerNum  uint32
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

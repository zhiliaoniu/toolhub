package timelineserver

import (
	"common"
	"db"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"

	"service/util"
	pb "twirprpc"
)

// VideoFullInfo video full information
type VideoFullInfo struct {
	VID            uint64
	Title          string
	CoverURL       string
	VideoURL       string
	Duration       uint32
	Width          uint32
	Height         uint32
	State          uint32
	OpState        uint32
	ViewNum        uint32
	FakeViewNum    uint32
	FavorNum       uint32
	FakeFavorNum   uint32
	CommentNum     uint32
	FakeCommentNum uint32
	ShareNum       uint32
	FakeShareNum   uint32
	CreateTime     time.Time
	CouldFavor     uint32
	CouldComment   uint32
	CouldShare     uint32
	TopicID        string
}

// QuestionInfo define the question information from mysql
type QuestionInfo struct {
	ID             uint64
	VideoID        uint64
	State          uint32
	HardLever      uint32
	Score          uint32
	RightAnswerNum uint32
	WrongAnswerNum uint32
	Content        string
	OptionType     uint32
}

// TopicInfo topic
type TopicInfo struct {
	TopicID     uint64
	Name        string
	Pic         string
	Disable     uint32
	Description string
	UserNum     uint32
	FakeUserNum uint32
	VideoNum    uint32
}

func (s *Server) CronTask() {
	ticker := time.NewTicker(time.Second * time.Duration(10))
	for {
		if s.cronTaskDisable {
			glog.Debug("exit config server crontab")
			break
		}
		select {
		case <-ticker.C:
			UpdateListItems()
		case <-time.After(time.Second):
		}
	}
}

// UpdateListItems update listItem
func UpdateListItems() {
	var err error
	var exist bool
	if exist, err = db.KExists(common.VID_DYNAMIC); err != nil {
		glog.Error("exists failed. err:%v", err)
		return
	}
	if exist == false {
		//glog.Debug("no vid need update in %s.", common.VID_DYNAMIC)
		return
	}

	glog.Debug("=============================begin crontab update list item ")
	//1.read ready_update_vid
	var r int
	if r, err = db.KRenameNX(common.VID_DYNAMIC, common.VID_DYNAMIC_TMP); err != nil {
		glog.Error("renamenx failed. err:%v", err)
		return
	}
	if r == 0 {
		glog.Debug("new key already exist.")
		return
	}
	defer func() {
		//delete tmp key
		if r, err := db.DelKey(common.VID_DYNAMIC_TMP); err != nil || r != 1 {
			glog.Error("del key:%s failed. err:%v, r:%d", common.VID_DYNAMIC_TMP, err, r)
		}
	}()
	readyUpdateVids, err := db.Smembers(common.VID_DYNAMIC_TMP)
	if err != nil {
		glog.Error("read redis set failed. err:%v", err)
		return
	}
	if len(readyUpdateVids) == 0 {
		glog.Debug("readyUpdateVids is empty.")
		return
	}
	glog.Debug("ready update readyUpdateVids:%v", readyUpdateVids)

	//2.read old listitem
	readyUpdateVidsInter := common.TransStrArrToInterface(readyUpdateVids)
	hashItems, err := db.HMGet(common.FULLVIDEOHASHKEYREDIS, readyUpdateVidsInter)
	if err != nil {
		glog.Error("get old listitem failed. err:%v", err)
		return
	}
	listItemMap := make(map[string]*pb.ListItem, 0) //<vid, *listItem>
	var vids, topicIds, questionIds []string
	for index, hashItem := range hashItems {
		if hashItem == "" {
			glog.Error("vid:%s videoItem is empty in redis", readyUpdateVids[index])
			continue
		}
		listItem := &pb.ListItem{}
		if err := proto.Unmarshal([]byte(hashItem), listItem); err != nil {
			glog.Error("Unmarshal pb message failed. vid:%s, err:%v", readyUpdateVids[index], err)
			continue
		}

		vid := listItem.VideoItem.VideoId
		vids = append(vids, vid)
		if listItem.QuestionItem != nil {
			questionId := listItem.QuestionItem.QuestionId
			for _, option := range listItem.QuestionItem.Options {
				composeId := fmt.Sprintf("%s_%s", questionId, option.OptionId)
				questionIds = append(questionIds, composeId)
			}
		}
		if listItem.TopicItem != nil {
			topicIds = append(topicIds, listItem.TopicItem.TopicId)
		}
		listItemMap[vid] = listItem
	}
	glog.Debug("len.listItemMap:%d", len(listItemMap))
	if len(listItemMap) == 0 {
		glog.Error("no adapt video")
		return
	}
	//glog.Debug("listItemMap:%v", listItemMap)

	//3.重组listitem
	if listItemMap, err = RecomposeListitems(vids, questionIds, topicIds, listItemMap); err != nil {
		glog.Error("RecomposeLisitems failed. err:%v", err)
		return
	}

	//4.write listitem to redis
	var fields []interface{}
	for vid, listItem := range listItemMap {
		listItemstr, err := proto.Marshal(listItem)
		if err != nil {
			glog.Error("marshal listitem:%v failed. err:%v", err)
			continue
		}
		fields = append(fields, vid, listItemstr)
	}

	if len(fields) == 0 {
		glog.Error("fields is empty, not update")
		return
	}
	if _, err := db.HMSet(common.FULLVIDEOHASHKEYREDIS, fields); err != nil {
		glog.Error("hmset failed. err:%v", err)
	}

	glog.Debug("=============================end crontab update list item ")
}

// RecomposeListitems 重组listitem
func RecomposeListitems(vids, questionIds []string, topicIds []string, listItemMap map[string]*pb.ListItem) (listItemMap2 map[string]*pb.ListItem, err error) {
	//1. 读取视频动态信息,赞、评、转、开关
	vidDynamicInfoMap, err := util.GetVideoDynamicInfo(vids)
	if err != nil {
		glog.Error("get dynamic info failed. vids:%v, err:%v", vids, err)
	}
	//2 读取最新热评信息
	vidHotCommentMap, err := util.GetVidsHotCommentQuickly(vids, 1)
	if err != nil {
		glog.Error("get video hot comment failed. vids:%v err:%v", vids, err)
	}

	//3 读取问题动态信息,选项选择人数
	questionMap, err := util.GetVideoQuestionDynamicInfo(questionIds)
	if err != nil {
		glog.Error("get video question dynamic info failed. vids:%v, questionIds:%v, err:%v", vids, questionIds, err)
	}

	//4 读取话题动态信息,订阅用户数
	topicMap, err := util.GetTopicDynamicInfo(topicIds)
	if err != nil {
		glog.Error("get topic info failed. err:%v", err)
		return
	}

	//5 fill dynamic
	for _, listItem := range listItemMap {
		vid := listItem.VideoItem.VideoId
		//5.1 fill video dynamic
		if dynamic, ok := vidDynamicInfoMap[vid]; ok {
			listItem.ArticleItem = &pb.ArticleItem{
				LikeCount:    dynamic.LikeCount,
				CommentCount: dynamic.CommentCount,
				ViewCount:    dynamic.ViewCount,
				ShareCount:   dynamic.ShareCount,
			}
			listItem.SwitchItem = &pb.SwitchItem{
				CommentDisabled: dynamic.CommentDisabled,
				LikeDisabled:    dynamic.LikeDisabled,
				ShareDisabled:   dynamic.ShareDisabled,
			}
		}

		//5.2 fill hot comment
		if commentItemArr, ok := vidHotCommentMap[vid]; ok {
			if listItem.ArticleItem == nil {
				listItem.ArticleItem = &pb.ArticleItem{}
			}
			listItem.ArticleItem.CommentItems = commentItemArr
		}

		//5.3 fill question dynamic
		questionItem := listItem.QuestionItem
		if questionItem != nil {
			if question, ok := questionMap[questionItem.QuestionId]; ok {
				total := 0
				for _, option := range questionItem.Options {
					if num, ok := question[option.OptionId]; ok {
						question[option.OptionId] = num
						total += num
					}
				}
				questionItem.AnsweredCount = uint32(total)
			}
		}

		//5.4 fill topic info
		if listItem.TopicItem != nil {
			if num, ok := topicMap[listItem.TopicItem.TopicId]; ok {
				listItem.TopicItem.SubscribeCount = uint32(num)
			}
		}
	}

	glog.Debug("listItemMap numbers:%d", len(listItemMap))

	return listItemMap, nil
}

// TimerLoadVideoFullInfo periodically loading video static info from mysql to redis
func TimerLoadVideoFullInfo() {
	videoTableNum := common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc["video_"]
	if exist, err := db.KExists(common.FULLVIDEOHASHKEYREDIS); err != nil || !exist {
		for i := 0; i < int(videoTableNum); i++ {
			videoTableName := fmt.Sprintf("video_%d", i)
			glog.Debug("start to load video static info.")
			err := LoadVideoFullInfotoRedis(videoTableName)
			if err != nil {
				glog.Error("LoadVideoStaticInfotoRedis failed. err:%v", err)
				continue
			}
			glog.Debug("end to load video static info.")
		}
	}

	ticker := time.NewTicker(time.Second * time.Duration(common.GetConfig().TimelineConf.CronTaskInternalSec))
	go func() {
		for {
			select {
			case <-ticker.C:
				for i := 0; i < int(videoTableNum); i++ {
					videoTableName := fmt.Sprintf("video_%d", i)
					err := LoadVideoFullInfotoRedis(videoTableName)
					if err != nil {
						glog.Error("LoadVideoStaticInfotoRedis failed. err:%v", err)
						continue
					}
				}
				glog.Debug("TimerLoadVideoFullInfo success")
			}
		}
	}()
}

// LoadVideoFullInfotoRedis load listItem full info into redis hash
//TODO 加载listitem全量信息，定时更新。同时更新video_dynamic_info、question、topic、comment_dynamic
func LoadVideoFullInfotoRedis(tableName string) (err error) {
	glog.Debug("LoadVideoFullInfotoRedis %s begin", tableName)
	execSQL := fmt.Sprintf("select vid, title, coverurl, videourl, duration, width, height, state, op_state, view_num, fake_view_num, favor_num, fake_favor_num, comment_num, fake_comment_num, share_num, fake_share_num, create_time, could_favor, could_comment, could_share, topic_id from %s", tableName)
	rows, err := db.Query(common.BUDAODB, execSQL)
	if err != nil {
		glog.Error("query %s failed. err:%v", tableName, err)
		return err
	}
	defer rows.Close()

	var vidArr []string
	listItemMap := make(map[string]*pb.ListItem, 0) // <vid, *listItem>
	vidTopicIDMap := make(map[string][]string, 0)   // <topicid, vid slice>

	for rows.Next() {
		videoInfo := VideoFullInfo{}
		rows.Scan(&videoInfo.VID, &videoInfo.Title, &videoInfo.CoverURL, &videoInfo.VideoURL, &videoInfo.Duration, &videoInfo.Width, &videoInfo.Height, &videoInfo.State, &videoInfo.OpState, &videoInfo.ViewNum, &videoInfo.FakeViewNum, &videoInfo.FavorNum, &videoInfo.FakeFavorNum, &videoInfo.CommentNum, &videoInfo.FakeCommentNum, &videoInfo.ShareNum, &videoInfo.FakeShareNum, &videoInfo.CreateTime, &videoInfo.CouldFavor, &videoInfo.CouldComment, &videoInfo.CouldShare, &videoInfo.TopicID)
		if videoInfo.State != 2 || videoInfo.OpState != 0 {
			exist, _ := db.HExists(common.FULLVIDEOHASHKEYREDIS, strconv.FormatUint(videoInfo.VID, 10))
			if exist == true {
				db.HDelete(common.FULLVIDEOHASHKEYREDIS, strconv.FormatUint(videoInfo.VID, 10))
			}
			continue
		}
		var topicID string
		exist := strings.Contains(videoInfo.TopicID, ",")
		if exist == true {
			temp := strings.Split(videoInfo.TopicID, ",")
			topicID = temp[0]
		} else {
			topicID = videoInfo.TopicID
		}

		videoShareURL := fmt.Sprintf("https://budao-web.yy.com/share-video-h5/index.html?topicid=%s&videoid=%d", topicID, videoInfo.VID)
		videoShareItem := &pb.ShareItem{
			Title:    videoInfo.Title,
			IconUrl:  videoInfo.CoverURL,
			ShareUrl: videoShareURL,
		}
		loc, _ := time.LoadLocation("Local")
		theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", videoInfo.CreateTime.Format("2006-01-02 15:04:05"), loc)
		videoItem := &pb.VideoItem{
			VideoId:     strconv.FormatUint(videoInfo.VID, 10),
			PictureUrl:  videoInfo.CoverURL,
			Title:       videoInfo.Title,
			Width:       videoInfo.Width,
			Height:      videoInfo.Height,
			Duration:    videoInfo.Duration,
			PublishTime: uint64(theTime.Unix()),
			ShareItem:   videoShareItem,
		}

		recommendItem := &pb.RecommendItem{
			Reason: "小编精选",
		}

		var maxRation float32
		var minRation float32
		var screen pb.SwitchItem_FullscreenType
		// temp
		if videoInfo.Width == 0 || videoInfo.Height == 0 {
			videoInfo.Height = 440
			videoInfo.Width = 680
		}
		if videoInfo.Height > videoInfo.Width {
			maxRation = float32(1)
			screen = 0
		} else {
			maxRation = float32(videoInfo.Height) / float32(videoInfo.Width)
			screen = 1
		}
		minRation = float32(9) / float32(16)
		switchItem := &pb.SwitchItem{
			ListVideoRatio:      float32(9) / float32(16),
			DetailVideoMaxRatio: maxRation,
			DetailVideoMinRatio: minRation,
			FullscreenType:      screen,
		}
		if videoInfo.CouldFavor == 0 {
			switchItem.LikeDisabled = true
		}
		if videoInfo.CouldComment == 0 {
			switchItem.CommentDisabled = true
		}
		if videoInfo.CouldShare == 0 {
			switchItem.ShareDisabled = true
		}

		articleItem := &pb.ArticleItem{
			LikeCount:    videoInfo.FavorNum + videoInfo.FakeFavorNum,
			CommentCount: videoInfo.CommentNum + videoInfo.FakeCommentNum,
			ViewCount:    videoInfo.ViewNum + videoInfo.FakeViewNum,
			ShareCount:   videoInfo.ShareNum + videoInfo.FakeShareNum,
		}

		listItem := &pb.ListItem{
			VideoItem:     videoItem,
			ArticleItem:   articleItem,
			RecommendItem: recommendItem,
			SwitchItem:    switchItem,
		}
		key := strconv.FormatUint(videoInfo.VID, 10)
		vidArr = append(vidArr, key)
		listItemMap[key] = listItem
		vidTopicIDMap[topicID] = append(vidTopicIDMap[topicID], key)

	}
	glog.Debug("read video table end")

	//------------- begin question ------------
	if common.GetConfig().TimelineConf.ShowQuestion == true {
		sqlString := fmt.Sprintf("select id, vid, state, hard_level, score, right_answer_num, wrong_answer_num, content, option_style from question")
		questionRows, err := db.Query(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("query question failed. err:%v", err)
			return err
		}
		defer questionRows.Close()

		for questionRows.Next() {
			question := QuestionInfo{}
			questionRows.Scan(&question.ID, &question.VideoID, &question.State, &question.HardLever, &question.Score, &question.RightAnswerNum, &question.WrongAnswerNum, &question.Content, &question.OptionType)

			if question.ID == 0 || question.VideoID == 0 {
				continue
			}

			var style pb.QuestionItem_Style
			if question.OptionType == 0 {
				style = pb.QuestionItem_SINGLE
			} else {
				style = pb.QuestionItem_DOUBLE
			}
			optionItems, err := util.GetOptionItemsByQuestionID(question.ID)
			if len(optionItems) == 0 || err != nil {
				continue
			}
			questionItem := &pb.QuestionItem{
				QuestionId:    strconv.FormatUint(question.ID, 10),
				VideoId:       strconv.FormatUint(question.VideoID, 10),
				HardLevel:     question.HardLever,
				Score:         question.Score,
				Content:       question.Content,
				Options:       optionItems,
				AnsweredCount: question.RightAnswerNum + question.WrongAnswerNum,
				Style:         style,
			}
			//glog.Debug("questionItem:%v", questionItem)
			videoID := strconv.FormatUint(question.VideoID, 10)
			if listItemMap[videoID] != nil {
				listItemMap[videoID].QuestionItem = questionItem
			} else {
				glog.Debug("videoID:%s is not valid. question:%v ", videoID, questionItem)
			}
		}
		glog.Debug("read question table end")
	}
	//------------- end question ------------

	sqlString := fmt.Sprintf("select topic_id, name, pic, disable, description, user_num, fake_user_num, video_num from topic")
	topicRows, err := db.Query(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("query topic failed. err:%v", err)
		return err
	}
	defer topicRows.Close()
	for topicRows.Next() {
		topicInfo := TopicInfo{}
		topicRows.Scan(&topicInfo.TopicID, &topicInfo.Name, &topicInfo.Pic, &topicInfo.Disable, &topicInfo.Description, &topicInfo.UserNum, &topicInfo.FakeUserNum, &topicInfo.VideoNum)
		if topicInfo.Disable == 1 {
			continue
		}

		topicShareURL := fmt.Sprintf("https://budao-web.yy.com/share-topic-h5/index.html?topicid=%d", topicInfo.TopicID)
		topicShareItem := &pb.ShareItem{
			Title:    topicInfo.Name,
			Desc:     topicInfo.Description,
			IconUrl:  topicInfo.Pic,
			ShareUrl: topicShareURL,
		}
		topicItem := &pb.TopicItem{
			TopicId:        strconv.FormatUint(topicInfo.TopicID, 10),
			Name:           topicInfo.Name,
			IconUrl:        topicInfo.Pic,
			Desc:           topicInfo.Description,
			SubscribeCount: topicInfo.UserNum + topicInfo.FakeUserNum,
			ShareItem:      topicShareItem,
			ShareDisabled:  common.GetConfig().TopicConf.ShareDisabled,
		}
		topicID := strconv.FormatUint(topicInfo.TopicID, 10)
		if vidArr, ok := vidTopicIDMap[topicID]; ok {
			for _, vid := range vidArr {
				if listItemMap[vid] != nil {
					listItemMap[vid].TopicItem = topicItem
					//glog.Debug("topicItem:%v", topicItem)
				} else {
					glog.Error("vid:%s, topicItem:%v", vid, topicItem)
				}
			}
		}
	}
	glog.Debug("read topic table end")

	//1: 外显个数
	glog.Debug("read hot comment table begin")
	vidHotCommentArr, err := util.GetVidsHotCommentQuickly(vidArr, 1)
	if err != nil {
		glog.Error("get vid hot comment failed. err:%v", err)
	} else {
		for vid, commentItems := range vidHotCommentArr {
			listItemMap[vid].ArticleItem.CommentItems = commentItems
		}
	}
	glog.Debug("read hot comment table end")

	loop := 0
	var listItemPackage []interface{}
	for _, vid := range vidArr {
		// Serialization videoItem
		data, err := proto.Marshal(listItemMap[vid])
		if err != nil {
			glog.Error("Serialization videoItem failed. err:%v", err)
			continue
		}
		listItemPackage = append(listItemPackage, vid, data)
		loop++
		if loop >= 500 {
			ret, err := db.HMSet(common.FULLVIDEOHASHKEYREDIS, listItemPackage)
			if ret != "ok" && err != nil {
				glog.Error("insert video static info into redis failed. err:%v", err)
				continue
			}
			loop = 0
			listItemPackage = make([]interface{}, 0)
		}
	}
	glog.Debug("marshal listItem end")
	if len(listItemPackage) != 0 {
		ret, err := db.HMSet(common.FULLVIDEOHASHKEYREDIS, listItemPackage)
		if ret != "ok" && err != nil {
			glog.Error("hmset failed. err:%v", err)
		}
	}
	glog.Debug("LoadVideoFullInfotoRedis %s end", tableName)

	return err
}

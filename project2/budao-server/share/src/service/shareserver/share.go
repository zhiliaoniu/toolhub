package shareserver

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sumaig/glog"

	"common"
	"db"
	pb "twirprpc"
)

// Server identify for ShareService RPC
type Server struct{}

const (
	COMMENTDYNAMICPREFIX = "comment_dynamic_"
	FAVORNUMPREFIX       = "favor_num_"
	REPLYNUMPREFIX       = "reply_num_"
	VIDEOLIKENUM         = "like_count_"
	VIDEOCOMMENTNUM      = "comment_count_"
	VIDDYNAMIC           = "video_dynamic"
)

// CommentDynamicInfo comment dynamic information
type CommentDynamicInfo struct {
	IsReply  bool
	FavorNum uint32
	ReplyNum uint32
	Weight   uint32
}

// GetServer return server of share service
func GetServer() *Server {
	server := &Server{}

	return server
}

// ShareVideoBottomPage share video bottom page
func (s *Server) ShareVideoBottomPage(ctx context.Context, req *pb.ShareVideoBottomPageRequest) (resp *pb.ShareVideoBottomPageResponse, err error) {
	resp = &pb.ShareVideoBottomPageResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	begin := time.Now()
	videoID := req.GetVideoId()
	topicID := req.GetTopicId()
	glog.Debug("videoID: %s", videoID)
	glog.Debug("topicID: %s", topicID)
	if videoID == "" {
		resp.Status.Code = pb.Status_OK
		resp.Status.Message = "request parameter videoID is empty"
		return
	}

	glog.Debug("###########################111")
	listHash, err := db.HGet(common.FULLVIDEOHASHKEYREDIS, videoID)
	if err != nil {
		glog.Error("get listItem from redis failed. err:%v", err)
		return
	}
	glog.Debug("###########################222")
	listItem := &pb.ListItem{}
	if err = proto.Unmarshal([]byte(listHash), listItem); err != nil {
		glog.Error("Unmarshal pb message failed. videoID:%s, err:%v", videoID, err)
		return
	}
	glog.Debug("###########################333")
	videoItem := &pb.VideoItem{
		VideoId:    listItem.VideoItem.VideoId,
		PictureUrl: listItem.VideoItem.PictureUrl,
		Title:      listItem.VideoItem.Title,
		Width:      listItem.VideoItem.Width,
		Height:     listItem.VideoItem.Height,
		Duration:   listItem.VideoItem.Duration,
	}

	if topicID == "" {
		if listItem.TopicItem != nil {
			topicID = listItem.TopicItem.TopicId
		}

	}

	glog.Debug("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	topicItem := &pb.TopicItem{}
	var toppicData string
	if topicID != "" {
		toppicData, err = db.HGet(common.TOPICHASH, topicID)
		if err != nil {
			glog.Error("get topicItem from redis failed. err:%v", err)
			return
		}
		if err := proto.Unmarshal([]byte(toppicData), topicItem); err != nil {
			glog.Error("Unmarshal pb message failed. topicID:%s, err:%v", topicID, err)
		}
	}

	vidLikeNumKey := fmt.Sprintf("%s%s", VIDEOLIKENUM, videoID)
	vidCommentNumKey := fmt.Sprintf("%s%s", VIDEOCOMMENTNUM, videoID)
	var vidDynaInter []interface{}
	vidDynaInter = append(vidDynaInter, vidLikeNumKey, vidCommentNumKey)
	vidDynaInfo, err := db.HMGetInt(VIDDYNAMIC, vidDynaInter)
	if err != nil {
		glog.Error("get video dynamic information from redis failed. err:%v", err)
		return
	}
	articleItem := &pb.ArticleItem{
		LikeCount:    uint32(vidDynaInfo[0]),
		CommentCount: uint32(vidDynaInfo[1]),
	}

	vchotString := fmt.Sprintf("%s%s", common.VCHOTPREFIX, videoID)
	cidsHash, err := db.ZRevRange(vchotString, 0, common.GetConfig().VideoBottomPage.CommentNUM-1)
	if err != nil {
		glog.Error("get cids from redis failed. err:%v", err)
		return
	}
	var hashComments []string
	var comments []*pb.CommentItem
	if len(cidsHash) != 0 {
		glog.Debug("commentnum:%d", len(cidsHash))
		cidsArr := common.TransStrArrToInterface(cidsHash)
		hashComments, err = db.HMGet(common.COMMENTITEMALLKEY, cidsArr)
		if err != nil {
			glog.Error("get commentItems from redis failed. err:%v", err)
			return
		}

		for _, hashComment := range hashComments {
			if hashComment == "" {
				glog.Error("commentItem is empty in redis")
				continue
			}
			commentItem := &pb.CommentItem{}
			if err := proto.Unmarshal([]byte(hashComment), commentItem); err != nil {
				glog.Error("Unmarshal pb message failed. err:%v", err)
				continue
			}
			targetItem := &pb.CommentItem{
				CommentId:   commentItem.CommentId,
				UserItem:    commentItem.UserItem,
				Content:     commentItem.Content,
				ReplyCount:  commentItem.ReplyCount,
				LikeCount:   commentItem.LikeCount,
				CommentTime: commentItem.CommentTime,
			}
			comments = append(comments, targetItem)
		}
	}

	videoIDInt, _ := strconv.ParseUint(videoID, 10, 64)
	dynamicMap, err := getCommentDynamicInfoWithVideoID(cidsHash, videoIDInt)
	for _, comment := range comments {
		if item, ok := dynamicMap[comment.CommentId]; ok {
			comment.ReplyCount = item.ReplyNum
			comment.LikeCount = item.FavorNum
		}
	}

	resp.Status.Code = pb.Status_OK
	resp.TopicItem = topicItem
	resp.ArticleItem = articleItem
	resp.VideoItem = videoItem
	resp.CommentItems = comments
	resp.DownloadUrl = common.GetConfig().DownLoadURL
	glog.Debug("################################")
	glog.Debug("===============ShareVideoBottomPage cost:%v", time.Now().Sub(begin))

	return
}

// ShareTopicPage share topic page
func (s *Server) ShareTopicPage(ctx context.Context, req *pb.ShareTopicPageRequest) (resp *pb.ShareTopicPageResponse, err error) {
	resp = &pb.ShareTopicPageResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	begin := time.Now()
	topicID := req.GetTopicId()
	if topicID == "" {
		resp.Status.Code = pb.Status_OK
		resp.Status.Message = "request parameter topicID is empty"
		return
	}

	// 1. get topicItem
	toppicData, err := db.HGet(common.TOPICHASH, topicID)
	if err != nil {
		glog.Error("get topicItem from redis failed. err:%v", err)
		return
	}
	topicItem := &pb.TopicItem{}
	if err := proto.Unmarshal([]byte(toppicData), topicItem); err != nil {
		glog.Error("Unmarshal pb message failed. topicID:%s, err:%v", topicID, err)
	}

	// 2. get videoID under the topicID
	topicVideoString := fmt.Sprintf("%s%s", common.TOPICWITHVIDEOID, topicID)
	vidsHash, err := db.ZRevRange(topicVideoString, 0, common.GetConfig().TopicPage.VideoNUM-1)
	if err != nil {
		glog.Error("get vids from redis failed. err:%v", err)
		return
	}

	// 3. get listItem from redis
	var hashListItems []string
	var listItems []*pb.ListItem
	if len(vidsHash) != 0 {
		glog.Debug("videonum:%d", len(vidsHash))
		vidsArr := common.TransStrArrToInterface(vidsHash)
		hashListItems, err = db.HMGet(common.FULLVIDEOHASHKEYREDIS, vidsArr)
		if err != nil {
			glog.Error("get listItems from redis failed. err:%v", err)
			return
		}

		for _, hashListItem := range hashListItems {
			if hashListItem == "" {
				glog.Error("listItem is empty in redis")
				continue
			}

			listItem := &pb.ListItem{}
			if err := proto.Unmarshal([]byte(hashListItem), listItem); err != nil {
				glog.Error("Unmarshal pb message failed. err:%v", err)
				continue
			}

			videoItem := &pb.VideoItem{
				VideoId:    listItem.VideoItem.VideoId,
				PictureUrl: listItem.VideoItem.PictureUrl,
				Title:      listItem.VideoItem.Title,
				Width:      listItem.VideoItem.Width,
				Height:     listItem.VideoItem.Height,
				Duration:   listItem.VideoItem.Duration,
			}
			articleItem := &pb.ArticleItem{
				LikeCount:    listItem.ArticleItem.LikeCount,
				CommentCount: listItem.ArticleItem.CommentCount,
			}

			targetListItem := &pb.ListItem{
				VideoItem:   videoItem,
				ArticleItem: articleItem,
			}

			listItems = append(listItems, targetListItem)
		}
	}

	resp.Status.Code = pb.Status_OK
	resp.TopicItem = topicItem
	resp.ListItems = listItems
	resp.DownloadUrl = common.GetConfig().DownLoadURL
	glog.Debug("===============ShareTopicPage cost:%v", time.Now().Sub(begin))

	return
}

func getCommentDynamicInfoWithVideoID(commentIds []string, videoID uint64) (dynamicMap map[string]*CommentDynamicInfo, err error) {
	key := fmt.Sprintf("%s%d", COMMENTDYNAMICPREFIX, videoID%1000)
	var fields []interface{}
	var field string
	for _, cid := range commentIds {
		field = fmt.Sprintf("%s%s", FAVORNUMPREFIX, cid)
		fields = append(fields, field)
		field = fmt.Sprintf("%s%s", REPLYNUMPREFIX, cid)
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

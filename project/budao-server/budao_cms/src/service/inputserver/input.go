package inputserver

import (
	"common"
	"context"
	"db"
	"fmt"
	"github.com/sumaig/glog"
	"math/rand"
	"service"
	"strings"
	"time"
	pb "twirprpc"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) InputVideoData(ctx context.Context, req *pb.InputVideoRequest) (resp *pb.InputVideoResponse, err error) {
	resp = &pb.InputVideoResponse{
		Status: common.GetInitStatus(),
		MapId:  make(map[string]string),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
		}
	}()

	if len(req.VideoDatas) <= 0 {
		err = fmt.Errorf("request data is null")
		return
	}

	//插入数据
	rs := rand.NewSource(int64(time.Now().Nanosecond()))
	r := rand.New(rs)
	vsids := make([]string, 0, len(req.VideoDatas))

	for _, video := range req.VideoDatas {
		//判断视频尺寸是否正确
		if video.VideoWidth <= 0 || video.VideoHeight <= 0 {
			continue
		}

		var i_id string
		i_id = fmt.Sprintf("%v", time.Now().Nanosecond()+r.Intn(100000000000))
		v_source_id := fmt.Sprintf("%v", int32(video.Source)) + "_" + video.Vid
		vsids = append(vsids, v_source_id)
		CreateTime := time.Unix(int64(video.CreateTime), 0).Format("2006-01-02 15:04:05")

		execSql := fmt.Sprintf(`insert into video_data (channel_id,topic,question_type,meta_data,play_url,video_id,source_type,i_id,media_name,media_id,video_title,play_count,video_duration,share_url,video_cover,source,video_width,video_height,v_source_id, praise_count,fav_count,share_count,comment_count,create_time,parse_type) value('%v','%v','%v','%v','%v','%v','%v',%v,'%v','%v','%v',%v,%v,'%v','%v',%v,%v,%v,'%v',%v,%v,%v,%v,'%v',%v)
ON DUPLICATE KEY UPDATE  %s`,
			video.ChannelId, video.Topic, video.QuestionType, video.MetaData, video.VideoUrl, video.Vid, video.SourceType, i_id, video.MediaName, video.MediaId, video.VideoTitle, video.PlayCount, video.VideoDuration, video.ShareUrl, video.VideoCover, int32(video.Source), video.VideoWidth, video.VideoHeight, v_source_id, video.PraiseCount, video.FavCount, video.ShareCount, video.CommentCount, CreateTime, int32(video.ParseType),
			getVideoUpdateSql(*video),
		)
		fmt.Println(execSql)
		db.Exec(service.SPIDERDB, execSql)
	}

	querySql := fmt.Sprintf(`select i_id,video_id from video_data where v_source_id in ("%v")`, strings.Join(vsids, `","`))
	rows, err := db.Query(service.SPIDERDB, querySql)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		return
	}
	for rows.Next() {
		var i_id, vid string
		rows.Scan(&i_id, &vid)
		resp.MapId[vid] = i_id
	}

	return
}

func (s *Server) InputCommentData(ctx context.Context, req *pb.InputCommentRequest) (resp *pb.InputCommentResponse, err error) {
	resp = &pb.InputCommentResponse{
		Status: common.GetInitStatus(),
		MapId:  make(map[string]string),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
		}
	}()

	if len(req.CommentDatas) <= 0 {
		err = fmt.Errorf("request data is null")
		return
	}
	vsids := make([]string, 0, len(req.CommentDatas))
	cids := make([]string, 0, len(req.CommentDatas))

	for _, comment := range req.CommentDatas {
		v_source_id := fmt.Sprintf("%v", int32(comment.Source)) + "_" + comment.Vid
		CreateTime := time.Unix(int64(comment.CreateTime), 0).Format("2006-01-02 15:04:05")
		comment.Content = strings.Replace(comment.Content, `"`, `\"`, -1)
		comment.UserName = strings.Replace(comment.UserName, `'`, `\'`, -1)
		var isHot = 0
		if comment.IsHot {
			isHot = 1
		}
		vsids = append(vsids, v_source_id)
		cids = append(cids, comment.Cid)

		execSql := fmt.Sprintf(`insert into comment_data (v_source_id,vid,cid,source,content,favor_num,user_id,user_name,user_photo,reply_num,is_hot,create_time) value ('%v','%v','%v',%v,"%v",%v,'%v','%v','%v',%v,%v,'%v') 
ON DUPLICATE KEY UPDATE content="%v",favor_num=%v,user_id='%v',user_name='%v',user_photo='%v',reply_num=%v,is_hot=%v,create_time='%v';`,
			v_source_id, comment.Vid, comment.Cid, int32(comment.Source), comment.Content, comment.FavorNum, comment.UserId, comment.UserName, comment.UserPhoto, comment.ReplyNum, isHot, CreateTime,
			comment.Content, comment.FavorNum, comment.UserId, comment.UserName, comment.UserPhoto, comment.ReplyNum, isHot, CreateTime)
		db.Exec(service.SPIDERDB, execSql)
	}
	if err != nil {
		glog.Error(err)
		return
	}

	querySql := fmt.Sprintf(`select id,cid from comment_data where v_source_id in ("%v") and cid in ("%v") `, strings.Join(vsids, `","`), strings.Join(cids, `","`))
	rows, err := db.Query(service.SPIDERDB, querySql)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		return
	}
	for rows.Next() {
		var id, cid string
		rows.Scan(&id, &cid)
		resp.MapId[cid] = id
	}

	return
}

func (s *Server) InputAuditOfflineVideoData(ctx context.Context, req *pb.InputAuditVideoRequest) (resp *pb.InputAuditVideoResponse, err error) {
	resp = &pb.InputAuditVideoResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
		}
	}()
	now := time.Now().Format(`2006-01-02 15:04:05`)
	if len(req.VideoDatas) <= 0 {
		err = fmt.Errorf("request data is null")
		return
	}
	for _, video := range req.VideoDatas {
		exclSql := fmt.Sprintf(`insert into audit_to_offline (vid,vsource_vid) value (%v,'%v') ON DUPLICATE KEY UPDATE insert_time='%v'`, video.PostVid, video.VsourceVid, now)
		_, err = db.Exec(service.BUDAODB, exclSql)
		if err != nil {
			return
		}
	}

	return
}

func getVideoUpdateSql(data pb.VideoData) string {
	var updataSql string
	if data.VideoUrl != "" {
		updataSql += fmt.Sprintf(`,play_url='%v'`, data.VideoUrl)
	}

	if data.SourceType != "" {
		updataSql += fmt.Sprintf(`,source_type='%v'`, data.SourceType)
	}

	if data.MediaName != "" {
		updataSql += fmt.Sprintf(`,media_name='%v'`, data.MediaName)
	}

	if data.MediaId != "" {
		updataSql += fmt.Sprintf(`,media_id='%v'`, data.MediaId)
	}

	if data.VideoTitle != "" {
		updataSql += fmt.Sprintf(`,video_title='%v'`, data.VideoTitle)
	}

	if data.PlayCount > 0 {
		updataSql += fmt.Sprintf(`,play_count=%v`, data.PlayCount)
	}

	if data.VideoDuration > 0 {
		updataSql += fmt.Sprintf(`,video_duration=%v`, data.VideoDuration)
	}

	if data.ShareUrl != "" {
		updataSql += fmt.Sprintf(`,share_url='%v'`, data.ShareUrl)
	}

	if data.VideoCover != "" {
		updataSql += fmt.Sprintf(`,video_cover='%v'`, data.VideoCover)
	}

	if data.VideoWidth > 0 {
		updataSql += fmt.Sprintf(`,video_width=%v`, data.VideoWidth)
	}

	if data.VideoHeight > 0 {
		updataSql += fmt.Sprintf(`,video_height=%v`, data.VideoHeight)
	}

	if data.PraiseCount > 0 {
		updataSql += fmt.Sprintf(`,praise_count=%v`, data.PraiseCount)
	}
	if data.FavCount > 0 {
		updataSql += fmt.Sprintf(`,fav_count=%v`, data.FavCount)
	}
	if data.ShareCount > 0 {
		updataSql += fmt.Sprintf(`,share_count=%v`, data.ShareCount)
	}

	if data.ParseType != 99 {
		updataSql += fmt.Sprintf(`,parse_type=%v`, int32(data.ParseType))
	}

	if data.CommentCount > 0 {
		updataSql += fmt.Sprintf(`,comment_count=%v`, data.CommentCount)
	}
	if data.CreateTime > 0 {
		CreateTime := time.Unix(int64(data.CreateTime), 0).Format("2006-01-02 15:04:05")
		updataSql += fmt.Sprintf(`,create_time='%v'`, CreateTime)
	}

	//channel_id,topic,question_type,meta_data
	if data.ChannelId != "" {
		updataSql += fmt.Sprintf(`,channel_id='%v'`, data.ChannelId)
	}
	if data.Topic != "" {
		updataSql += fmt.Sprintf(`,topic='%v'`, data.Topic)
	}
	if data.QuestionType != "" {
		updataSql += fmt.Sprintf(`,question_type='%v'`, data.QuestionType)
	}
	if data.MetaData != "" {
		updataSql += fmt.Sprintf(`,meta_data='%v'`, data.MetaData)
	}

	return updataSql[1:]
}

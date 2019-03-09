package spiderservice

import (
	"common"
	"context"
	"db"
	"fmt"
	"github.com/sumaig/glog"
	"github.com/twitchtv/twirp"
	"log"
	"net/http"
	"service"
	"service/api"
	"strconv"
	pb "twirprpc"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) VideoList(ctx context.Context, req *api.QueryListRequest) (resp *api.VideoListResponse, err error) {
	resp = &api.VideoListResponse{
		Data: &api.VideoList{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	var m = map[string]string{
		"status":       `status,=%v`,
		"playCount":    `play_count,=%v`,
		"title":        `video_title,like '%%%v%%'`,
		"fromTime":     `insert_time,>='%v'`,
		"toTime":       `insert_time,<='%v'`,
		"source":       "source,=%v",
		"sourceType":   `source_type,like '%%%v%%'`,
		"topic":        "topic,like '%%%v%%'",
		"durationS":    "video_duration,>=%v",
		"durationE":    "video_duration,<=%v",
		"titleLenS":    "char_length(video_title),>=%v",
		"mediaName":    "media_name,like '%%%v%%'",
		"praiseCount":  "praise_count,>=%v",
		"favCount":     "fav_count,>=%v",
		"shareCount":   "share_count,>=%v",
		"commentCount": "comment_count,>=%v",
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		return
	}

	querySql := fmt.Sprintf(`select ifnull(insert_time,""), ifnull(play_url,""),ifnull(source_type,""),ifnull(post_vid,""),ifnull(v_source_id,""),ifnull(media_name,""),praise_count,fav_count,share_count,comment_count,char_length(video_title) as title_len,ifnull(video_duration,0),topic,id,video_title,play_count,share_url as video_url,video_cover,i_id,status,source,video_width,video_height, parse_type from video_data where 1=1 %v %v %v`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.SPIDERDB, querySql)

	if err != nil {
		err = service.MysqlError
		return
	}
	defer rows.Close()
	data := make([]*api.Video, 0, 10)
	for rows.Next() {
		temp := api.Video{}
		err := rows.Scan(&temp.InsertTime, &temp.PlayUrl, &temp.SourceType, &temp.PostVid, &temp.VSourceId, &temp.MediaName, &temp.PraiseCount, &temp.FavCount, &temp.ShareCount, &temp.CommentCount, &temp.TitleLen, &temp.Duration, &temp.Topic, &temp.Id, &temp.Title, &temp.PlayCount, &temp.VideoUrl, &temp.VideoCover, &temp.Vid, &temp.Status, &temp.Source, &temp.VideoWidth, &temp.VideoHeight, &temp.ParseType)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			break
		}
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.SPIDERDB, fmt.Sprintf(`select id from video_data where 1=1 %s`, filterSql))
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = count
	return
}

func (s *Server) Modify(ctx context.Context, req *api.ModifyResponse) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	execSql := fmt.Sprintf(`UPDATE video_data set video_title='%s' WHERE i_id=%s`, req.NewTitle, req.Id)
	_, err = db.Exec(service.SPIDERDB, execSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	resp.Code = "200"
	resp.Msg = "OK"
	return
}

func (s *Server) GiveUP(ctx context.Context, req *api.Identity) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	execSql := fmt.Sprintf(`UPDATE video_data set status=2 WHERE id=%s`, req.Id)
	_, err = db.Exec(service.SPIDERDB, execSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	resp.Code = "200"
	resp.Msg = "OK"
	return
}

func (s *Server) Post(ctx context.Context, req *api.Video) (resp *api.PostVideResponse, err error) {
	resp = &api.PostVideResponse{}
	var (
		rpcReq  *pb.PostVideosRequest
		rpcResp *pb.PostVideosResponse
		rpcErr  error
	)
	if req.VideoHeight == 0 || req.VideoWidth == 0 || req.Duration == 0 {
		//更新video_data发布失败原因
		//updateVideoData(param.Vid, "video height or weight or duration 0", 3)
		err = service.ParamError
		return
	}
	rpcReq = &pb.PostVideosRequest{}
	tVideos := make([]*pb.PostVideo, 0)
	postVideo := &pb.PostVideo{}

	if len(req.VSourceId) <= 0 {
		//更新video_data发布失败原因
		//updateVideoData(param.Vid, "v_source_id is empty = "+ param.VSourceID, 3)
		err = fmt.Errorf("v_source_id is empty")
		return
	}
	postVideo.VsourceVid = req.VSourceId
	postVideo.SourceVid = req.Vid
	postVideo.Pic = req.VideoCover
	//postVideo.VideoUrl = param.VideoUrl
	postVideo.Title = req.Title
	postVideo.Width = req.VideoWidth
	postVideo.Height = req.VideoHeight
	postVideo.Duration = req.Duration
	postVideo.EVideoSource = pb.EVideoSource(req.Source)
	postVideo.EVideoParseRule = pb.VideoParseRule(req.ParseType)

	if len(req.PlayUrl) > 0 {
		postVideo.VideoUrl = req.PlayUrl
	} else {
		postVideo.VideoUrl = req.VideoUrl
	}

	tVideos = append(tVideos, postVideo)
	rpcReq.TVideos = tVideos

	client := pb.NewTransferProtobufClient(common.GetConfig().Extension["transferUrl"], &http.Client{})

	rpcResp, rpcErr = client.PostVideos(context.Background(), rpcReq)
	if rpcErr != nil {
		//更新video_data发布失败原因
		//updateVideoData(param.Vid, fmt.Sprintf("rpc return err: %v", rpcErr), 3)
		if twerr, ok := rpcErr.(twirp.Error); ok {
			if twerr.Meta("retryable") != "" {
				// Log the error and go again.
				log.Printf("got error %q, retrying", twerr)
			}
		}
		log.Fatal(err)
		err = rpcErr
		return
	}
	if len(rpcResp.GetTResults()) > 0 {
		result := rpcResp.GetTResults()[0]
		data := &api.Video{
			Vid: result.Vid,
		}
		if data.Vid != "" {
			resp.Code = "200"
			resp.Msg = "OK"
			resp.Data = data
		} else {
			//更新video_data发布失败原因
			//updateVideoData(param.Vid, fmt.Sprintf("retun result: %d", pb.EPostResult(result.Result)), 3)
			err = service.HasPostError
		}
		updateSql := fmt.Sprintf(`update video_data set status=1,post_vid=%v where i_id=%v`, data.Vid, req.Vid)
		db.Exec(service.SPIDERDB, updateSql)
	} else {
		//更新video_data发布失败原因
		//updateVideoData(param.Vid, "server not return value", 3)
		err = fmt.Errorf("server not return value")
	}
	return
}

func (s *Server) VideoCommentList(ctx context.Context, req *api.QueryListRequest) (resp *api.CommentListResponse, err error) {
	resp = &api.CommentListResponse{
		Data: &api.CommentList{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	var m = map[string]string{
		"status":    `status,=%v`,
		"vSourceId": `v_source_id,='%v'`,
		"favorNum":  `favor_num,>=%v`,
		"replyNum":  `reply_num,>=%v`,
		"ishot":     `is_hot,=%v`,
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	querySql := fmt.Sprintf(`select cid,source,user_id,v_source_id,status,content,favor_num,user_name,user_photo,reply_num,is_hot,create_time from comment_data where 1=1 %v %v %v`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.SPIDERDB, querySql)
	if err != nil {
		err = service.MysqlError
		return
	}
	defer rows.Close()
	data := make([]*api.Comment, 0, 10)
	type result struct {
		Data  []*api.Comment `json:"data"`
		Total uint64         `json:"total"`
	}

	for rows.Next() {
		temp := api.Comment{}
		err := rows.Scan(&temp.Cid, &temp.Source, &temp.Uid, &temp.VSourceId, &temp.Status, &temp.Content, &temp.FavorNum, &temp.UName, &temp.UPhoto, &temp.ReplyNum, &temp.IsHot, &temp.CTime)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}
		temp.CTime = common.GetTimeStr(temp.CTime)
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.SPIDERDB, fmt.Sprintf(`select id from comment_data where 1=1 %v`, filterSql))
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = count
	return
}

func (s *Server) ModifyCommentStatus(ctx context.Context, req *api.ModifyCommentStatusRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	for _, comment := range req.Comments {
		db.Exec(service.SPIDERDB, fmt.Sprintf(`update comment_data set status=%v where v_source_id='%v' and cid='%v'`, req.Status, comment.VSourceId, comment.Cid))
	}
	return
}

func (s *Server) PostComment(ctx context.Context, req *api.PostCommentRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	insertNum := 0
	for _, comment := range req.Comments {
		cid, autoIncreId, tableName, err := common.GetItemId(service.BUDAODB, service.COMMENT_)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}
		u := service.GetRandUserByHash(comment.Source + "_" + comment.Uid)
		updateSql := fmt.Sprintf(`update %v set content="%v",cid=%v,vid=%v,from_uid=%v,from_name="%v",from_photo="%v",create_time='%v',favor_num=%v,reply_num=%v where id=%v`, tableName, comment.Content, cid, req.Vid, u.Uid, u.Name, u.Photo, comment.CTime, comment.FavorNum, comment.ReplyNum, autoIncreId)
		_, err = db.Exec(service.BUDAODB, updateSql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}

		insertNum++

		db.Exec(service.SPIDERDB, fmt.Sprintf(`update comment_data set status=1,post_cid=%v where v_source_id='%v' and cid='%v'`, cid, comment.VSourceId, comment.Cid))

		uid, _ := strconv.ParseUint(u.Uid, 10, 64)
		tableName, _ = db.GetTableName(service.USER_, uid)
		db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+1 where uid=%v`, tableName, uid))
	}

	vid, _ := strconv.ParseUint(req.Vid, 10, 64)
	tableName, _ := db.GetTableName(service.VIDEO_, vid)
	db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+%v where vid =%v`, tableName, insertNum, vid))

	return
}

func (s *Server) PostAllComment(ctx context.Context, req *api.PostAllCommentReques) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	cidSet := make(map[string]bool)
	for _, cid := range req.Cids {
		cidSet[cid] = true
	}
	//查询出所以评论
	querySql := fmt.Sprintf(`select  v_source_id,cid,source,content,favor_num,user_id,reply_num,create_time from comment_data where v_source_id='%v' and status=0`, req.VSourceId)
	rows, err := db.Query(service.SPIDERDB, querySql)
	defer rows.Close()
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	insertNum := 0
	for rows.Next() {
		temp := api.Comment{}
		rows.Scan(&temp.VSourceId, &temp.Cid, &temp.Source, &temp.Content, &temp.FavorNum, &temp.Uid, &temp.ReplyNum, &temp.CTime)
		if cidSet[temp.Cid] {
			continue
		}

		temp.CTime = common.GetTimeStr(temp.CTime)
		cid, autoIncreId, tableName, err := common.GetItemId(service.BUDAODB, service.COMMENT_)
		u := service.GetRandUserByHash(temp.Source + "_" + temp.Uid)

		updateSql := fmt.Sprintf(`update %v set content="%v",cid=%v,vid=%v,from_uid=%v,from_name="%v",from_photo="%v",create_time='%v',favor_num=%v where id=%v`, tableName, temp.Content, cid, req.Vid, u.Uid, u.Name, u.Photo, temp.CTime, temp.FavorNum, autoIncreId)
		_, err = db.Exec(service.BUDAODB, updateSql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}

		insertNum++

		db.Exec(service.SPIDERDB, fmt.Sprintf(`update comment_data set status=1,post_cid=%v where v_source_id='%v' and cid='%v'`, cid, temp.VSourceId, temp.Cid))

		uid, _ := strconv.ParseUint(u.Uid, 10, 64)
		tableName, _ = db.GetTableName(service.USER_, uid)
		db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+1 where uid=%v`, tableName, uid))
	}

	vid, _ := strconv.ParseUint(req.Vid, 10, 64)
	tableName, _ := db.GetTableName(service.VIDEO_, vid)
	db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+%v where vid =%v`, tableName, insertNum, vid))

	return
}

package userservice

import (
	"common"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"github.com/sumaig/glog"
	"service"
	"service/api"
	"strconv"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

//评论举报操作
func (s *Server) ReportCommentOpera(ctx context.Context, req *api.ReportCommentOperaRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	maxNum := common.GetConfig().DB.MySQL[service.BUDAODB].TableDesc[service.COMMENT_]
	switch req.Type {
	case "del":
		for i := uint64(0); i < maxNum; i++ {
			tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
			querySql := fmt.Sprintf(`update %v set state=4 where cid = '%v'`, tableName, req.Cid)
			_, err = db.Exec(service.BUDAODB, querySql)
			if err != nil {
				glog.Error(err)
				err = service.ParamError
				return
			}
		}

		sql := fmt.Sprintf(`update user_report_comment set state=1 where cid = '%v' and id = '%v'`, req.Cid, req.Id)
		_, err = db.Exec(service.BUDAODB, sql)
		if err != nil {
			glog.Error(err)
			err = service.ParamError
			return
		}
	case "retain":
		rsql := fmt.Sprintf(`select uid, cid, state from user_report_comment where id=%v`, req.Id)
		row, err := db.QueryRow(service.BUDAODB, rsql)
		if err != nil {
			glog.Error(err)
			return nil, service.ParamError
		}
		report := service.UserReportComment{}
		err = row.Scan(&report.Uid, &report.Cid, &report.Status)
		if err != nil {
			glog.Error(err)
			return nil, service.MysqlError
		}

		if report.Status == "1" {
			for i := uint64(0); i < maxNum; i++ {
				tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
				querySql := fmt.Sprintf(`update %v set state=2 where cid = '%v'`, tableName, req.Cid)
				_, err = db.Exec(service.BUDAODB, querySql)
				if err != nil {
					glog.Error(err)
					return nil, service.ParamError
				}
			}

			sql := fmt.Sprintf(`update user_report_comment set state=0 where cid='%v' and id = '%v'`, req.Cid, req.Id)
			_, err = db.Exec(service.BUDAODB, sql)
			if err != nil {
				glog.Error(err)
				return nil, service.MysqlError
			}
		} else if report.Status == "2" {
			rsql := fmt.Sprintf(`update user_report_comment set state=0 where cid='%v' and id = '%v'`, req.Cid, req.Id)
			_, err = db.Exec(service.BUDAODB, rsql)
			if err != nil {
				glog.Error(err)
				return nil, service.ParamError
			}
		}
	case "report":
		sql := fmt.Sprintf(`update user_report_comment set state=2 where cid='%v' and id = '%v'`, req.Cid, req.Id)
		_, err = db.Exec(service.BUDAODB, sql)
		if err != nil {
			glog.Error(err)
			return nil, service.MysqlError
		}
	default:
		err = fmt.Errorf("type值不存在")
		return
	}
	return
}

//评论举报
func (s *Server) ReportComment(ctx context.Context, req *api.QueryListRequest) (resp *api.ReportCommentResponse, err error) {
	resp = &api.ReportCommentResponse{
		Data: &api.ReportCommentList{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	var m = map[string]string{
		"uid":          `uid,='%v'`,
		"reportReason": `reason,='%v'`,
		"startTime":    `update_time,>='%v'`,
		"endTime":      `update_time,<='%v'`,
		"operaStatus":  `state,=%v`,
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	//对条件是vid的进行处理
	if len(req.Filter["vid"]) > 0 {
		res, count := getCommentCids(req.Filter["vid"], filterSql)
		resp.Code = service.SUCCESS_CODE
		resp.Msg = "OK"
		resp.Data.Data = res
		resp.Data.Count = count
		return
	}

	querySql := fmt.Sprintf(`select id, uid, device_id, cid, reason, state, update_time from user_report_comment where 1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()
	data := make([]*api.ReportCommentMessage, 0, 10)

	for rows.Next() {
		temp := api.ReportCommentMessage{}
		err := rows.Scan(&temp.Id, &temp.Uid, &temp.DeviceId, &temp.Cid, &temp.ReportReason, &temp.OperaStatus, &temp.ReportTime)
		if err != nil {
			err = service.MysqlError
			break
		}
		//时间格式转化
		temp.ReportTime = common.GetTimeStr(temp.ReportTime)

		//查询评论表
		maxNum := common.GetConfig().DB.MySQL[service.BUDAODB].TableDesc[service.COMMENT_]
		for i := uint64(0); i < maxNum; i++ {
			tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
			commSql := fmt.Sprintf(`select vid, content from %v where cid = '%v'`, tableName, temp.Cid)
			row2, _ := db.QueryRow(service.BUDAODB, commSql)
			glog.Info(row2.Scan(&temp.Vid, &temp.Content))

		}

		//根据vid获取视频信息
		vSql := fmt.Sprintf(`select coverurl, videourl from video_0 where vid = '%v'`, temp.Vid)
		row, err := db.QueryRow(service.BUDAODB, vSql)
		if err != nil {
			glog.Error(err)
			return nil, service.MysqlError
		}
		row.Scan(&temp.CoverUrl, &temp.VideoUrl)

		//根据uid获取用户信息
		uSql := fmt.Sprintf(`select name, photo, phone from user_0 where uid = '%v'`, temp.Uid)
		row1, err := db.QueryRow(service.BUDAODB, uSql)
		if err != nil {
			glog.Error(err)
			return nil, service.MysqlError
		}
		row1.Scan(&temp.Name, &temp.Photo, &temp.Phone)

		data = append(data, &temp)
	}

	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from user_report_comment where 1=1 %s`, filterSql))
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = count
	return
}

//根据vid获取所有的cid
func getCommentCids(vid string, filter string) ([]*api.ReportCommentMessage, uint64) {
	result := make([]*api.ReportCommentMessage, 0, 10)
	var total uint64
	maxNum := common.GetConfig().DB.MySQL[service.BUDAODB].TableDesc[service.COMMENT_]
	for i := uint64(0); i < maxNum; i++ {
		tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
		uSql := fmt.Sprintf(`select cid, vid, content from %v where vid = '%v'`, tableName, vid)
		rows, _ := db.Query(service.BUDAODB, uSql)
		defer rows.Close()
		for rows.Next() {
			temp := api.ReportCommentMessage{}
			rows.Scan(&temp.Cid, &temp.Vid, &temp.Content)

			//根据cid获取report_comment
			commInfos, count := getCidCommentInfo(temp.Cid, filter)
			glog.Info(commInfos)
			for _, commInfo := range commInfos {
				if len(commInfo.Id) > 0 {
					//获取视频信息
					vSql := fmt.Sprintf(`select coverurl, videourl from video_0 where vid = '%v'`, temp.Vid)
					row, _ := db.QueryRow(service.BUDAODB, vSql)

					row.Scan(&temp.CoverUrl, &temp.VideoUrl)

					commInfo.CoverUrl = temp.CoverUrl
					commInfo.VideoUrl = temp.VideoUrl

					//评论内容
					commInfo.Content = temp.Content
					commInfo.Vid = temp.Vid

					//根据uid获取用户信息
					uSql := fmt.Sprintf(`select name, photo, phone from user_0 where uid = '%v'`, commInfo.Uid)
					row1, _ := db.QueryRow(service.BUDAODB, uSql)

					row1.Scan(&temp.Name, &temp.Photo, &temp.Phone)

					commInfo.Name = temp.Name
					commInfo.Photo = temp.Photo
					commInfo.Phone = temp.Phone

					result = append(result, commInfo)
					total = count
				}
			}
		}
	}

	return result, total
}

//根据cid获取举报信息
func getCidCommentInfo(cid string, filter string) ([]*api.ReportCommentMessage, uint64) {
	cQuery := fmt.Sprintf("select id, uid, device_id, cid, reason, state, update_time from user_report_comment where cid = %v and 1=1 %v", cid, filter)
	rows, _ := db.Query(service.BUDAODB, cQuery)
	defer rows.Close()
	res := make([]*api.ReportCommentMessage, 0, 1)
	for rows.Next() {
		temp := api.ReportCommentMessage{}
		rows.Scan(&temp.Id, &temp.Uid, &temp.DeviceId, &temp.Cid, &temp.ReportReason, &temp.OperaStatus, &temp.ReportTime)
		temp.ReportTime = common.GetTimeStr(temp.ReportTime)
		res = append(res, &temp)
	}

	count, _ := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from user_report_comment where 1=1 %s`, filter))

	return res, count
}

//反馈列表
func (s *Server) FeedbackList(ctx context.Context, req *api.QueryListRequest) (resp *api.FeedbackListResponse, err error) {
	resp = &api.FeedbackListResponse{
		Data: &api.FeedbackList{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	var m = map[string]string{
		"startTime": `update_time,>='%v'`,
		"endTime":   "update_time,<='%v'",
	}

	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	uuid := make(map[string]string)
	if len(req.Filter["mobile"]) > 0 {
		//根据手机号获取uid
		uuid = getUserMobileInfo(req.Filter["mobile"])
		filterSql = filterSql + " and uid = '" + uuid["uid"] + "'"
	}

	querySql := fmt.Sprintf(`select uid, feedback, update_time from user_feedback where 1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()

	data := make([]*api.FeedbackMessage, 0, 10)

	for rows.Next() {
		temp := api.FeedbackMessage{}
		err := rows.Scan(&temp.Uid, &temp.FeedBack, &temp.FCreateTime)

		//获取用户表的名
		tID, _ := strconv.ParseUint(temp.Uid, 10, 64)
		tableName, err := db.GetTableName(service.USER_, tID)
		glog.Error(err)
		uInfo := getUserInfo(tableName, temp.Uid)
		uInfo.FCreateTime = common.GetTimeStr(uInfo.FCreateTime)

		data = append(data, &uInfo)
		glog.Info(data)
	}

	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select uid, feedback from user_feedback where 1=1 %s`, filterSql))
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = count
	return
}

//根据uid获取用户信息
func getUserInfo(tableName, uid string) api.FeedbackMessage {
	uSql := fmt.Sprintf(`select name, photo, phone from %v where uid = '%v'`, tableName, uid)
	row, _ := db.QueryRow(service.BUDAODB, uSql)
	temp := api.FeedbackMessage{}
	row.Scan(&temp.Name, &temp.Photo, &temp.Phone)

	return temp
}

//根据手机号获取uid
func getUserMobileInfo(phone string) map[string]string {
	uSql := fmt.Sprintf(`select uid from user_0 where phone = '%v'`, phone)
	row, err := db.QueryRow(service.BUDAODB, uSql)
	if err != nil {
		return nil
	}

	m := make(map[string]string)
	type user struct {
		Uid string `json:"uid"`
	}
	temp := &user{}
	row.Scan(&temp.Uid)

	j, _ := json.Marshal(temp)
	json.Unmarshal(j, &m)

	return m
}

//举报视频操作
func (s *Server) ReportVideoOpera(ctx context.Context, req *api.ReportVideoOperaRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	switch req.Type {
	case "del":
		querySql := fmt.Sprintf(`update video_0 set state=4 where vid=%v`, req.Vid)
		_, err = db.Exec(service.BUDAODB, querySql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}

		sql := fmt.Sprintf(`update user_report_video set state=1 where vid=%v`, req.Vid)
		_, err = db.Exec(service.BUDAODB, sql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
	case "retain":
		rsql := fmt.Sprintf(`select uid, vid, state from user_report_video where vid=%v`, req.Vid)
		row, err := db.QueryRow(service.BUDAODB, rsql)
		if err != nil {
			glog.Error(err)
			return nil, service.MysqlError
		}
		report := service.UserReportVideo{}
		err = row.Scan(&report.Uid, &report.Vid, &report.Status)
		if err != nil {
			glog.Error(err)
			return nil, service.MysqlError
		}

		if report.Status == "1" {
			querySql := fmt.Sprintf(`update video_0 set state=2 where vid=%v`, req.Vid)
			_, err = db.Exec(service.BUDAODB, querySql)
			if err != nil {
				glog.Error(err)
				return nil, service.MysqlError
			}

			sql := fmt.Sprintf(`update user_report_video set state=0 where vid=%v`, req.Vid)
			_, err = db.Exec(service.BUDAODB, sql)
			if err != nil {
				glog.Error(err)
				return nil, service.MysqlError
			}
		} else if report.Status == "2" {
			rsql := fmt.Sprintf(`update user_report_video set state=0 where vid=%v`, req.Vid)
			_, err = db.Exec(service.BUDAODB, rsql)
			if err != nil {
				glog.Error(err)
				return nil, service.MysqlError
			}
		}
	case "report":
		sql := fmt.Sprintf(`update user_report_video set state=2 where vid=%v`, req.Vid)
		_, err = db.Exec(service.BUDAODB, sql)
		if err != nil {
			glog.Error(err)
			return nil, service.MysqlError
		}
	default:
		err = fmt.Errorf("type值不存在")
		return
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"

	return
}

//视频举报
func (s *Server) ReportVideo(ctx context.Context, req *api.QueryListRequest) (resp *api.ReportVideoResponse, err error) {
	resp = &api.ReportVideoResponse{
		Data: &api.ReportVideoList{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	var m = map[string]string{
		"uid":          `u.uid,='%v'`,
		"vid":          `v.vid,='%v'`,
		"reportReason": `u.reason,='%v'`,
		"startTime":    `u.update_time,>='%v'`,
		"endTime":      "u.update_time,<='%v'",
		"operaStatus":  `u.state,=%v`,
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	querySql := fmt.Sprintf(`select v.vid, v.title, v.coverurl, v.videourl, v.create_time, v.view_num, v.favor_num, v.comment_num, v.duration, v.type, v.state ,u.uid, u.reason, u.update_time, u.state from video_0 as v inner join user_report_video as u on v.vid = u.vid where 1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)

	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()
	data := make([]*api.ReportVideoMessage, 0, 10)
	for rows.Next() {
		temp := api.ReportVideoMessage{}
		err := rows.Scan(&temp.Vid, &temp.Title, &temp.VideoCover, &temp.VideoUrl, &temp.VCreateTime, &temp.ViewNum, &temp.FavorNum, &temp.CommentNum, &temp.Duration, &temp.Source, &temp.Status, &temp.Uid, &temp.Reason, &temp.ReportTime, &temp.OperaStatus)
		if err != nil {
			glog.Info(err)
			err = service.MysqlError
			break
		}
		//时间格式转化
		temp.VCreateTime = common.GetTimeStr(temp.VCreateTime)
		temp.ReportTime = common.GetTimeStr(temp.ReportTime)

		data = append(data, &temp)
	}
	csql := fmt.Sprintf(`select v.vid from video_0 as v inner join user_report_video as u on v.vid = u.vid where 1=1 %s`, filterSql)
	count, err := db.QuerySqlCount(service.BUDAODB, csql)
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

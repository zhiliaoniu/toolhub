package reportservice

import (
	"common"
	"context"
	"db"
	"fmt"
	"github.com/sumaig/glog"
	"github.com/tealeg/xlsx"
	"reflect"
	"service"
	"service/api"
	"strconv"
	"time"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

const (
	FILE_PATH = "/data/houshanjie/"
)

// app统计每日视频相关数-导出execl
func (s *Server) ReportVideoStatisticsExecl(ctx context.Context, req *api.QueryListRequest) (resp *api.ReportVideoExeclResponse, err error) {
	resp = &api.ReportVideoExeclResponse{
		Data: &api.ReportVideoExeclBack{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	var m = map[string]string{
		"sTime":        `create_time,>='%v'`,
		"eTime":        `create_time,<='%v'`,
		"vExposeNum":   `video_expose_num,>='%v'`,  //视频播放
		"vClictNum":    `video_clict_num,>='%v'`,   //视频点击
		"vViewNum":     `video_view_num,>='%v'`,    //视频观看
		"vFavorNum":    `video_favor_num,>='%v'`,   //视频点赞
		"commFavorNum": `comment_favor_num,>='%v'`, //评论点赞
		"commNum":      `comment_num,>='%v'`,       //评论发表
		"tFollowNum":   `topic_follow_num,>='%v'`,  //话题订阅数
	}

	filterSql, _, _, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	vSql := fmt.Sprintf(`select id, video_expose_num, video_clict_num, video_view_num, video_favor_num, comment_favor_num, comment_num, topic_follow_num, create_time from statis_biz_daily where 1=1 %s`, filterSql)
	rows, err := db.Query(service.BUDAODB, vSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	//现将数据写入excel里面
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	row := sheet.AddRow()
	row.SetHeightCM(1)
	cell := row.AddCell()
	cell.Value = "Id"
	cell = row.AddCell()
	cell.Value = "视频播放数"
	cell = row.AddCell()
	cell.Value = "视频点击数"
	cell = row.AddCell()
	cell.Value = "视频观看数"
	cell = row.AddCell()
	cell.Value = "视频点赞数"
	cell = row.AddCell()
	cell.Value = "评论点赞数"
	cell = row.AddCell()
	cell.Value = "评论发表数"
	cell = row.AddCell()
	cell.Value = "话题订阅数"
	cell = row.AddCell()
	cell.Value = "时间"

	for rows.Next() {
		temp := api.StatisBizDaily{}
		err := rows.Scan(&temp.Id, &temp.VideoExposeNum, &temp.VideoClictNum, &temp.VideoViewNum, &temp.VideoFavorNum, &temp.CommentFavorNum, &temp.CommentNum, &temp.TopicFollowNum, &temp.CreateTime)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			break
		}
		temp.CreateTime = common.GetTimeStr(temp.CreateTime)
		t := reflect.TypeOf(temp)
		v := reflect.ValueOf(temp)
		row := sheet.AddRow()
		row.SetHeightCM(1)
		for i := 0; i < t.NumField(); i++ {
			cell = row.AddCell()
			cell.Value = fmt.Sprintf("%v", v.Field(i).Interface())
		}
	}
	//./doc/downFile/report_video.xlsx
	fileName := time.Now().Format("2006-01-02_150405")
	filePath := FILE_PATH + "report_video_" + fileName + "xlsx"
	err1 := file.Save(filePath)
	if err1 != nil {
		panic(err1)
	}
	var fileUrl string
	if !common.GetConfig().Condition {
		fileUrl = "http://jxz-cms.yy.com" + filePath
	} else {
		fileUrl = "http://116.31.122.113" + filePath
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data.FileUrl = fileUrl
	return
}

// app统计每日用户相关数-导出execl
func (s *Server) ReportUserStatisticsExecl(ctx context.Context, req *api.QueryListRequest) (resp *api.ReportUserExeclResponse, err error) {
	resp = &api.ReportUserExeclResponse{
		Data: &api.ReportUserExeclBack{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	var m = map[string]string{
		"sTime":     `create_time,>='%v'`,
		"eTime":     `create_time,<='%v'`,
		"totalNum":  `total_num,>='%v'`,  //当前用户总数
		"activeNum": `active_num,>='%v'`, //前一天活跃用户数
		"newNum":    `new_num,>='%v'`,    //前一天新增用户数
	}

	filterSql, _, _, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	vSql := fmt.Sprintf(`select id, total_num, active_num, new_num, create_time from statis_new_user where 1=1 %s`, filterSql)
	rows, err := db.Query(service.BUDAODB, vSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	//现将数据写入excel里面
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	row := sheet.AddRow()
	row.SetHeightCM(1)
	cell := row.AddCell()
	cell.Value = "Id"
	cell = row.AddCell()
	cell.Value = "当前用户总数"
	cell = row.AddCell()
	cell.Value = "前一天活跃用户数"
	cell = row.AddCell()
	cell.Value = "前一天新增用户数"
	cell = row.AddCell()
	cell.Value = "时间"

	for rows.Next() {
		temp := api.StatisUser{}
		err = rows.Scan(&temp.Id, &temp.TotalNum, &temp.ActiveNum, &temp.NewNum, &temp.CreateTime)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		temp.CreateTime = common.GetTimeStr(temp.CreateTime)
		//根据映射遍历struct
		t := reflect.TypeOf(temp)
		v := reflect.ValueOf(temp)
		row := sheet.AddRow()
		row.SetHeightCM(1)
		for i := 0; i < t.NumField(); i++ {
			cell = row.AddCell()
			cell.Value = fmt.Sprintf("%v", v.Field(i).Interface())
		}
	}
	//./doc/downFile/report_user.xlsx
	fileName := time.Now().Format("2006-01-02_150405")
	filePath := FILE_PATH + "report_user_" + fileName + ".xlsx"
	err1 := file.Save(filePath)
	if err1 != nil {
		panic(err1)
	}
	var fileUrl string
	if !common.GetConfig().Condition {
		fileUrl = "http://jxz-cms.yy.com" + filePath
	} else {
		fileUrl = "http://116.31.122.113" + filePath
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data.FileUrl = fileUrl

	return
}

//每日视频数据统计
func (s *Server) ReportVideoStatistics(ctx context.Context, req *api.QueryListRequest) (resp *api.ReportVideoStatisticsResponse, err error) {
	resp = &api.ReportVideoStatisticsResponse{
		Data: &api.ReportVideoStatisticsBack{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	var m = map[string]string{
		"sTime":        `create_time,>='%v'`,
		"eTime":        `create_time,<='%v'`,
		"vExposeNum":   `video_expose_num,>='%v'`,  //视频播放
		"vClictNum":    `video_clict_num,>='%v'`,   //视频点击
		"vViewNum":     `video_view_num,>='%v'`,    //视频观看
		"vFavorNum":    `video_favor_num,>='%v'`,   //视频点赞
		"commFavorNum": `comment_favor_num,>='%v'`, //评论点赞
		"commNum":      `comment_num,>='%v'`,       //评论发表
		"tFollowNum":   `topic_follow_num,>='%v'`,  //话题订阅数

	}

	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	vSql := fmt.Sprintf(`select id, video_expose_num, video_clict_num, video_view_num, video_favor_num, comment_favor_num, comment_num, topic_follow_num, create_time from statis_biz_daily where 1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, vSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	data := make([]*api.StatisBizDaily, 0, 10)
	for rows.Next() {
		temp := api.StatisBizDaily{}
		err = rows.Scan(&temp.Id, &temp.VideoExposeNum, &temp.VideoClictNum, &temp.VideoViewNum, &temp.VideoFavorNum, &temp.CommentFavorNum, &temp.CommentNum, &temp.TopicFollowNum, &temp.CreateTime)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		temp.CreateTime = common.GetTimeStr(temp.CreateTime)

		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from statis_biz_daily where 1=1 %s`, filterSql))
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

//统计每日用户数
func (s *Server) ReportUserStatistics(ctx context.Context, req *api.QueryListRequest) (resp *api.ReportUserStatisticsResponse, err error) {
	resp = &api.ReportUserStatisticsResponse{
		Data: &api.ReportUserStatisticsBack{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	var m = map[string]string{
		"sTime":     `create_time,>='%v'`,
		"eTime":     `create_time,<='%v'`,
		"totalNum":  `total_num,>='%v'`,  //当前用户总数
		"activeNum": `active_num,>='%v'`, //前一天活跃用户数
		"newNum":    `new_num,>='%v'`,    //前一天新增用户数
	}

	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	vSql := fmt.Sprintf(`select id, total_num, active_num, new_num, create_time from statis_new_user where 1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, vSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	data := make([]*api.StatisUser, 0, 10)
	for rows.Next() {
		temp := api.StatisUser{}
		err = rows.Scan(&temp.Id, &temp.TotalNum, &temp.ActiveNum, &temp.NewNum, &temp.CreateTime)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		temp.CreateTime = common.GetTimeStr(temp.CreateTime)

		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from statis_new_user where 1=1 %s`, filterSql))
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

//内容数据  话题数据
func (s *Server) ReportContentTopicData(ctx context.Context, req *api.QueryListRequest) (resp *api.ReportContentTopicDataResponse, err error) {
	resp = &api.ReportContentTopicDataResponse{
		Data: &api.ReportContentTopicDataBack{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	param := make(map[string]string)
	_, sortSql, pageSql, err := service.GetSqlParam(*req, param)
	tSql := fmt.Sprintf(`select topic_id, name, user_num, video_num from topic %s %s`, sortSql, pageSql)

	//话题总数
	tpiceCount, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select topic_id from topic`))
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}

	//话题启用数
	topicOpenCount, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select topic_id from topic  where disable = 0`))
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}

	//话题禁用数
	topicCloseCount, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select topic_id from topic  where disable = 1`))
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}

	//每个话题下的视频数  视频问题数
	rows, err := db.Query(service.BUDAODB, tSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()

	type TopicCount struct {
		TopicID  string `json:"topicID"`
		Name     string `json:"name"`
		UserNum  string `json:"userCount"`
		VideoNum string `json:"videoCount"`
		QueCount string `json:"queCount"`
	}

	data := make([]*api.TopicOther, 0, 10)
	for rows.Next() {
		m := api.TopicOther{}
		temp := TopicCount{}
		err = rows.Scan(&temp.TopicID, &temp.Name, &temp.UserNum, &temp.VideoNum)
		if err != nil {
			err = service.MysqlError
			return
		}
		//获取话题下的视频
		tID, _ := strconv.ParseUint(temp.TopicID, 10, 64)
		tableName, err := db.GetTableName(service.TOPIC_VIDEO_, tID)
		glog.Error(err)

		queCount, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select t.id from %v as t inner join question as q on t.vid = q.vid where t.topic_id = '%v'`, tableName, temp.TopicID))
		if err != nil {
			glog.Error(err)
			return nil, service.MysqlError
		}
		m.TopicName = temp.Name
		m.VideoCount = temp.VideoNum
		m.QueCount = strconv.FormatUint(queCount, 10)
		m.UserCount = temp.UserNum

		//temp.QueCount =  strconv.FormatUint(queCount,10)
		data = append(data, &m)
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = tpiceCount
	resp.Data.TOCount = topicOpenCount
	resp.Data.TCCount = topicCloseCount
	return
}

//内容数据  发布数据
func (s *Server) ReportContentPostData(ctx context.Context, req *api.QueryListRequest) (resp *api.ReportContentPostDataResponse, err error) {
	resp = &api.ReportContentPostDataResponse{
		Data: &api.ReportContentPostDataBack{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	var m = map[string]string{
		"sTime":   `v.create_time,>='%v'`,
		"eTime":   `v.create_time,<='%v'`,
		"source":  `v.source,='%v'`,
		"topicID": `v.topic_id,like '%%%v%%'`,
	}
	filterSql, _, _, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	//发布的视频数
	querySql := fmt.Sprintf(`select v.id from video_0 as v where 1=1 %s`, filterSql)
	countTotal, err := db.QuerySqlCount(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}

	//发布视频含有问题的数
	qesSql := fmt.Sprintf(`select v.id from video_0 as v inner join question as q on v.vid = q.vid where 1=1 %s`, filterSql)
	countQ, err := db.QuerySqlCount(service.BUDAODB, qesSql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data.PostVideoCount = countTotal
	resp.Data.PostVideoQueCount = countQ
	return
}

//pushVV列表
func (s *Server) ReportPushVVList(ctx context.Context, req *api.ReportPushVVListRequest) (resp *api.ReportPushVVListResponse, err error) {
	resp = &api.ReportPushVVListResponse{}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	querySql := fmt.Sprintf(`select a.*,ifnull(v.title,''),ifnull(v.duration,'') from (select vid,sum(vv) as sumVV from budao_test.push_video_vv_stat where %v='%v' and photo_type='%v' group by vid) a left join video_0 v on a.vid=v.vid`, `date_format(update_time,'%Y-%m-%d')`, req.Date, req.PhotoType)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	result := make([]*api.ReportPushVVListBack, 0, 10)
	for rows.Next() {
		temp := api.ReportPushVVListBack{}
		err = rows.Scan(&temp.Vid, &temp.SumVV, &temp.Title, &temp.Duration)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		result = append(result, &temp)
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data = result

	return
}

//pushVV
func (s *Server) ReportPushVV(ctx context.Context, req *api.ReportPushVVRequest) (resp *api.ReportPushVVResponse, err error) {
	resp = &api.ReportPushVVResponse{}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	if req.STime == "" {
		today := time.Now()
		req.STime = today.Add(-7*24*time.Hour).Format(`2006-01-02`) + ` 00:00:00`
		req.ETime = today.Add(-24*time.Hour).Format(`2006-01-02`) + ` 23:59:59`
	}
	var filterSql string
	filterSql = fmt.Sprintf(` and update_time>='%v' and update_time<='%v'`, req.STime, req.ETime)
	if req.PhotoType != "" {
		filterSql += fmt.Sprintf(` and photo_type='%v'`, req.PhotoType)
	}
	querySql := fmt.Sprintf(`select %v date,sum(vv) as sumVV from budao_test.push_video_vv_stat where 1=1 %v group by %v)`, `date_format(update_time,'%Y-%m-%d')`, filterSql, `date_format(update_time,'%Y-%m-%d'`)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	result := make([]*api.ReportPushVVBack, 0, 10)
	for rows.Next() {
		temp := api.ReportPushVVBack{}
		err = rows.Scan(&temp.Date, &temp.SumVV)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		result = append(result, &temp)
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data = result
	return
}

//运营今日数据
func (s *Server) ReportOpToday(ctx context.Context, req *api.NullReq) (resp *api.ReportOpTodayResponse, err error) {
	resp = &api.ReportOpTodayResponse{
		Data: &api.ReportOpBack{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	querySql := fmt.Sprintf(`select count(a.vid) as allPush, count(q.vid) as hasQuestion from (select vid from video_0 where TO_DAYS(create_time)=TO_DAYS(NOW()) and type!=100 and type!=99) a left join
(select vid from question where TO_DAYS(create_time)=TO_DAYS(NOW()) group by vid) q on a.vid=q.vid `)

	rows, err := db.QueryRow(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	result := api.ReportOpBack{}
	err = rows.Scan(&result.AllPush, &result.HasQuestion)
	if err != nil {
		err = service.MysqlError
		return
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data = &result
	return
}

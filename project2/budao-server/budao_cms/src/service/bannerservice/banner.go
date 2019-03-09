package bannerservice

import (
	"common"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"github.com/sumaig/glog"
	"service"
	"service/api"
	"strings"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

//删除banner
func (s *Server) BannerDel(ctx context.Context, req *api.BannerDelRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	querySql := fmt.Sprintf(`update banner set status = 1 where id = '%v'`, req.Id)
	_, err = db.Exec(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"

	return
}

//修改banner
func (s *Server) BannerModify(ctx context.Context, req *api.Banner) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	var m = map[string]string{
		"id":          `id,='%v'`,
		"picUrl":      `pic_url,='%v'`,
		"clickUrl":    `link,='%v'`,
		"position":    `position,='%v'`,
		"description": `description,='%v'`,
		"status":      `status,='%v'`,
		"fromTime":    `effective_time,='%v'`,
		"toTime":      `ineffective_time,='%v'`,
	}
	b, _ := json.Marshal(req)
	param, updateSql, err := service.GetUpdateSql(b, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	querySql := fmt.Sprintf(`update banner set %v where id = '%v'`, updateSql, param["id"])
	_, err = db.Exec(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"

	return
}

//获取banner信息
func (s *Server) BannerInfo(ctx context.Context, req *api.BannerInfoRequest) (resp *api.BannerInfoResponse, err error) {
	resp = &api.BannerInfoResponse{
		Data: &api.Banner{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	querySql := fmt.Sprintf(`select id, pic_url, link, position, status, description, effective_time, ineffective_time from banner where id = '%v'`, req.Id)
	row, err := db.QueryRow(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}

	temp := api.Banner{}
	err = row.Scan(&temp.Id, &temp.PicUrl, &temp.ClickUrl, &temp.Position, &temp.Status, &temp.Description, &temp.FromTime, &temp.ToTime)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	temp.FromTime = common.GetTimeStr(temp.FromTime)
	temp.ToTime = common.GetTimeStr(temp.ToTime)

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data = &temp
	return
}

//banner列表
func (s *Server) BannerList(ctx context.Context, req *api.QueryListRequest) (resp *api.BannerListResponse, err error) {
	resp = &api.BannerListResponse{
		Data: &api.BannerList{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	var m = map[string]string{
		"position":    `position,='%v'`,
		"description": `description,like '%%v%%'`,
		"status":      `status,='%v'`,
		"crTime":      `created_at,<='%v'`,
		"fromTime":    `effective_time,>='%v'`,   //起始时间
		"toTime":      `ineffective_time,<='%v'`, //终止时间
	}

	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	querySql := fmt.Sprintf(`select id, pic_url, link, position, status, description, created_at, updated_at, effective_time, ineffective_time from banner where 1=1 %v %v %v`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	data := make([]*api.Banner, 0, 10)
	for rows.Next() {
		temp := api.Banner{}
		err := rows.Scan(&temp.Id, &temp.PicUrl, &temp.ClickUrl, &temp.Position, &temp.Status, &temp.Description, &temp.CreatedAt, &temp.UpdatedAt, &temp.FromTime, &temp.ToTime)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			break
		}

		temp.CreatedAt = common.GetTimeStr(temp.CreatedAt)
		temp.UpdatedAt = common.GetTimeStr(temp.UpdatedAt)
		temp.FromTime = common.GetTimeStr(temp.FromTime)
		temp.ToTime = common.GetTimeStr(temp.ToTime)

		data = append(data, &temp)
	}

	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from banner where 1=1 %s`, filterSql))
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = count
	return
}

//添加banner
func (s *Server) BannerAdd(ctx context.Context, req *api.Banner) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	if !strings.Contains(req.ClickUrl, "://") {
		err = fmt.Errorf(`link url is error`)
		return
	}
	querySql := fmt.Sprintf(`insert into banner (pic_url, link, position, description, effective_time, ineffective_time) values('%v', '%v', '%v', '%v', '%v', '%v')`, req.PicUrl, req.ClickUrl, req.Position, req.Description, req.FromTime, req.ToTime)

	_, err = db.Exec(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	return
}

package pushservice

import (
	"common"
	"context"
	"db"
	"encoding/json"
	"fmt"
	"github.com/sumaig/glog"
	"regexp"
	"service"
	"service/api"
	"strings"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

//删除push
func (s *Server) PushDel(ctx context.Context, req *api.PushDelRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	querySql := fmt.Sprintf(`update push set status = 1 where id = '%v'`, req.PushId)
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

//修改Push
func (s *Server) PushModify(ctx context.Context, req *api.Push) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	var m = map[string]string{
		"id":            `id,='%v'`,
		"opUid":         `op_uid,='%v'`,
		"pushChannel":   `push_channel,='%v'`,
		"pushType":      `push_type,='%v'`,
		"pushTitle":     `push_title,='%v'`,
		"pushContent":   `push_content,='%v'`,
		"pushDescribe":  `push_describe,='%v'`,
		"pushTime":      `push_time,='%v'`,
		"pushEndTime":   `push_end_time,='%v'`,
		"pushVideoType": `push_video_type,='%v'`,
		"pushUrlType":   `push_url_type,='%v'`,
		"device":        `device,='%v'`,
		"status":        `status,='%v'`,
	}
	b, _ := json.Marshal(req)
	param, updateSql, err := service.GetUpdateSql(b, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	//根据推送点击方式拼接不同的连接地址
	pushUrl := dealPushUrl(param["pushUrlType"], param["pushVideoType"], req.CommonStr)
	updateSql += " ,push_url = '" + pushUrl + "'"
	glog.Info(updateSql)
	querySql := fmt.Sprintf(`update push set %v where id = '%v'`, updateSql, param["id"])
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

//获取push信息
func (s *Server) PushInfo(ctx context.Context, req *api.PushInfoRequest) (resp *api.PushInfoResponse, err error) {
	resp = &api.PushInfoResponse{
		Data: &api.Push{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	querySql := fmt.Sprintf(`select id, op_uid, push_obj, push_channel, push_type, push_status, push_title, push_content, push_describe, push_time, push_end_time, push_video_type, push_url_type, push_url, device, status from push where id = '%v'`, req.PushId)
	row, err := db.QueryRow(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}

	temp := api.Push{}
	err = row.Scan(&temp.Id, &temp.OpUid, &temp.PushObj, &temp.PushChannel, &temp.PushType, &temp.PushStatus, &temp.PushTitle, &temp.PushContent, &temp.PushDescribe, &temp.PushTime, &temp.PushEndTime, &temp.PushVideoType, &temp.PushUrlType, &temp.PushUrl, &temp.Device, &temp.Status)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	//获取地址上的id
	reg := regexp.MustCompile(`[\d]+`)
	temp.CommonStr = reg.FindString(temp.PushUrl)
	temp.PushTime = common.GetTimeStr(temp.PushTime)
	temp.PushEndTime = common.GetTimeStr(temp.PushEndTime)

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data = &temp
	return
}

//push列表
func (s *Server) PushList(ctx context.Context, req *api.QueryListRequest) (resp *api.PushListResponse, err error) {
	resp = &api.PushListResponse{
		Data: &api.PushList{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	var m = map[string]string{
		"pushDescribe":  `push_describe,like '%%%v%%'`,
		"pushChannel":   `push_channel,='%v'`,
		"pushTitle":     `push_title,like '%%%v%%'`,
		"fromTime":      `push_time,>='%v'`,
		"toTime":        `push_time,<='%v'`,
		"pushType":      `push_type,='%v'`,
		"pushStatus":    `push_status,='%v'`,
		"pushVideoType": `push_video_type,='%v'`,
		"pushUrlType":   `push_url_type,='%v'`,
		"status":        `status,='%v'`,
	}

	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	querySql := fmt.Sprintf(`select id, op_uid, push_obj, push_channel, push_type, push_status, push_title, push_content, push_describe, push_time, push_end_time, push_video_type, push_url_type, push_url, device, status, created_at, updated_at from push where 1=1 %v %v %v`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	data := make([]*api.Push, 0, 10)
	for rows.Next() {
		temp := api.Push{}
		err := rows.Scan(&temp.Id, &temp.OpUid, &temp.PushObj, &temp.PushChannel, &temp.PushType, &temp.PushStatus, &temp.PushTitle, &temp.PushContent, &temp.PushDescribe, &temp.PushTime, &temp.PushEndTime, &temp.PushVideoType, &temp.PushUrlType, &temp.PushUrl, &temp.Device, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			break
		}
		temp.PushTime = common.GetTimeStr(temp.PushTime)
		temp.PushEndTime = common.GetTimeStr(temp.PushEndTime)
		temp.CreatedAt = common.GetTimeStr(temp.CreatedAt)
		temp.UpdatedAt = common.GetTimeStr(temp.UpdatedAt)

		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from push where 1=1 %s`, filterSql))

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = count
	return
}

//添加push
func (s *Server) PushAdd(ctx context.Context, req *api.Push) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	//根据推送点击方式拼接不同的连接地址
	req.PushUrl = dealPushUrl(req.PushUrlType, req.PushVideoType, req.CommonStr)

	querySql := fmt.Sprintf(`insert into push (push_obj, push_channel, push_type, push_title, push_content, push_describe, push_time, push_end_time, push_video_type, push_url_type, push_url, device) 
values('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')`, req.PushObj, req.PushChannel, req.PushType, req.PushTitle, req.PushContent, req.PushDescribe, req.PushTime, req.PushEndTime, req.PushVideoType, req.PushUrlType, req.PushUrl, req.Device)

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

//对pushUrl进行处理
func dealPushUrl(pushUrlType, pushVideoType, commonStr string) string {
	var pushUrl string
	//根据推送点击方式拼接不同的连接地址
	switch pushUrlType {
	case "0":
		switch pushVideoType {
		case "0":
			pushUrl = "zvideo://topicdetail?topicId=" + strings.Trim(commonStr, " ") + "&from=push"
		case "1":
			pushUrl = "zvideo://videodetail?videoId=" + strings.Trim(commonStr, " ") + "&from=push"
		default:
			pushUrl = "zvideo://"
		}
	case "1":
		pushUrl = "zvideo://"
	case "2":
		pushUrl = commonStr
	default:
		pushUrl = "zvideo://"
	}

	return pushUrl
}

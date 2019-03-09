package commonservice

import (
	"context"
	"db"
	"fmt"
	"github.com/sumaig/glog"
	"service"
	"service/api"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) CategoryList(ctx context.Context, req *api.NullReq) (resp *api.CategoryListResponse, err error) {
	resp = &api.CategoryListResponse{}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	querySql := fmt.Sprintf(`select * from question_category`)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	result := make([]*api.QuestionCategory, 0, 10)
	for rows.Next() {
		temp := api.QuestionCategory{}
		err = rows.Scan(&temp.Id, &temp.CategoryName)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}
		result = append(result, &temp)
	}
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data = result
	return
}

func (s *Server) TypeList(ctx context.Context, req *api.NullReq) (resp *api.TypeListResponse, err error) {
	resp = &api.TypeListResponse{}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	querySql := fmt.Sprintf(`select * from question_type`)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	result := make([]*api.QuestionType, 0, 10)
	for rows.Next() {
		temp := api.QuestionType{}
		err = rows.Scan(&temp.Id, &temp.TypeName)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		result = append(result, &temp)
	}

	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data = result
	return
}

func (s *Server) AddCategory(ctx context.Context, req *api.AddCategoryReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	querySql := fmt.Sprintf(`insert into question_category (category_name) values ('%s')`, req.CategoryName)
	_, err = db.Exec(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	return
}

func (s *Server) AddType(ctx context.Context, req *api.AddTypeReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	querySql := fmt.Sprintf(`insert into question_type (type_name) values ('%s')`, req.TypeName)
	_, err = db.Exec(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	return
}

func (s *Server) InternalUserList(ctx context.Context, req *api.QueryListRequest) (resp *api.InternalUserListResponse, err error) {
	resp = &api.InternalUserListResponse{
		Data: &api.UserList{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	allUser := service.IuCache.Get()
	count := int32(len(allUser))
	start := (req.Num - 1) * req.Size
	end := start + req.Size
	if end > count {
		end = count
	}
	if start >= count {
		start = 0
		end = 0
	}

	data := allUser[start:end]
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Count = uint64(len(allUser))
	for _, v := range data {
		u := &api.User{}
		u.Uid = v.Uid
		u.Photo = v.Photo
		u.Name = v.Name
		resp.Data.Data = append(resp.Data.Data, u)
	}
	return
}

package keyword

import (
	"common"
	"context"
	"db"
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

func (s *Server) QueryVideoTitle(ctx context.Context, req *api.QueryListRequest) (resp *api.VideoListResponse, err error) {
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
		"keyworld": `title,like '%%%v%%'`,
	}
	filterSql, _, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	querySql := fmt.Sprintf(`select vid,title,coverurl,videourl from video_0 where 1=1 and state=2 %v %v`, filterSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()
	data := make([]*api.Video, 0, 10)

	for rows.Next() {
		temp := api.Video{}

		err = rows.Scan(&temp.Vid, &temp.Title, &temp.CoverUrl, &temp.VideoUrl)

		if err != nil {
			err = service.MysqlError
			return
		}
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from video_0 where 1=1 %v`, filterSql))
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = count

	return
}

func (s *Server) ReplaceVideoTitle(ctx context.Context, req *api.OptionVideoTitleReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	updateSql := fmt.Sprintf(`update video_0 set title= replace(title,'%v','%v')`, req.Keyword, req.NewWord)
	_, err = db.Exec(service.BUDAODB, updateSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		err = nil
	}
	return
}

func (s *Server) QueryCommentContent(ctx context.Context, req *api.QueryCommentContentReq) (resp *api.CommentListResponse, err error) {
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

	allDate := make([]*api.Comment, 0, 50)
	maxNum := common.GetConfig().DB.MySQL[service.BUDAODB].TableDesc[service.COMMENT_]

	for i := uint64(0); i < maxNum; i++ {
		tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
		querySql := fmt.Sprintf(`select cid,vid,from_uid,from_name,from_photo,content,favor_num from %v where state=2 and content like '%%%v%%'`, tableName, req.Keyword)
		rows, err := db.Query(service.BUDAODB, querySql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}
		for rows.Next() {
			temp := new(api.Comment)
			err = rows.Scan(&temp.Cid, &temp.Vid, &temp.Uid, &temp.UName, &temp.UPhoto, &temp.Content, &temp.FavorNum)
			if err != nil {
				glog.Error(err)
				err = service.MysqlError
				return resp, err
			}

			allDate = append(allDate, temp)

		}
	}

	count := int32(len(allDate))
	start := (req.Num - 1) * req.Size
	end := start + req.Size
	if end > count {
		end = count
	}
	if start >= count {
		start = 0
		end = req.Size
	}
	data := allDate[start:end]
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = uint64(len(allDate))
	return
}

func (s *Server) OptionCommentContent(ctx context.Context, req *api.OptionCommentContentReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	var execSql string
	if req.Option == "DEL" {
		execSql = fmt.Sprintf(`update tableName set state=4 where content like '%%%v%%'`, req.Keyword)
	} else {
		execSql = fmt.Sprintf(`update tableName set content= replace(content,'%v','%v')`, req.Keyword, req.NewWord)
	}
	maxNum := common.GetConfig().DB.MySQL[service.BUDAODB].TableDesc[service.COMMENT_]

	for i := uint64(0); i < maxNum; i++ {
		tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
		execSql = strings.Replace(execSql, `tableName`, tableName, -1)
		_, err = db.Exec(service.BUDAODB, execSql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
	}
	return
}

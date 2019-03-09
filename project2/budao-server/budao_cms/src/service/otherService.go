package service

import (
	"common"
	"context"
	"db"
	"fmt"
	"github.com/sumaig/glog"
	"github.com/tealeg/xlsx"
	"github.com/twitchtv/twirp/ctxsetters"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type OtherServie struct {
	service
}

func NewOtherServie() *OtherServie {
	return &OtherServie{
		service{},
	}
}

func (s *OtherServie) ServiceDescriptor() ([]byte, int) {
	return []byte{}, 0
}

func (s *OtherServie) ProtocGenTwirpVersion() string {
	return "v5.3.0"
}
func (s *OtherServie) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "budao")
	ctx = ctxsetters.WithServiceName(ctx, "CommonServie")
	ctx = ctxsetters.WithResponseWriter(ctx, resp)
	url, buff, err := s.verfiy(req)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	resp.Header().Add("Access-Control-Allow-Origin", "*")
	switch url {
	//上传图片
	case "/other/uploadPic":
		ctx = ctxsetters.WithMethodName(ctx, "uploadPic")
		s.serviceUploadPic(ctx, resp, req)
	//视频相关导出execl
	case "/other/video/xlsx":
		ctx = ctxsetters.WithMethodName(ctx, "videoXlsx")
		s.serviceGetVideoExcelFile(ctx, resp, buff, req)
	//用户相关导出execl
	case "/other/user/xlsx":
		ctx = ctxsetters.WithMethodName(ctx, "userXlsx")
		s.serviceGetUserExcelFile(ctx, resp, buff, req)
	default:
		resp.WriteHeader(http.StatusNotFound)
	}

}

//用户相关导出execl
func (s *OtherServie) serviceGetUserExcelFile(ctx context.Context, resp http.ResponseWriter, buff []byte, req *http.Request) {
	var m = map[string]string{
		"sTime":     `create_time,>='%v'`,
		"eTime":     `create_time,<='%v'`,
		"totalNum":  `total_num,>='%v'`,  //当前用户总数
		"activeNum": `active_num,>='%v'`, //前一天活跃用户数
		"newNum":    `new_num,>='%v'`,    //前一天新增用户数
	}

	filterSql, _, _, err := getSqlParam(buff, m)
	if err != nil {
		glog.Error(err)
		s.writeError(ctx, resp, ParamError)
		return
	}

	vSql := fmt.Sprintf(`select id, total_num, active_num, new_num, create_time from statis_new_user where 1=1 %s`, filterSql)
	rows, err := db.Query(BUDAODB, vSql)
	if err != nil {
		glog.Error(err)
		s.writeError(ctx, resp, MysqlError)
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
		temp := StatisNewUser{}
		err := rows.Scan(&temp.Id, &temp.TotalNum, &temp.ActiveNum, &temp.NewNum, &temp.CreateTime)
		if err != nil {
			glog.Error(err)
			s.writeError(ctx, resp, MysqlError)
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

	err1 := file.Save("./doc/downFile/report_user.xlsx")
	if err1 != nil {
		panic(err1)
	}

	//下载文件
	common.DownloadFile("./doc/downFile/report_user.xlsx", resp)

	writeJsonResp(resp, nil)
}

//视频相关导出execl
func (s *OtherServie) serviceGetVideoExcelFile(ctx context.Context, resp http.ResponseWriter, buff []byte, req *http.Request) {
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

	filterSql, _, _, err := getSqlParam(buff, m)
	if err != nil {
		glog.Error(err)
		s.writeError(ctx, resp, ParamError)
		return
	}

	vSql := fmt.Sprintf(`select id, video_expose_num, video_clict_num, video_view_num, video_favor_num, comment_favor_num, comment_num, topic_follow_num, create_time from statis_biz_daily where 1=1 %s`, filterSql)
	rows, err := db.Query(BUDAODB, vSql)
	if err != nil {
		glog.Error(err)
		s.writeError(ctx, resp, MysqlError)
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
		temp := StatisBizDaily{}
		err := rows.Scan(&temp.Id, &temp.VideoExposeNum, &temp.VideoClictNum, &temp.VideoViewNum, &temp.VideoFavorNum, &temp.CommentFavorNum, &temp.CommentNum, &temp.TopicFollowNum, &temp.CreateTime)
		if err != nil {
			glog.Error(err)
			s.writeError(ctx, resp, MysqlError)
			return
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

	err1 := file.Save("./doc/downFile/report_video.xlsx")
	if err1 != nil {
		panic(err1)
	}

	//下载文件
	common.DownloadFile("./doc/downFile/report_video.xlsx", resp)

	writeJsonResp(resp, nil)
}

type UploadPicParams struct {
	PicUrl string `json:"picUrl"`
}

//上传图片
func (s *OtherServie) serviceUploadPic(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var (
		picurl string
	)

	fileType := req.FormValue("fileType")
	if fileType == "" {
		uploadFile, hander, err := req.FormFile("picName")
		glog.Info(err)
		if err != nil {
			s.writeError(ctx, resp, ParamError)
			return
		}
		picName := hander.Filename

		buff, err := ioutil.ReadAll(uploadFile)
		if err != nil {
			s.writeError(ctx, resp, err)
			return
		}
		picurl, _ = uploadPic("", picName, buff)
	} else {
		filepath := req.FormValue("picName")
		if !strings.Contains(filepath, "http") {
			s.writeError(ctx, resp, fmt.Errorf("url格式不正确"))
			return
		}

		picurl, _ = uploadPic(filepath, "", nil)
	}

	result := UploadPicParams{
		PicUrl: picurl,
	}

	writeJsonResp(resp, result)
}

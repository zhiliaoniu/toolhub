package blacklistservice

import (
	"common"
	"context"
	"db"
	"fmt"
	"github.com/sumaig/glog"
	"service"
	"service/api"
	"strings"
	"sync"
	"time"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) AddBlackList(ctx context.Context, req *api.AddBlackListReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	status := 0
	if req.Duration == 0 {
		status = 1
	}
	if req.Duration == -1 {
		req.ETime = `NULL`
	} else {
		req.ETime = fmt.Sprintf(`'%v'`, req.ETime)
	}
	insertSql := fmt.Sprintf(`insert into user_blacklist (uid,blacklist_type,status,blacklist_duration,reason,start_time,end_time) values `)
	for _, uid := range req.Uids {
		insertSql += fmt.Sprintf(` (%v,%v,%v,%v,'%v','%v',%v),`, uid, req.BlacklistType, status, req.Duration, req.Reason, req.STime, req.ETime)
	}

	insertSql = insertSql[:len(insertSql)-1]

	_, err = db.Exec(service.BUDAODB, insertSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	return
}

func (s *Server) BlacklistUserList(ctx context.Context, req *api.QueryListRequest) (resp *api.BlacklistUserListResponse, err error) {
	resp = &api.BlacklistUserListResponse{
		Data: &api.BlackListUserList{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	var m = map[string]string{
		"uid":           `uid,like '%%%v%%'`,
		"blacklistType": `blacklist_type,=%v`,
		"status":        `status,=%v`,
		"sTime":         `start_time,>='%v'`,
		"eTime":         `end_time,>='%v'`,
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	querySql := fmt.Sprintf(`select uid,blacklist_type,status,blacklist_duration,reason,ifnull(start_time,''),ifnull(end_time,'') from user_blacklist where 1=1 %v %v %v`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	type BlackList struct {
		Uid           string `json:"uid"`
		BlacklistType string `json:"blacklistType"`
		Status        string `json:"status"`
		Duration      int    `json:"duration"`
		Reason        string `json:"reason"`
		STime         string `json:"sTime"`
		ETime         string `json:"eTime"`
	}

	data := make([]*api.BlackList, 0, 10)
	for rows.Next() {
		temp := api.BlackList{}
		err = rows.Scan(&temp.Uid, &temp.BlacklistType, &temp.Status, &temp.Duration, &temp.Reason, &temp.STime, &temp.ETime)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		data = append(data, &temp)
	}

	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from user_blacklist where 1=1 %s`, filterSql))
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = count

	return
}

func (s *Server) UnlockUser(ctx context.Context, req *api.UnlockUserReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	updateSql := fmt.Sprintf(`update user_blacklist set status=1 where uid in (%v);`, strings.Join(req.Uids, `,`))
	_, err = db.Exec(service.BUDAODB, updateSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	return
}

func (s *Server) UserCommonList(ctx context.Context, req *api.CommentListReq) (resp *api.CommentListResponse, err error) {
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
	if data, ok := userCommentCache.get(req.Uid); ok {
		allDate = data
		goto returnData
	}

	for i := uint64(0); i < maxNum; i++ {
		tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
		querySql := fmt.Sprintf(`select cid,vid,from_uid,from_name,from_photo,to_comment_id,parentcomid,content,favor_num,weight,reply_num,state,create_time,update_time from %v where from_uid=%v`, tableName, req.Uid)
		rows, err := db.Query(service.BUDAODB, querySql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}
		for rows.Next() {
			temp := new(api.Comment)
			err = rows.Scan(&temp.Cid, &temp.Vid, &temp.Uid, &temp.UName, &temp.UPhoto, &temp.PId, &temp.MasterId, &temp.Content, &temp.FavorNum, &temp.Weight, &temp.ReplyNum, &temp.State, &temp.CTime, &temp.UTime)
			if err != nil {
				glog.Error(err)
				err = service.MysqlError
				return resp, err
			}

			allDate = append(allDate, temp)
			//allDate[temp.Cid] = temp
		}
	}
	userCommentCache.put(req.Uid, allDate)
returnData:
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

func (s *Server) RemoveUserCommon(ctx context.Context, req *api.RemoveUserCommonReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	userCommentCache.del(req.Uid)
	whereSql := fmt.Sprintf(`where  from_uid=%v `, req.Uid)
	if len(req.In) > 0 {
		whereSql += fmt.Sprintf(`and cid in (%v) `, strings.Join(req.In, `,`))
	}
	if len(req.Out) > 0 {
		whereSql += fmt.Sprintf(`and cid not in  (%v) `, strings.Join(req.In, `,`))
	}
	maxNum := common.GetConfig().DB.MySQL[service.BUDAODB].TableDesc[service.COMMENT_]
	for i := uint64(0); i < maxNum; i++ {
		tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
		updateSql := fmt.Sprintf(`update %v set state=4 %v `, tableName, whereSql)
		_, err = db.Exec(service.BUDAODB, updateSql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
	}
	return
}

//视频评论缓存
var userCommentCache *commentCache = func() *commentCache {

	cache := &commentCache{
		data: make(map[string]*commentCacheDate),
	}
	go func() {
		for {
			now := time.Now()
			for vid, v := range cache.data {
				if now.Sub(v.inTime) > 5*time.Minute {
					delete(cache.data, vid)
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	return cache
}()

type commentCache struct {
	sync.Mutex
	data map[string]*commentCacheDate
}
type commentCacheDate struct {
	inTime time.Time
	data   []*api.Comment
}

func (c commentCache) get(uid string) ([]*api.Comment, bool) {
	c.Lock()
	defer c.Unlock()
	data, ok := c.data[uid]
	if !ok {
		return nil, ok
	}
	now := time.Now()
	if now.Sub(data.inTime) > 5*time.Minute {
		delete(c.data, uid)
		return nil, false
	}
	return data.data, ok
}
func (c commentCache) del(uid string) {
	c.Lock()
	defer c.Unlock()
	delete(c.data, uid)
}

func (c commentCache) put(uid string, data []*api.Comment) {
	c.Lock()
	defer c.Unlock()
	cacheDate := commentCacheDate{
		inTime: time.Now(),
		data:   data,
	}
	c.data[uid] = &cacheDate
}

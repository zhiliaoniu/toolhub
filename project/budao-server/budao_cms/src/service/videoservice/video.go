package videoservice

import (
	"bytes"
	"common"
	"context"
	"db"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sumaig/glog"
	"service"
	"service/api"
	"strconv"
	"strings"
	"sync"
	"time"
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
		"vid":         `v.vid,='%v'`,
		"postType":    `v.post_type,='%v'`,
		"opStatus":    `v.op_state,='%v'`,
		"status":      `v.state,=%v`,
		"durationS":   `v.duration,>=%v`,
		"durationE":   `v.duration,<=%v`,
		"title":       `v.title,like '%%%v%%'`,
		"cTime":       `v.create_time,>='%v'`,
		"titleLenS":   "char_length(v.title),>=%v",
		"titleLenE":   "char_length(v.title),<=%v",
		"hasQuestion": `q.qCount,%v`,
		"hasTopic":    `t.tCount,%v`,
		"topicId":     `t.allTid,like '%%%v%%'`,
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	querySql := fmt.Sprintf(`select post_type,op_state,char_length(title) as title_len,v.vid,source_vid,title,coverurl,videourl,state,duration,ifnull(q.qCount,0),ifnull(t.tCount,0),ifnull(t.allTopic,''),ifnull(t.allTid,'')
 from video_0 v  left join 
(select vid,count(vid) as qCount  from question group by vid) q on v.vid=q.vid
left join (select vid,count(vid) as tCount,group_concat(t.name) allTopic,group_concat(t.topic_id) allTid  from topic_video_0 tv left join topic t on tv.topic_id=t.topic_id group by vid) t on v.vid=t.vid
where 1=1 %s %s %s`, filterSql, sortSql, pageSql)
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
		var (
			qcount    int
			tcount    int
			topicName string
			topicId   string
		)
		err = rows.Scan(&temp.PostType, &temp.OpStatus, &temp.TitleLen, &temp.Vid, &temp.SourceVid, &temp.Title, &temp.VideoCover, &temp.VideoUrl, &temp.Status, &temp.Duration, &qcount, &tcount, &topicName, &topicId)

		if err != nil {
			err = service.MysqlError
			return
		}
		if qcount > 0 {
			temp.HasQuestion = true
		}
		if tcount > 0 {
			temp.HasTopic = true
		}
		if topicName != "" {
			names := strings.Split(topicName, `,`)
			tips := strings.Split(topicId, `,`)
			for k, v := range names {
				topic := api.Topic{
					Name:    v,
					TopicId: tips[k],
				}
				temp.Topics = append(temp.Topics, &topic)
			}
		}
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select v.vid from video_0 v  left join (select vid,count(vid) as qCount  from question group by vid) q on v.vid=q.vid
left join (select vid,count(vid) as tCount,group_concat(t.name) allTopic,group_concat(t.topic_id) allTid  from topic_video_0 tv left join topic t on tv.topic_id=t.topic_id group by vid) t on v.vid=t.vid
where 1=1 %s`, filterSql))
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
func (s *Server) VideoModify(ctx context.Context, req *api.Video) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{
		Code: "200",
		Msg:  "OK",
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	querySql := fmt.Sprintf(`update video_0 set title='%v' where vid=%v`, req.Title, req.Vid)
	_, err = db.Exec(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	return
}
func (s *Server) QuestionList(ctx context.Context, req *api.QueryListRequest) (resp *api.QuestionListResponse, err error) {
	resp = &api.QuestionListResponse{
		Data: &api.QuestionList{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	paramMap := map[string]string{
		"vid":          "q.vid,=%v",
		"id":           "q.id,=%v",
		"questionType": "q.question_type,=%v",
		"category":     "q.question_category,=%v",
		"hardLevel":    "q.hard_level,=%v",
		"state":        "q.state,=%v",
		"cTime":        "q.create_time,>='%v'",
		"sTime":        "q.create_time,>='%v'",
		"eTime":        "q.create_time,<='%v'",
		"rightNum":     "q.right_answer_num",
		"wrongNum":     "q.wrong_answer_num",
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, paramMap)

	querySql := fmt.Sprintf(`select q.answer_type,q.state,q.id,q.content,q.create_time,q.update_time,q.question_category,q.question_type,q.hard_level,q.score,ifnull(c.category_name,''),ifnull(t.type_name,''),q.option_id,q.right_answer_num,q.wrong_answer_num,
ifnull(o.option_id,'') as 'index',o.option_content,o.is_answer,o.answer_num
	from question q 
		left join question_option o on q.id=o.question_id 
			left join question_category c on q.question_category=c.id
	 			left join question_type t on q.question_type=t.id 
	 			where  1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()
	qs := make([]*api.Question, 0, 10)
	questionSet := make(map[string]*api.Question)
	for rows.Next() {
		q := &api.Question{}
		o := &api.QuestionOption{}
		err = rows.Scan(&q.AnswerType, &q.State, &q.Id, &q.Content, &q.CTime, &q.UTime, &q.Category, &q.QuestionType, &q.HardLevel, &q.Score, &q.CategoryName, &q.TypeName, &q.OptionID, &q.RightNum, &q.WrongNum,
			&o.Index, &o.Content, &o.IsAnswer, &o.AnswerNum)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		if _, ok := questionSet[q.Id]; !ok {
			questionSet[q.Id] = q
			qs = append(qs, q)
		}
		tempQ := questionSet[q.Id]
		tempQ.Options = append(tempQ.Options, o)
		if o.IsAnswer == "1" {
			tempQ.Answers = append(tempQ.Answers, o.Content)
		}
	}
	countSql := fmt.Sprintf(`select q.id from question q where 1=1 %s`, filterSql)
	count, err := db.QuerySqlCount(service.BUDAODB, countSql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Data = qs
	resp.Data.Count = count
	return
}

func (s *Server) AddQuestion(ctx context.Context, req *api.Question) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	if req.Score == "" {
		req.Score = "0"
	}
	insertSql := fmt.Sprintf(`insert into question (vid,question_category,question_type,answer_type,hard_level,score,content,option_id)
values(%v,%v,%v,%v,%v,%v,'%v','%v')`, req.Vid, req.Category, req.QuestionType, req.AnswerType, req.HardLevel, req.Score, req.Content, req.OptionID)
	result, err := db.Exec(service.BUDAODB, insertSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	err = insertQuestionOption(fmt.Sprintf(`%v`, id), req.Options...)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	uintVid, err := strconv.ParseUint(req.Vid, 10, 64)
	tableName, _ := db.GetTableName(service.VIDEO_, uintVid)
	db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+1 where vid=%v`, tableName, uintVid))
	return
}
func insertQuestionOption(qid string, options ...*api.QuestionOption) error {
	if len(options) <= 0 {
		return nil
	}
	insertSql := `insert into question_option (question_id,option_id,option_content,is_answer) values `
	for _, o := range options {
		insertSql += fmt.Sprintf(`(%v,%v,'%v',%v),`, qid, o.Index, o.Content, o.IsAnswer)
	}
	insertSql = insertSql[:len(insertSql)-1]
	_, err := db.Exec(service.BUDAODB, insertSql)
	return err
}
func (s *Server) DelQuestion(ctx context.Context, req *api.Question) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	execSql := fmt.Sprintf(`delete from question where id =%v`, req.Id)
	_, err = db.Exec(service.BUDAODB, execSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	db.Exec(service.BUDAODB, fmt.Sprintf(`delete from question_option where question_id =%v`, req.Id))

	return
}
func (s *Server) GetQuestion(ctx context.Context, req *api.Question) (resp *api.GetQuestionResponse, err error) {
	resp = &api.GetQuestionResponse{
		Data: &api.Question{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	selectSql := fmt.Sprintf(`select t.type_name,category_name,q.id,content,options,answer,answer_type,question_category,question_type,hard_level,score from question q
left join question_category c on q.question_category=c.id 
left join question_type t on q.question_type=t.id where q.id=%v`, req.Id)
	rows, err := db.Query(service.BUDAODB, selectSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	q := api.Question{}
	if rows.Next() {
		rows.Scan(&q.TypeName, &q.CategoryName, &q.Id, &q.Content, &q.Options, &q.Answer, &q.AnswerType, &q.Category, &q.QuestionType, &q.HardLevel, &q.Score)
	} else {
		err = fmt.Errorf("query db not data")
		return
	}

	return
}
func (s *Server) ModifyQuestion(ctx context.Context, req *api.Question) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	updateSql := fmt.Sprintf(`update question set question_category=%v,question_type=%v,answer_type=%v,hard_level=%v,score=%v,content='%v',option_id='%v' where id=%v`,
		req.Category, req.QuestionType, req.AnswerType, req.HardLevel, req.Score, req.Content, req.OptionID, req.Id)
	_, err = db.Exec(service.BUDAODB, updateSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	//delete old options
	deleteSql := fmt.Sprintf(`delete from question_option where question_id=%v`, req.Id)
	db.Exec(service.BUDAODB, deleteSql)

	//update options
	err = insertQuestionOption(req.Id, req.Options...)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	return
}

func (s *Server) GiveUp(ctx context.Context, req *api.Video) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	execSql := fmt.Sprintf(`UPDATE video_0 set state=4, op_state=2 WHERE vid=%s`, req.Vid)
	_, err = db.Exec(service.BUDAODB, execSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	return
}

func (s *Server) AuditVideoList(ctx context.Context, req *api.QueryListRequest) (resp *api.VideoListResponse, err error) {
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
		"status": `audit_state,=%v`,
		"cTime":  `insert_time,>='%v'`,
		"uTime":  `update_time,>='%v'`,
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	querySql := fmt.Sprintf(`select v.vid,v.vsource_vid,v.title,v.coverurl,v.videourl,a.audit_state,a.insert_time,a.update_time from (select vid,vsource_vid,audit_state,insert_time,update_time from audit_to_offline where 1=1 %v %v %v)a left join video_0 v  on a.vid=v.vid`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	data := make([]*api.Video, 0, 10)
	for rows.Next() {
		temp := api.Video{}
		err = rows.Scan(&temp.Vid, &temp.VSourceId, &temp.Title, &temp.CoverUrl, &temp.VideoUrl, &temp.Status, &temp.CTime, &temp.UTime)
		if err != nil {
			err = service.MysqlError
			return
		}
		temp.CTime = common.GetTimeStr(temp.CTime)
		temp.UTime = common.GetTimeStr(temp.UTime)
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select vid,vsource_vid,audit_state,insert_time,update_time from audit_to_offline where 1=1 %v`, filterSql))
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
func (s *Server) AuditVideoOffline(ctx context.Context, req *api.AuditVideoOfflineReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	vids := make([]string, 0, len(req.Data))
	vsids := make([]string, 0, len(req.Data))
	for _, v := range req.Data {
		vids = append(vids, v.Vid)
		vsids = append(vsids, `'`+v.VSourceId+`'`)
	}

	//更新视频状态为删除
	updateSql := fmt.Sprintf(`update video_0 set state=4 where vid in (%v)`, strings.Join(vids, `,`))
	_, err = db.Exec(service.BUDAODB, updateSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	//更新视频所在的话题的视频数
	querySql := fmt.Sprintf(`select topic_id,count(id) count from topic_video_0 where vid in (%v) and disable=0 group by topic_id `, strings.Join(vids, `,`))
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	for rows.Next() {
		var (
			tid   string
			count string
		)
		err := rows.Scan(&tid, &count)
		if err != nil {
			continue
		}
		db.Exec(service.BUDAODB, fmt.Sprintf(`update topic set video_num=video_num-%v where topic_id=%v`, count, tid))
	}

	//更新视频在话题下的状态为disable
	updateSql = fmt.Sprintf(`update topic_video_0 set disable=1 where vid in (%v)`, strings.Join(vids, `,`))
	db.Exec(service.BUDAODB, updateSql)

	//更新视频在抓取库的状态
	updateSql = fmt.Sprintf(`update video_data set status=2 where v_source_id in (%v)`, strings.Join(vsids, `,`))
	db.Exec(service.SPIDERDB, updateSql)

	//更新视频在待审核表的状态
	updateSql = fmt.Sprintf(`update audit_to_offline set audit_state=1 where vid in (%v)`, strings.Join(vids, `,`))
	db.Exec(service.BUDAODB, updateSql)

	return
}
func (s *Server) AuditVideoNormal(ctx context.Context, req *api.AuditVideoOfflineReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	vids := make([]string, 0, len(req.Data))
	for _, v := range req.Data {
		vids = append(vids, v.Vid)
	}
	_, err = db.Exec(service.BUDAODB, fmt.Sprintf(`update audit_to_offline set audit_state=2 where vid in (%v)`, strings.Join(vids, `,`)))
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	return
}

func (s *Server) AddTopic(ctx context.Context, req *api.Topic) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	topicId, id, err := common.GetTopicItemId(service.BUDAODB)
	if err != nil {
		glog.Error(err)
		return
	}
	execSql := fmt.Sprintf(`update topic set topic_id=%v,name='%v',pic='%v',disable=%v,need_login=%v,weight=%v,description='%v',rule='%v', hide='%v' where id=%v`,
		topicId, req.Name, req.Pic, req.Disable, req.NeedLogin, req.Weight, req.Description, req.Rule, req.Hide, id)
	_, err = db.Exec(service.BUDAODB, execSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	return
}
func (s *Server) TopicList(ctx context.Context, req *api.QueryListRequest) (resp *api.TopicListResponse, err error) {
	resp = &api.TopicListResponse{
		Data: &api.TopicList{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	var m = map[string]string{
		"needLogin": `need_login,=%v`,
		"disable":   `disable,=%v`,
		"weightS":   `weight,>=%v`,
		"weightE":   `weight,<=%v`,
		"name":      `name,like '%%%v%%'`,
		"cTime":     `create_time,>='%v'`,
		"userNum":   `user_num,>=%v`,
		"videoNum":  `video_num,>=%v`,
		"hide":      `hide,='%v'`,
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	var querySql string
	querySql = fmt.Sprintf(`select hide, id,topic_id,name,pic,disable,need_login,weight,description,rule,user_num,video_num,create_time,update_time from topic where 1=1 %v %v %v`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	data := make([]*api.Topic, 0, 10)
	for rows.Next() {
		temp := api.Topic{}
		err = rows.Scan(&temp.Hide, &temp.Id, &temp.TopicId, &temp.Name, &temp.Pic, &temp.Disable, &temp.NeedLogin, &temp.Weight, &temp.Description, &temp.Rule, &temp.UserNum, &temp.VideoNum, &temp.CTime, &temp.UTime)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select id from topic where 1=1 %s`, filterSql))
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
func (s *Server) TopicInfo(ctx context.Context, req *api.Topic) (resp *api.TopicInfoResponse, err error) {
	resp = &api.TopicInfoResponse{}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	var querySql string
	querySql = fmt.Sprintf(`select hide,id,topic_id,name,pic,disable,need_login,weight,description,rule,user_num,video_num,create_time,update_time from topic where id=%v`, req.Id)
	row, err := db.QueryRow(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return nil, err
	}
	temp := api.Topic{}
	err = row.Scan(&temp.Hide, &temp.Id, &temp.TopicId, &temp.Name, &temp.Pic, &temp.Disable, &temp.NeedLogin, &temp.Weight, &temp.Description, &temp.Rule, &temp.UserNum, &temp.VideoNum, &temp.CTime, &temp.UTime)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data = &temp

	return
}
func (s *Server) TopicModify(ctx context.Context, req *api.Topic) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	var m = map[string]string{
		"name":        `name,='%v'`,
		"pic":         `pic,='%v'`,
		"disable":     `disable,=%v`,
		"needLogin":   `need_login,=%v`,
		"weight":      `weight,=%v`,
		"description": `description,='%v'`,
		"rule":        `rule,='%v'`,
		"hide":        `hide,='%v'`,
	}

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: false}
	err = marshaler.Marshal(&buf, req)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	param, updateSql, err := service.GetUpdateSql(buf.Bytes(), m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	execSql := fmt.Sprintf(`update topic set %v where id=%v`, updateSql, param["id"])
	_, err = db.Exec(service.BUDAODB, execSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	return
}
func (s *Server) TopicAddVideo(ctx context.Context, req *api.TopicVideo) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	tID, err := strconv.ParseUint(req.TopicId, 10, 64)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	tableName, err := db.GetTableName(service.TOPIC_VIDEO_, tID)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	//获取topic
	querySql := fmt.Sprintf(`select topic_id from video_0 where vid =%v`, req.Vid)
	row, err := db.QueryRow(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	var topics string
	err = row.Scan(&topics)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	insertSql := fmt.Sprintf(`insert into %v (topic_id,vid,weight,disable) values (%v,%v,%v,%v)`, tableName, req.TopicId, req.Vid, req.Weight, req.Disable)
	_, err = db.Exec(service.BUDAODB, insertSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	updateSql := fmt.Sprintf(`update topic set video_num=video_num+1 where topic_id=%v`, req.TopicId)
	db.Exec(service.BUDAODB, updateSql)

	//更新topic
	topics = strings.Trim(topics+`,`+req.TopicId, ",")
	updateSql = fmt.Sprintf(`update video_0 set topic_id='%v' where vid=%v`, topics, req.Vid)
	db.Exec(service.BUDAODB, updateSql)

	return
}
func (s *Server) TopicDelVideo(ctx context.Context, req *api.TopicVideo) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	tID, err := strconv.ParseUint(req.TopicId, 10, 64)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	tableName, err := db.GetTableName(service.TOPIC_VIDEO_, tID)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	delSql := fmt.Sprintf(`delete from %v where topic_id=%v and vid=%v`, tableName, req.TopicId, req.Vid)
	_, err = db.Exec(service.BUDAODB, delSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	//更新topic_id字段
	querySql := fmt.Sprintf(`select topic_id from video_0 where vid =%v`, req.Vid)
	row, err := db.QueryRow(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	var topics string
	err = row.Scan(&topics)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	updateSql := fmt.Sprintf(`update topic set video_num=video_num-1 where topic_id=%v`, req.TopicId)
	db.Exec(service.BUDAODB, updateSql)
	tids := strings.Split(topics, ",")
	topics = ""
	for _, tid := range tids {
		if tid == req.TopicId {
			continue
		}
		topics += tid
	}
	topics = strings.Trim(topics, ",")
	updateSql = fmt.Sprintf(`update video_0 set topic_id='%v' where vid=%v`, topics, req.Vid)
	db.Exec(service.BUDAODB, updateSql)

	return
}
func (s *Server) TopicVideoModify(ctx context.Context, req *api.TopicVideo) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()
	tID, err := strconv.ParseUint(req.TopicId, 10, 64)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	tableName, err := db.GetTableName(service.TOPIC_VIDEO_, tID)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	insertSql := fmt.Sprintf(`update %v set weight=%v,disable=%v where topic_id=%v and vid=%v`, tableName, req.Weight, req.Disable, req.TopicId, req.Vid)
	_, err = db.Exec(service.BUDAODB, insertSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	return
}
func (s *Server) TopicVideoList(ctx context.Context, req *api.QueryListRequest) (resp *api.VideoListResponse, err error) {
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
		"vid":       `v.vid,='%v'`,
		"ruleID":    `t.rule_id,='%v'`,
		"status":    `v.state,=%v`,
		"topicId":   `t.topic_id,=%v`,
		"weight":    `t.weight,>=%v`,
		"disable":   `t.disable,=%v`,
		"durationS": `v.duration,>=%v`,
		"durationE": `v.duration,<=%v`,
		"title":     `v.title,like '%%%v%%'`,
		"cTime":     `v.create_time,>='%v'`,
		"titleLenS": "v.char_length(title),>=%v",
		"titleLenE": "v.char_length(title),<=%v",
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	tid, err := strconv.ParseUint(req.Filter["topicId"], 10, 64)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	tableName, err := db.GetTableName(service.TOPIC_VIDEO_, tid)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	querySql := fmt.Sprintf(`select char_length(title) as title_len,v.vid,source_vid,title,coverurl,videourl,state,duration,t.weight,t.disable from %v t left join video_0 v on t.vid=v.vid where 1=1 %v %v %v`, tableName, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	data := make([]*api.Video, 0, 10)
	for rows.Next() {
		temp := api.Video{}
		rows.Scan(&temp.TitleLen, &temp.Vid, &temp.SourceVid, &temp.Title, &temp.VideoCover, &temp.VideoUrl, &temp.Status, &temp.Duration, &temp.Weight, &temp.Disable)
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.BUDAODB, fmt.Sprintf(`select v.vid from %v t left join video_0 v on t.vid=v.vid where 1=1 %v `, tableName, filterSql))
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

func (s *Server) CommentList(ctx context.Context, req *api.CommentListReq) (resp *api.CommentListResponse, err error) {
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

	allDate := make(map[string]*api.Comment)
	masterDate := make([]*api.Comment, 0, 50)
	maxNum := common.GetConfig().DB.MySQL[service.BUDAODB].TableDesc[service.COMMENT_]
	if data, ok := videoCommentCache.get(req.Vid); ok {
		masterDate = data
		goto returnData
	}

	for i := uint64(0); i < maxNum; i++ {
		tableName := fmt.Sprintf(service.COMMENT_+"%v", i)
		querySql := fmt.Sprintf(`select cid,vid,from_uid,from_name,from_photo,to_comment_id,parentcomid,content,favor_num,weight,reply_num,state,create_time,update_time from %v where vid=%v`, tableName, req.Vid)
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
			if temp.PId == "0" {
				masterDate = append(masterDate, temp)
			}
			allDate[temp.Cid] = temp
		}
	}
	for _, v := range allDate {
		if v.PId != "0" {
			temp := allDate[v.PId]
			temp.Sub = append(temp.Sub, v)
			allDate[v.PId] = temp
		}
	}
	videoCommentCache.put(req.Vid, masterDate)
returnData:
	count := int32(len(masterDate))
	start := (req.Num - 1) * req.Size
	end := start + req.Size
	if end > count {
		end = count
	}
	if start >= count {
		start = 0
		end = req.Size
	}
	data := masterDate[start:end]
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Data = data
	resp.Data.Count = uint64(len(masterDate))
	return
}
func (s *Server) CommentModifyState(ctx context.Context, req *api.CommentModifyStateReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	tableMap := make(map[string][]string)
	for _, cid := range req.Cids {
		id, _ := strconv.ParseUint(cid, 10, 64)
		tableName, err := db.GetTableName(service.COMMENT_, id)
		if err != nil {
			glog.Error(err)
			return resp, err
		}
		tableMap[tableName] = append(tableMap[tableName], cid)
	}

	for tableName, cids := range tableMap {
		updateSql := fmt.Sprintf(`update %v set state=%v where cid in (%v)`, tableName, req.State, strings.Join(cids, `,`))
		_, err = db.Exec(service.BUDAODB, updateSql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
	}
	return
}
func (s *Server) CommentModify(ctx context.Context, req *api.Comment) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	var m = map[string]string{
		"content":  `content,='%v'`,
		"cTime":    `create_time,='%v'`,
		"state":    `state,=%v`,
		"favorNum": `favor_num,=%v`,
		"weight":   `weight,=%v`,
	}
	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: false}
	err = marshaler.Marshal(&buf, req)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	param, updateSql, err := service.GetUpdateSql(buf.Bytes(), m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	cid, err := strconv.ParseUint(param["cid"], 10, 64)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	tableName, err := db.GetTableName(service.COMMENT_, cid)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	execSql := fmt.Sprintf(`update %v set %v where cid=%v`, tableName, updateSql, cid)
	_, err = db.Exec(service.BUDAODB, execSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	if v, ok := param["vid"]; ok {
		videoCommentCache.del(v)
	}

	return
}
func (s *Server) CommentInfo(ctx context.Context, req *api.Comment) (resp *api.CommentInfoResponse, err error) {
	resp = &api.CommentInfoResponse{}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	cid, err := strconv.ParseUint(req.Cid, 10, 64)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	tableName, err := db.GetTableName(service.COMMENT_, cid)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	var querySql string
	querySql = fmt.Sprintf(`select cid,vid,from_uid,from_name,from_photo,to_comment_id,parentcomid,content,favor_num,weight,reply_num,state,create_time,update_time from %v where cid=%v`, tableName, req.Cid)
	row, err := db.QueryRow(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	temp := api.Comment{}
	err = row.Scan(&temp.Cid, &temp.Vid, &temp.Uid, &temp.UName, &temp.UPhoto, &temp.PId, &temp.MasterId, &temp.Content, &temp.FavorNum, &temp.Weight, &temp.ReplyNum, &temp.State, &temp.CTime, &temp.UTime)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data = &temp
	return
}
func (s *Server) CommentReply(ctx context.Context, req *api.Comment) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	uintVid, _ := strconv.ParseUint(req.Vid, 10, 64)
	cid, autoIncreID, tableName, err := common.GetItemId(service.BUDAODB, service.COMMENT_)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	updateSql := fmt.Sprintf(`update %v set content="%v",cid=%v,vid=%v,from_uid=%v,from_name="%v",from_photo="%v",to_comment_id=%v,parentcomid=%v where id=%v`,
		tableName, req.Content, cid, req.Vid, req.Uid, req.UName, req.UPhoto, req.PId, req.MasterId, autoIncreID)
	_, err = db.Exec(service.BUDAODB, updateSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	tableName, _ = db.GetTableName(service.VIDEO_, uintVid)
	db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+1 where vid=%v`, tableName, uintVid))
	uintUid, _ := strconv.ParseUint(req.Uid, 10, 64)
	tableName, _ = db.GetTableName(service.USER_, uintUid)
	db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+1 where uid=%v`, tableName, uintUid))

	return
}
func (s *Server) CommentAdd(ctx context.Context, req *api.CommentAddReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	uintVid, err := strconv.ParseUint(req.Vid, 10, 64)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	users := service.GetRandUser(len(req.Comments))
	for index, comment := range req.Comments {
		cid, autoIncreId, tableName, err := common.GetItemId(service.BUDAODB, service.COMMENT_)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}
		u := users[index]
		updateSql := fmt.Sprintf(`update %v set content="%v",cid=%v,vid=%v,from_uid=%v,from_name="%v",from_photo="%v" where id=%v`, tableName, comment, cid, req.Vid, u.Uid, u.Name, u.Photo, autoIncreId)
		_, err = db.Exec(service.BUDAODB, updateSql)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return resp, err
		}

		tableName, _ = db.GetTableName(service.VIDEO_, uintVid)
		db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+1 where vid=%v`, tableName, uintVid))
		uintUid, _ := strconv.ParseUint(u.Uid, 10, 64)
		tableName, _ = db.GetTableName(service.USER_, uintUid)
		db.Exec(service.BUDAODB, fmt.Sprintf(`update %v set comment_num=comment_num+1 where uid=%v`, tableName, uintUid))

	}
	return
}

func (s *Server) VideoReview(ctx context.Context, req *api.VideoReviewReq) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	var videoStatus string
	if req.OpStatus == "0" {
		videoStatus = "2"
	} else {
		videoStatus = "4"
	}
	execSql := fmt.Sprintf(`UPDATE video_0 set op_state = '%s', state = '%s' WHERE vid = %s`, req.OpStatus, videoStatus, req.Vid)
	_, err = db.Exec(service.BUDAODB, execSql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	return
}
func (s *Server) GetTopicRule(ctx context.Context, req *api.GetTopicRuleReq) (resp *api.GetTopicRuleResponse, err error) {
	resp = &api.GetTopicRuleResponse{}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	querySql := fmt.Sprintf(`select id, name from topic_rule where topic_id = '%v'`, req.TopicId)
	glog.Info(querySql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()

	data := make([]*api.RuleName, 0, 10)
	for rows.Next() {
		temp := api.RuleName{}
		err = rows.Scan(&temp.RuleId, &temp.RuleName)
		if err != nil {
			glog.Error(err)
			err = service.ParamError
			return
		}
		data = append(data, &temp)
	}
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data = data
	return
}
func (s *Server) GetTopicName(ctx context.Context, req *api.NullReq) (resp *api.TopicListResponse, err error) {
	resp = &api.TopicListResponse{
		Data: &api.TopicList{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
			err = nil
		}
	}()

	querySql := fmt.Sprintf(`select topic_id, name from topic where hide = 0`)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	data := make([]*api.Topic, 0, 10)
	for rows.Next() {
		temp := api.Topic{}
		err = rows.Scan(&temp.TopicId, &temp.Name)
		if err != nil {
			glog.Error(err)
			err = service.MysqlError
			return
		}
		data = append(data, &temp)
	}
	resp.Code = "200"
	resp.Msg = "OK"
	resp.Data.Data = data
	return
}

//视频评论缓存
var videoCommentCache *commentCache = func() *commentCache {

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

func (c commentCache) get(vid string) ([]*api.Comment, bool) {
	c.Lock()
	defer c.Unlock()
	data, ok := c.data[vid]
	if !ok {
		return nil, ok
	}
	now := time.Now()
	if now.Sub(data.inTime) > 5*time.Minute {
		delete(c.data, vid)
		return nil, false
	}
	return data.data, ok
}
func (c commentCache) del(vid string) {
	c.Lock()
	defer c.Unlock()
	delete(c.data, vid)
}

func (c commentCache) put(vid string, data []*api.Comment) {
	c.Lock()
	defer c.Unlock()
	cacheDate := commentCacheDate{
		inTime: time.Now(),
		data:   data,
	}
	c.data[vid] = &cacheDate
}

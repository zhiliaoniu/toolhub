package topicservice

import (
	"context"
	"db"
	"encoding/json"
	"fmt"
	"github.com/sumaig/glog"
	"math"
	"regexp"
	"service"
	"service/api"
	"strconv"
	"strings"
	"time"
)

func GetServer() *Server {
	return &Server{}
}

type Server struct {
}

//发布
func (s *Server) PostRuleVideo(ctx context.Context, req *api.TopicRulePostVideoRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	rSql := fmt.Sprintf(`select id, topic_id, name, rules from topic_rule where state = 0 and id = '%v'`, req.RuleID)
	row, err := db.QueryRow(service.BUDAODB, rSql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	temp := api.Rule{}
	err1 := row.Scan(&temp.Id, &temp.TopicId, &temp.RuleName, &temp.RuleCondition)
	if err1 != nil {
		glog.Error(err1)
		return
	}

	//对条件进行处理
	var sql string
	m := make(map[string]string)
	err = json.Unmarshal([]byte(temp.RuleCondition), &m)
	if err != nil {
		glog.Error(err)
		return
	}
	for _, v := range m {
		//将值拼起来组成一个sql
		sql += v + " and "
	}
	sql = sql + " 1=1"

	topicId := temp.TopicId
	ruleId := req.RuleID

	total, err := db.QuerySqlCount(service.SPIDERDB, fmt.Sprintf(`select id from video_data where %s`, sql))
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	limit := 100
	num := math.Ceil(float64(total) / float64(limit))
	idSql := fmt.Sprintf(`select id from video_data where %s`, sql)
	querySql := fmt.Sprintf(`select id, i_id, ifnull(v_source_id,""), topic, video_title, play_count, video_duration, share_url as video_url, video_cover, source, video_width, video_height, parse_type from video_data where %s`, sql)

	for i := 0; i < int(num); i++ {
		offest := i * limit
		go func() {
			lSql := fmt.Sprintf(`%s and id >= ( %s limit %s, 1) limit %s`, querySql, idSql, strconv.Itoa(offest), strconv.Itoa(limit))
			rows, err := db.Query(service.SPIDERDB, lSql)
			if err != nil {
				err = service.MysqlError
				return
			}
			defer rows.Close()

			for rows.Next() {
				temp := api.Video{}
				err := rows.Scan(&temp.Id, &temp.SourceVid, &temp.VSourceId, &temp.Topic, &temp.Title, &temp.PlayCount, &temp.VideoDuration, &temp.VideoUrl, &temp.VideoCover, &temp.Source, &temp.VideoWidth, &temp.VideoHeight, &temp.ParseType)
				if err != nil {
					glog.Error(err)
					err = service.MysqlError
					continue
				}
				//获取每个视频的信息，开始发布视频、评论
				res := service.PostVideo(temp)
				if len(res.GetTResults()) > 0 {
					result := res.GetTResults()[0]
					data := struct {
						Vid string `json:"vid"`
					}{result.Vid}
					if data.Vid != "" {
						glog.Info(result)
					} else {
						//更新video_data发布失败原因
						//updateVideoData(temp.SourceVid, fmt.Sprintf("retun result: %d", pb.EPostResult(result.Result)), 3)
						err = service.HasPostError
						continue
					}
					glog.Info(data)
					//1.获取到视频的vid,修改爬虫状态
					updateSql := fmt.Sprintf(`update video_data set status=1 and post_vid=%v where i_id=%v`, data.Vid, temp.SourceVid)
					db.Exec(service.SPIDERDB, updateSql)
					//2.将发布的视频放入topic_video_0里面
					service.PostVideoTopic(topicId, data.Vid, ruleId)
					//3.发布视频下的评论
					service.PostVideoComment(temp.VSourceId, data.Vid)
					//4.更新video_0里面记录，以哪种方式存入的(规则，手动)
					db.Exec(service.BUDAODB, fmt.Sprintf(`update video_0 set post_type = 1 where vid='%v'`, data.Vid))

				} else {
					//更新video_data发布失败原因
					//updateVideoData(temp.SourceVid, "server not return value", 3)
					err = fmt.Errorf("server not return value")
				}
			}
		}()
	}
	//更新规则发布数
	updateSql := fmt.Sprintf(`update topic_rule set post_count=post_count+1 where id=%v`, req.RuleID)
	db.Exec(service.BUDAODB, updateSql)

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"

	return
}

//规则浏览列表
func (s *Server) RuleViewList(ctx context.Context, req *api.QueryListRequest) (resp *api.TopicRuleViewListResponse, err error) {
	resp = &api.TopicRuleViewListResponse{
		Data: &api.ViewList{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	var m = map[string]string{
		"media_name":     `media_name,='%v'`,
		"titleLenS":      "char_length(video_title),>=%v",
		"titleLenE":      "char_length(video_title),<=%v",
		"play_count":     `play_count,>='%v'`,
		"video_duration": `video_duration,>='%v'`,
		"source":         `source,='%v'`,
		"source_type":    `source_type,='%v'`,
		"praise_count":   `praise_count,>='%v'`,
		"share_count":    `share_count,>='%v'`,
		"fav_count":      `fav_count,>='%v'`,
		"comment_count":  `comment_count,>='%v'`,
		//"title":         `video_title,like '%%%v%%'`,
	}

	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	querySql := fmt.Sprintf(`select char_length(video_title) as title_len,ifnull(video_duration,0),topic,id,video_title,play_count,share_url as video_url,video_cover,i_id,status,source,video_width,video_height, ifnull(source_type,""), praise_count, share_count, fav_count, comment_count from video_data where 1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.SPIDERDB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()

	data := make([]*api.Video, 0, 10)
	for rows.Next() {
		temp := api.Video{}
		err := rows.Scan(&temp.TitleLen, &temp.Duration, &temp.Topic, &temp.Id, &temp.Title, &temp.PlayCount, &temp.VideoUrl, &temp.VideoCover, &temp.Vid, &temp.Status, &temp.Source, &temp.VideoWidth, &temp.VideoHeight, &temp.SourceType, &temp.PraiseCount, &temp.ShareCount, &temp.FavCount, &temp.CommentCount)
		if err != nil {
			glog.Error(err)
			err = service.ParamError
			break
		}
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.SPIDERDB, fmt.Sprintf(`select id from video_data where 1=1 %s`, filterSql))
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

//获取媒介名称
func (s *Server) RuleMediaName(ctx context.Context, req *api.QueryListRequest) (resp *api.TopicRuleMediaNameResponse, err error) {
	resp = &api.TopicRuleMediaNameResponse{
		Data: &api.MediaName{},
	}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	var m = map[string]string{
		"mediaName": `media_name,like '%%%v%%'`,
		"source":    `source,='%v'`,
	}

	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}
	querySql := fmt.Sprintf(`select media_name, source from video_data group by media_name having 1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.SPIDERDB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()

	type MediaName struct {
		MediaName string `json:"mediaName"`
		Source    string `json:"source"`
	}

	data := make([]*api.MediaNameMessage, 0, 10)
	for rows.Next() {
		temp := api.MediaNameMessage{}
		err := rows.Scan(&temp.MediaName, &temp.Source)
		if err != nil {
			glog.Error(err)
			err = service.ParamError
			break
		}
		data = append(data, &temp)
	}
	count, err := db.QuerySqlCount(service.SPIDERDB, fmt.Sprintf(`select media_name from video_data group by media_name having 1=1 %s`, filterSql))
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

//rule修改
func (s *Server) RuleModify(ctx context.Context, req *api.TopicRuleRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	b, _ := json.Marshal(req)
	m := make(map[string]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	rm := make(map[string]string)
	for k, v := range m {
		if len(v) > 0 {
			switch {
			case k == "media_name" || k == "source" || k == "source_type":
				rm[k] = k + " = " + "'" + v + "'"
			case k == "titleLenE":
				rm[k] = "char_length(video_title)" + " <= " + "'" + v + "'"
			case k == "titleLenS":
				rm[k] = "char_length(video_title)" + " >= " + "'" + v + "'"
			case k == "id" || k == "name" || k == "topicID" || k == "topicName":
				continue
			default:
				rm[k] = k + " >= " + "'" + v + "'"
			}
		}
	}

	str, err := json.Marshal(rm)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	ruleStr := strings.Replace(strings.Replace(string(str), "\\u003e", ">", -1), "\\u003c", "<", -1)
	rule := strings.TrimSpace(strings.Replace(ruleStr, "'", "\\'", -1))
	addQuerySql := fmt.Sprintf(`update topic_rule set topic_id='%v', name='%v', rules='%v' where id='%v'`, m["topicID"], m["name"], rule, m["id"])
	_, err = db.Exec(service.BUDAODB, addQuerySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"

	return
}

//rule信息
func (s *Server) RuleInfo(ctx context.Context, req *api.TopicRuleInfoRequest) (resp *api.TopicRuleInfoResponse, err error) {
	resp = &api.TopicRuleInfoResponse{}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()

	querySql := fmt.Sprintf(`select r.id, r.topic_id, r.name, r.rules, r.state, t.name as topic_name from topic_rule as r inner join topic as t on r.topic_id = t.topic_id where r.id=%v`, req.Id)
	row, err := db.QueryRow(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	temp := api.Rule{}
	err = row.Scan(&temp.Id, &temp.TopicId, &temp.RuleName, &temp.RuleCondition, &temp.RuleStatus, &temp.TopicName)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	//对返回前端的json进行处理
	m := make(map[string]string)
	err1 := json.Unmarshal([]byte(temp.RuleCondition), &m)
	if err1 == nil {
		for k, v := range m {
			pat := "=(.)+"
			reg := regexp.MustCompile(pat)
			value := strings.Replace(reg.FindString(v), "=", "", -1)
			m[k] = strings.TrimSpace(strings.Replace(value, "'", "", -1))
		}
	}

	m["name"] = temp.RuleName
	m["id"] = temp.Id
	m["status"] = temp.RuleStatus
	m["topicID"] = temp.TopicId
	m["topicName"] = temp.TopicName

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	resp.Data = m
	return
}

//删除
func (s *Server) RuleDel(ctx context.Context, req *api.TopicRuleDelRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	//更新topic_rule的规则状态
	querySql := fmt.Sprintf(`update topic_rule set state = 1 where id='%v'`, req.Id)
	_, err = db.Exec(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	gSql := fmt.Sprintf(`select topic_id from topic_rule where id = '%v'`, req.Id)
	row, err := db.QueryRow(service.BUDAODB, gSql)
	if err != nil {
		glog.Error(err)
		return nil, service.MysqlError
	}
	var topicId string
	err = row.Scan(&topicId)
	glog.Error(err)
	//将规则下的所有发布视频进行删除
	tID, _ := strconv.ParseUint(topicId, 10, 64)
	tableName, err := db.GetTableName(service.TOPIC_VIDEO_, tID)
	glog.Error(err)

	tSql := fmt.Sprintf(`update %v set disable = 1 where rule_id = '%v'`, tableName, req.Id)
	_, err = db.Exec(service.BUDAODB, tSql)
	glog.Error(err)
	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"

	return
}

//rule列表
func (s *Server) RuleList(ctx context.Context, req *api.QueryListRequest) (resp *api.TopicRuleListResponse, err error) {
	resp = &api.TopicRuleListResponse{
		Data: &api.TopicRuleList{},
	}

	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	var m = map[string]string{
		"topicID": `r.topic_id,='%v'`,
		"name":    `r.name,='%v'`,
		"status":  `r.state,='%v'`,
	}
	filterSql, sortSql, pageSql, err := service.GetSqlParam(*req, m)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	querySql := fmt.Sprintf(`select r.post_count, r.id, r.topic_id, r.name, r.rules, r.state, t.name as topic_name from topic_rule as r inner join topic as t on r.topic_id = t.topic_id where 1=1 %s %s %s`, filterSql, sortSql, pageSql)
	rows, err := db.Query(service.BUDAODB, querySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}
	defer rows.Close()

	data := make([]*api.TopicBack, 0, 10)
	for rows.Next() {
		temp := api.Rule{}
		err := rows.Scan(&temp.PostCount, &temp.Id, &temp.TopicId, &temp.RuleName, &temp.RuleCondition, &temp.RuleStatus, &temp.TopicName)
		if err != nil {
			err = service.MysqlError
			break
		}
		//对返回前端的json进行处理
		m := api.TopicBack{}

		err1 := json.Unmarshal([]byte(temp.RuleCondition), &m.Data)
		if err1 == nil {
			for k, v := range m.Data {
				pat := "=(.)+"
				reg := regexp.MustCompile(pat)
				value := strings.Replace(reg.FindString(v), "=", "", -1)
				m.Data[k] = strings.TrimSpace(strings.Replace(value, "'", "", -1))
			}
		}
		m.Data["name"] = temp.RuleName
		m.Data["id"] = temp.Id
		m.Data["status"] = temp.RuleStatus
		m.Data["topicID"] = temp.TopicId
		m.Data["topName"] = temp.TopicName
		m.Data["postCount"] = temp.PostCount

		rules, _ := json.Marshal(m)
		temp.RuleCondition = string(rules)
		data = append(data, &m)
	}

	csql := fmt.Sprintf(`select id from topic_rule where 1=1 %s`, filterSql)
	count, err := db.QuerySqlCount(service.BUDAODB, csql)
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

//添加rule
func (s *Server) RuleAdd(ctx context.Context, req *api.TopicRuleRequest) (resp *api.CommonResponse, err error) {
	resp = &api.CommonResponse{Code: "200", Msg: "OK"}
	defer func() {
		if err != nil {
			resp.Code = service.FAIL_CODE
			resp.Msg = fmt.Sprintf("%s", err)
		}
	}()
	b, _ := json.Marshal(req)
	m := make(map[string]string)
	err = json.Unmarshal(b, &m)

	rm := make(map[string]string)
	for k, v := range m {
		if len(v) > 0 {
			switch {
			case k == "media_name" || k == "source" || k == "source_type":
				rm[k] = k + " = " + "'" + v + "'"
			case k == "titleLenE":
				rm[k] = "char_length(video_title)" + " <= " + "'" + v + "'"
			case k == "titleLenS":
				rm[k] = "char_length(video_title)" + " >= " + "'" + v + "'"
			case k == "id" || k == "name" || k == "topicID" || k == "topicName":
				continue
			default:
				rm[k] = k + " >= " + "'" + v + "'"
			}
		}
	}

	str, err := json.Marshal(rm)
	if err != nil {
		glog.Error(err)
		err = service.ParamError
		return
	}

	ruleStr := strings.Replace(strings.Replace(string(str), "\\u003e", ">", -1), "\\u003c", "<", -1)
	rule := strings.TrimSpace(strings.Replace(ruleStr, "'", "\\'", -1))
	addQuerySql := fmt.Sprintf(`insert into topic_rule(topic_id, name, rules) values ('%v', '%v', '%v')`, m["topicID"], m["name"], rule)
	_, err = db.Exec(service.BUDAODB, addQuerySql)
	if err != nil {
		glog.Error(err)
		err = service.MysqlError
		return
	}

	resp.Code = service.SUCCESS_CODE
	resp.Msg = "OK"
	return
}

func timPostRuleVideo() {
	/**
	 * 1.从topic_rule表中获取规则
	 * 2.话题的规则下的视频进行发布，和评论发布
	 */
	for {
		querySql := fmt.Sprintf(`select id, topic_id, name, rules from topic_rule where state = 0`)
		rows, err := db.Query(service.BUDAODB, querySql)
		if err != nil {
			glog.Error(err)
		}
		defer rows.Close()

		for rows.Next() {
			temp := api.Rule{}
			err := rows.Scan(&temp.Id, &temp.TopicId, &temp.RuleName, &temp.RuleCondition)
			if err != nil {
				glog.Error(err)
				return
			}

			//对返回前端的json进行处理
			var str string
			m := make(map[string]string)
			err1 := json.Unmarshal([]byte(temp.RuleCondition), &m)
			if err1 == nil {
				for _, v := range m {
					//将值拼起来组成一个sql
					str += v + " and "
				}
				str = str + " 1=1"
				total, _ := db.QuerySqlCount(service.SPIDERDB, fmt.Sprintf(`select id from video_data where %s`, str))

				limit := 100
				num := math.Ceil(float64(total) / float64(limit))
				idSql := fmt.Sprintf(`select id from video_data where %s`, str)
				querySql := fmt.Sprintf(`select id, i_id, ifnull(v_source_id,""), topic, video_title, play_count, video_duration, video_url, video_cover, source, video_width, video_height from video_data where %s`, str)
				for i := 0; i < int(num); i++ {
					offest := i * limit
					go func() {
						lSql := fmt.Sprintf(`%s and id >= ( %s limit %s, 1) limit %s`, querySql, idSql, strconv.Itoa(offest), strconv.Itoa(limit))
						rows, err := db.Query(service.SPIDERDB, lSql)
						glog.Error(err)
						defer rows.Close()

						for rows.Next() {
							temp := api.Video{}
							err := rows.Scan(&temp.Id, &temp.SourceVid, &temp.VSourceId, &temp.Topic, &temp.Title, &temp.PlayCount, &temp.VideoDuration, &temp.VideoUrl, &temp.VideoCover, &temp.Source, &temp.VideoWidth, &temp.VideoHeight)
							glog.Error(err)

							//获取每个视频的信息，开始发布视频、评论
							res := service.PostVideo(temp)
							if len(res.GetTResults()) > 0 {
								result := res.GetTResults()[0]
								data := struct {
									Vid string `json:"vid"`
								}{result.Vid}
								if data.Vid != "" {
									glog.Info(data)
								} else {
									glog.Info(data)
									continue
								}
								//1.获取到视频的vid,修改爬虫状态
								updateSql := fmt.Sprintf(`update video_data set status=1 and post_vid=%v where i_id=%v`, data.Vid, temp.SourceVid)
								db.Exec(service.SPIDERDB, updateSql)
								//2.将发布的视频放入topic_video_0里面
								service.PostVideoTopic(m["topicID"], data.Vid, m["ruleID"])
								//3.发布视频下的评论
								service.PostVideoComment(temp.VSourceId, data.Vid)
								//4.更新video_0里面记录，以哪种方式存入的(规则，手动)
								db.Exec(service.BUDAODB, fmt.Sprintf(`update video_0 set post_type = 1 where vid='%v'`, data.Vid))

							} else {
								glog.Error("server not return value")
							}
						}
					}()
				}
			}
		}

		//睡一会
		time.Sleep(1 * time.Minute)
	}

}

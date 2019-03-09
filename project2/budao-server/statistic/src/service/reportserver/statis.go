package reportserver

import (
	"common"
	"db"
	"fmt"

	"github.com/robfig/cron"
	"github.com/sumaig/glog"
)

//cron使用文档：https://godoc.org/github.com/robfig/cron
/*
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/

func (s *Server) cronTask() {
	//TestStatis()
	s.c = cron.New()
	s.c.AddFunc("0 5 0 * * *", StatisUser)
	s.c.AddFunc("0 6 0 * * *", StatisBizDaily)
	s.c.Start()
}

func TestStatis() {
	StatisUser()
	StatisBizDaily()
}

func StatisUser() {
	glog.Debug("begin Statis User")
	var totalUserCount, yestordayActiveUserCount, newUserCount, oldUserCount int
	var err error
	//1.统计截止今天凌晨的用户总量
	if totalUserCount, err = GetDintinctRecordNumByTime("user_statistic_", "dev_info", "create_time", "", GetYesterdayEnd()); err != nil {
		return
	}
	//2.统计昨天的活跃数
	if yestordayActiveUserCount, err = GetYesterdayDintinctRecordNum("user_statistic_", "dev_info"); err != nil {
		return
	}
	//3.统计截止昨天凌晨的用户总增量
	if oldUserCount, err = GetBeforeYesterdayDintinctRecordNum("user_statistic_", "dev_info"); err != nil {
		return
	}

	newUserCount = totalUserCount - oldUserCount
	execSql := fmt.Sprintf("insert into statis_new_user(total_num, active_num, new_num, create_time) values(%d, %d, %d, '%s')", totalUserCount, yestordayActiveUserCount, newUserCount, GetYesterday())
	_, err = db.Exec(common.BUDAODB, execSql)
	if err != nil {
		glog.Error("insert mysql failed. execSql:%s, err:%v", execSql, err)
		return
	}

	glog.Debug("end Statis User")
}

func StatisBizDaily() {
	glog.Debug("begin StatisBizDaily")
	//1.视频
	//1.1 视频曝光
	videoExposeNum, err := GetYesterdayVideoRecordNumWithFlag("user_statistic_", 3)
	if err != nil {
		return
	}
	//1.2 视频点击
	videoClickNum, err := GetYesterdayVideoRecordNumWithFlag("user_statistic_", 1)
	if err != nil {
		return
	}
	//1.3 视频播放
	videoViewNum, err := GetYesterdayVideoRecordNumWithFlag("user_statistic_", 2)
	if err != nil {
		return
	}
	//1.4 视频点赞
	videoFavorNum, err := GetYesterdayRecordNumWithName("user_favor_video_", "updatetime")
	if err != nil {
		return
	}
	//1.5 视频播放时长
	videoPlayTime, err := GetYesterdayVideoViewTimeWithFlag("user_statistic_", 2)
	if err != nil {
		return
	}
	//2.评论
	//2.1 评论点赞
	commentFavorNum, err := GetYesterdayRecordNumWithName("user_favor_comment_", "updatetime")
	if err != nil {
		return
	}
	//2.1 评论发表
	commentNum, err := GetYesterdayRecordNum("comment_")
	if err != nil {
		return
	}
	//3.话题
	//3.1 话题关注数
	topicFollowNum, err := GetYesterdayRecordNumWithName("user_follow_topic_", "update_time")
	if err != nil {
		return
	}
	execSql := fmt.Sprintf("insert into statis_biz_daily(video_expose_num, video_click_num, video_view_num, video_favor_num, video_play_time, comment_favor_num, comment_num, topic_follow_num, create_time) values(%d,%d,%d,%d,%d,%d,%d,'%s')", videoExposeNum, videoClickNum, videoViewNum, videoFavorNum, videoPlayTime, commentFavorNum, commentNum, topicFollowNum, GetYesterday())
	_, err = db.Exec(common.BUDAODB, execSql)
	if err != nil {
		glog.Error("insert mysql failed. execSql:%s, err:%v", execSql, err)
		return
	}

	glog.Debug("end StatisBizDaily")
}

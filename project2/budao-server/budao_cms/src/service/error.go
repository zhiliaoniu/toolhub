package service

import "github.com/twitchtv/twirp"

var (
	ParamError   twirp.Error = twirp.NewError(twirp.NoError, "参数错误")
	MysqlError   twirp.Error = twirp.NewError(twirp.NoError, "数据库操作错误")
	HasPostError twirp.Error = twirp.NewError(twirp.NoError, "视频已经发布过")
	IoError      twirp.Error = twirp.NewError(twirp.NoError, "系统IO错误")
)

const USER_ = "user_"
const COMMENT_ = "comment_"
const TOPIC_VIDEO_ = "topic_video_"
const VIDEO_ = "video_"

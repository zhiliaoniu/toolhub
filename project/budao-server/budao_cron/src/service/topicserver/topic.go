package topicserver

import "common"

// TopicInfo storage topic info from mysql
type TopicInfo struct {
	AutoID      uint64
	TopicID     uint64
	Name        string
	Picture     string
	Disable     uint32
	Weight      uint32
	Desc        string
	UserNUM     uint32
	FakeUserNUM uint32
	VideoNUM    uint64
}

// Server identify for TopicService RPC
type Server struct {
	cronTaskDisable     bool //定时任务是否可用
	cronTaskInternalSec int  //定时任务执行间隔
}

// GetServer get server
func GetServer() *Server {
	server := &Server{}
	server.initServer()

	//do cron
	go server.cronTask()

	return server
}

func (s *Server) initServer() {
	topicConf := common.GetConfig().TopicConf
	s.cronTaskDisable = topicConf.CronTaskDisable
	s.cronTaskInternalSec = topicConf.CronTaskInternalSec

	//do init
	go InitUserSubscribedTopic()
	go InitTopicDynamicInfo()
}

// Close server
func (s *Server) Close() {
	s.cronTaskDisable = true
}

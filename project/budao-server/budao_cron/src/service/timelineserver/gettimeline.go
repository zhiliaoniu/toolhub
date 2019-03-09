package timelineserver

import (
	"common"
)

type Server struct {
	cronTaskDisable     bool
	cronTaskInternalSec int
}

func GetServer() *Server {
	server := &Server{}
	server.initServer()

	// cron update listitem by vid
	go server.CronTask()

	return server
}

func (s *Server) initServer() {
	timelineConf := common.GetConfig().TimelineConf
	s.cronTaskDisable = timelineConf.CronTaskDisable
	s.cronTaskInternalSec = timelineConf.CronTaskInternalSec

	go InitVideoDynamicInfo()
	go InitUserAct()
}

func (s *Server) Close() {
	s.cronTaskDisable = true
}

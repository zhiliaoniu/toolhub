package commentserver

import (
	"common"
	"time"
)

func GetServer() *Server {
	server := &Server{}
	server.initServer()
	return server
}

// Server identify for comment RPC
type Server struct {
	cronTaskDisable     bool //定时任务是否可用
	cronTaskInternalSec int  //评论相关定时任务执行间隔
}

func (s *Server) initServer() {
	commentConf := common.GetConfig().CommentConf
	s.cronTaskDisable = commentConf.CronTaskDisable
	s.cronTaskInternalSec = commentConf.CronTaskInternalSec

	go InitCommentInfoByTraverseMysql()

	go s.cronTask()
}

func (s *Server) Close() {
	s.cronTaskDisable = true
}

func (s *Server) cronTask() {
	ticker := time.NewTicker(time.Minute * time.Duration(s.cronTaskInternalSec))
	for {
		if s.cronTaskDisable {
			break
		}
		select {
		case <-ticker.C:
			UpdateCommentInfoByTraverseMysql()
		case <-time.After(time.Second):
		}
	}
}

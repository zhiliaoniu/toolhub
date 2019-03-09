package pushserver

import (
	"github.com/robfig/cron"
)

type Server struct {
	c *cron.Cron
}

func GetServer() *Server {
	server := &Server{}
	server.initServer()
	return server
}

func (s *Server) initServer() {
	go s.cronTask()
}

func (s *Server) Close() {
	s.c.Stop()
}

func (s *Server) cronTask() {
	s.c = cron.New()
	s.c.AddFunc("*/10 * * * * *", PushTaskScan)
	s.c.Start()
}

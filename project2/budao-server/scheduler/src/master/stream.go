package master

import (
	"context"

	"github.com/sumaig/glog"
	pb "proto/scheduler"
)

//stream 对象是对grpc stream的一层封装
type Stream struct {
	internalStream pb.MessageService_WorkerHeartBeartServer
	Ctx            context.Context
	Cancel         func()
	Id             string
}

func newStream(internalStream pb.MessageService_WorkerHeartBeartServer) (s *Stream) {
	s = &Stream{
		internalStream: internalStream,
	}
	s.Ctx, s.Cancel = context.WithCancel(context.Background())
	return
}

func (s *Stream) Run() {
	for {
		select {
		case <-s.Ctx.Done():
			glog.Debug("stream %v canceled", s.Id)
			return
		}
	}
}

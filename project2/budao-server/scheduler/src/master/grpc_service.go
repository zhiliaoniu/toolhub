package master

import (
	"context"
	pb "proto/scheduler"
)

func (m *Master) GetMembers(ctx context.Context, req *pb.GetMembersRequest) (resp *pb.GetMembersResponse, err error) {
	return
}
func (m *Master) WorkerHeartBeart(stream pb.MessageService_WorkerHeartBeartServer) (err error) {
	s := newStream(stream)
	//firstPacket := s.Recv()
	//check first packet
	//bind to coordinator
	<-s.Ctx.Done()
	return
}

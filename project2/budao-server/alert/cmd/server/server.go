package main

import (
	"context"
	"fmt"
	"net/http"
	pb "rpc/budao"
)

const corpId = "wwd84057ff109fac91"
const corpSecret = "E0CwOI-CGe8GgYuAaVbaxz68qS2Vdr__vLrEpNQz95s"
const agentId = 1000015

type Server struct{}

func (s *Server) Alert(ctx context.Context, req *pb.AlertRequest) (res *pb.AlertResponse, err error) {

	fmt.Printf("req: %v\n", req)
	toUser := req.Touser
	content := req.Content
	if toUser == "" || content == "" {
		return &pb.AlertResponse{
			Code: 1,
			Msg:  "user/content empty",
		}, nil
	}

	for i := 0; i < 3; i++ {
		if !send(agentId, toUser, content) {
			refreshToken(corpId, corpSecret)
		} else {
			break
		}
	}

	return &pb.AlertResponse{
		Code: 0,
		Msg:  content,
	}, nil
}

// TODO undefined func
func (s *Server) AlertGroup(ctx context.Context, req *pb.AlertGroupRequest) (res *pb.AlertResponse, err error) {
	return &pb.AlertResponse{
		Code: 0,
		Msg:  "AlertGroup",
	}, nil
}

func main() {
	server := &Server{} // implements Haberdasher interface
	twirpHandler := pb.NewAlertWeChatWorkServer(server, nil)

	initToken()
	refreshToken(corpId, corpSecret)
	http.ListenAndServe(":8124", twirpHandler)

}

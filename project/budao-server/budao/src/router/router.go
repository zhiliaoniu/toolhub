package router

import (
	"net/http"

	"base"
	"service/commentserver"
	"service/configserver"
	"service/likeserver"
	"service/miscserver"
	"service/parseurlserver"
	"service/questionserver"
	"service/timelineserver"
	"service/topicserver"
	"service/transfer"
	"service/userserver"
	pb "twirprpc"
)

var (
	Mux *http.ServeMux // With read-write lock
)

// Registermulroutes function register multiple route pattern.
func Registermulroutes() {
	hook := base.NewStatsdServerHooks(base.GetLogStater())

	//timeline server
	timelineServer := timelineserver.GetServer()
	timelineHandler := pb.NewTimeLineServiceServer(timelineServer, hook)

	//transfer server
	transfer := transfer.GetTransfer()
	transferHandler := pb.NewTransferServer(transfer, hook)

	parseurlServer := parseurlserver.GetServer()
	parseurlHandler := pb.NewParseURLServiceServer(parseurlServer, hook)

	//common.WG.Add(1)
	commentServer := commentserver.GetServer()
	commentHandler := pb.NewCommentServiceServer(commentServer, hook)

	likeServer := likeserver.GetServer()
	likeHandler := pb.NewLikeServiceServer(likeServer, hook)

	miscServer := &miscserver.Server{}
	miscHandler := pb.NewMiscServiceServer(miscServer, hook)

	questionServer := questionserver.GetServer()
	questionHandle := pb.NewQuestionServiceServer(questionServer, hook)

	userServer := userserver.GetServer()
	userHandler := pb.NewUserServiceServer(userServer, hook)

	configServer := configserver.GetServer()
	configHandler := pb.NewConfigServiceServer(configServer, hook)

	topicServer := topicserver.GetServer()
	topicHandler := pb.NewTopicServiceServer(topicServer, hook)

	//glog.Debug("last wait")
	//common.WG.Wait()
	//glog.Debug("end wait")
	//// regularly load video's static information
	//timelineserver.TimerLoadVideoFullInfo()

	Mux = http.DefaultServeMux
	Mux.Handle(pb.TimeLineServicePathPrefix, timelineHandler)
	Mux.Handle(pb.TransferPathPrefix, transferHandler)
	Mux.Handle(pb.ParseURLServicePathPrefix, parseurlHandler)
	Mux.Handle(pb.CommentServicePathPrefix, commentHandler)
	Mux.Handle(pb.UserServicePathPrefix, userHandler)
	Mux.Handle(pb.LikeServicePathPrefix, likeHandler)
	Mux.Handle(pb.MiscServicePathPrefix, miscHandler)
	Mux.Handle(pb.QuestionServicePathPrefix, questionHandle)
	Mux.Handle(pb.ConfigServicePathPrefix, configHandler)
	Mux.Handle(pb.TopicServicePathPrefix, topicHandler)
}

// GetMux function get global Mux variables.
func GetMux() *http.ServeMux {
	return Mux
}

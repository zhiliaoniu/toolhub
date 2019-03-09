package router

import (
	"net/http"

	"base"
	"service/recommendserver"
	pb "twirprpc"
)

var (
	Mux *http.ServeMux // With read-write lock
)

// Registermulroutes function register multiple route pattern.
func Registermulroutes() {
	hook := base.NewStatsdServerHooks(base.GetLogStater())

	//recommend server
	recommendServer := recommendserver.GetServer()
	recommendHandler := pb.NewRecommendServiceServer(recommendServer, hook)

	Mux = http.DefaultServeMux
	Mux.Handle(pb.RecommendServicePathPrefix, recommendHandler)
}

// GetMux function get global Mux variables.
func GetMux() *http.ServeMux {
	return Mux
}

package router

import (
	"base"
	"net/http"
	"service/inputserver"
	pb "twirprpc"
)

var (
	Mux *http.ServeMux // With read-write lock
)

// Registermulroutes function register multiple route pattern.
func Registermulroutes() {
	hook := base.NewStatsdServerHooks(base.GetLogStater())

	//timeline server
	inputServer := inputserver.GetServer()
	inputHandler := pb.NewSpiderInputServiceServer(inputServer, hook)

	Mux = http.DefaultServeMux
	Mux.Handle(pb.SpiderInputServicePathPrefix, inputHandler)
}

// GetMux function get global Mux variables.
func GetMux() *http.ServeMux {
	return Mux
}

package router

import (
	"net/http"

	"base"
	"service/shareserver"
	pb "twirprpc"
)

var (
	Mux *http.ServeMux // With read-write lock
)

// Registermulroutes function register multiple route pattern.
func Registermulroutes() {
	hook := base.NewStatsdServerHooks(base.GetLogStater())

	shareServer := shareserver.GetServer()
	shareHandler := pb.NewShareServiceServer(shareServer, hook)

	Mux = http.DefaultServeMux
	Mux.Handle(pb.ShareServicePathPrefix, shareHandler)
}

// GetMux function get global Mux variables.
func GetMux() *http.ServeMux {
	return Mux
}

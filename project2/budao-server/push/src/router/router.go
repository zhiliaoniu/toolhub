package router

import (
	"net/http"
	"service/pushserver"
)

var (
	Mux *http.ServeMux // With read-write lock
)

// Registermulroutes function register multiple route pattern.
func Registermulroutes() {
	//hook := base.NewStatsdServerHooks(base.GetLogStater())

	_ = pushserver.GetServer()
	//pushHandler := pb.NewReportServiceServer(pushServer, hook)

	Mux = http.DefaultServeMux
	//Mux.Handle(pb.ReportServicePathPrefix, pushHandler)
}

// GetMux function get global Mux variables.
func GetMux() *http.ServeMux {
	return Mux
}

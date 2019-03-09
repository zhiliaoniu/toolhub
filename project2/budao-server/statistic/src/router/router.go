package router

import (
	"net/http"

	"base"
	"service/reportserver"
	pb "twirprpc"
)

var (
	Mux *http.ServeMux // With read-write lock
)

// Registermulroutes function register multiple route pattern.
func Registermulroutes() {
	hook := base.NewStatsdServerHooks(base.GetLogStater())

	reportServer := reportserver.GetServer()
	reportHandler := pb.NewReportServiceServer(reportServer, hook)

	Mux = http.DefaultServeMux
	Mux.Handle(pb.ReportServicePathPrefix, reportHandler)
}

// GetMux function get global Mux variables.
func GetMux() *http.ServeMux {
	return Mux
}

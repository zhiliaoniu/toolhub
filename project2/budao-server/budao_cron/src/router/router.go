package router

import (
	"common"
	"service/commentserver"
	"service/configserver"
	"service/timelineserver"
	"service/topicserver"

	"github.com/sumaig/glog"
)

func Registermulroutes() {
	_ = timelineserver.GetServer()

	common.WG.Add(1)
	_ = commentserver.GetServer()

	_ = configserver.GetServer()

	_ = topicserver.GetServer()

	glog.Debug("last wait")
	common.WG.Wait()
	glog.Debug("end wait")
	// regularly load video's static information
	timelineserver.TimerLoadVideoFullInfo()
}

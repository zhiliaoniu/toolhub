package base

import (
	"encoding/json"
	"log"

	"github.com/sumaig/glog"
)

/*
	init file logger with json config.
	sonConfig like:
	{
	"filename":"logs/beego.log",
	"maxLines":10000,
	"maxsize":1024,
	"daily":true,
	"maxDays":15,
	"rotate":true,
	    "perm":"0600"
	}
*/
//logger conf
type LoggerConf struct {
	Debug    bool   `json:"debug"`
	FileName string `json:"filename"`
	MaxDays  int    `json:"maxdays"`
}

type LoggerInnerConf struct {
	FileName string `json:"filename"`
	MaxDays  int    `json:"maxdays"`
}

// InitLog function set log level and switch.
func InitLogger(loggerConf *LoggerConf) {
	loggerInnerConf := &LoggerConf{
		FileName: loggerConf.FileName,
		MaxDays:  loggerConf.MaxDays,
	}
	loggerConfJson, err := json.Marshal(loggerInnerConf)
	if err != nil {
		log.Printf("json marshal failed. loggerConf:%v, err:%v\n", loggerConf, err)
		panic(err)
	}

	glog.SetLogger(glog.AdapterFile, string(loggerConfJson))
	glog.SetLogger(glog.AdapterConsole)
	glog.EnableFuncCallDepth(true)
	glog.SetLogFuncCallDepth(3)

	//if loggerConf.Debug == true {
	if true {
		glog.SetLevel(glog.LevelDebug)
	} else {
		glog.SetLevel(glog.LevelError)
	}

	log.Println("Log start success.")
}

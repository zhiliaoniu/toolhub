package service

import (
	"flag"

	"github.com/sumaig/glog"
	"github.com/zieckey/goini"

	"base"
)

var ConfPath = flag.String("c", "./conf/app.build.conf", "config file path")

//an internal http service
//parse conf and manage logger
type ServiceContext struct {
	Conf *goini.INI
}

func (s *ServiceContext) initConf() (err error) {
	if !flag.Parsed() {
		flag.Parse()
	}
	s.Conf, err = goini.LoadInheritedINI(*ConfPath)
	glog.Debug("iniConf:%v", s.Conf)
	return err
}

func (s *ServiceContext) initLogger() (err error) {
	debug, _ := s.Conf.SectionGetBool("log", "debug")
	fileName, _ := s.Conf.SectionGet("log", "file_name")
	maxDays, _ := s.Conf.SectionGetInt("log", "max_days")
	loggerConf := &base.LoggerConf{
		Debug:    debug,
		FileName: fileName,
		MaxDays:  maxDays,
	}
	base.InitLogger(loggerConf)
	return
}

func (s *ServiceContext) Init() (err error) {
	err = s.initConf()
	if err != nil {
		return err
	}

	err = s.initLogger()
	return
}

func (s *ServiceContext) Run() {
	/*
		for {
			time.Sleep(1000)
		}
	*/
	return
}

func (s *ServiceContext) Stop() {
}

func New() *ServiceContext {
	c := ServiceContext{}
	return &c
}

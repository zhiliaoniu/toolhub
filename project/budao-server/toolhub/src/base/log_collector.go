package base

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sumaig/glog"
)

const (
	AgentAddr string = "http://127.0.0.1:10039/metrics_api"
	BodyType  string = "application/json;charset=utf-8"
)

var CounterChan chan *Counter
var DefModelChan chan *DefModel
var HistogramChan chan *Histogram
var GlobalLogCollector *LogCollector

//call by main
func InitLogCollector(conf *LogCollectorConf) {
	CounterChan = make(chan *Counter, 2000)
	DefModelChan = make(chan *DefModel, 2000)
	HistogramChan = make(chan *Histogram, 2000)

	GlobalLogCollector = &LogCollector{}
	GlobalLogCollector.Init(conf)

	go GlobalLogCollector.LogCollectCronTask()
}

var GlobalLogCollectorConf *LogCollectorConf

//log collector conf
type LogCollectorConf struct {
	IsDisable    bool          `json:"isDisable"`
	AppName      string        `json:"appName"`
	ReportFormat *ReportFormat `json:"reportFormat"`
	MaxInterMSec int           `json:"maxInterMSec"`
	MaxRemainNum int           `json:"maxRemainNum"`
}

type ReportFormat struct {
	AppVersion string   `json:"appVersion"`
	Scale      []uint32 `json:"scale"`
}

//增加chan，把日志先发到chan中，由专门的协程将统计数据发往统计平台
//专门的协程定时或者定量的延时转发，并做初步计算，把相同数据合为一条
type LogCollector struct {
	maxInterMSec     int
	maxRemainNum     int
	currentRemainNum int
	statReport       *StatReport
	cronTaskDisable  bool
}

func (s *LogCollector) Init(conf *LogCollectorConf) {
	GlobalLogCollectorConf = conf
	//TODO use isDisable
	s.maxInterMSec = conf.MaxInterMSec
	s.maxRemainNum = conf.MaxRemainNum
	s.currentRemainNum = 0
	s.cronTaskDisable = false
	appVer := conf.ReportFormat.AppVersion
	appName := conf.AppName
	s.statReport = &StatReport{
		AppName:       appName,
		AppVer:        appVer,
		ServiceName:   appName,
		Step:          60,
		Ver:           "0.1",
		Skip1stPeriod: "false",
		DefModel:      make([]*DefModel, 0, 20000),
		Counter:       make([]*Counter, 0, 40000),
		Histogram:     make([]*Histogram, 0, 20000),
	}
}

func (s *LogCollector) Close() {
	s.cronTaskDisable = true
}

func (s *LogCollector) LogCollectCronTask() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(s.maxInterMSec))
	for {
		if s.cronTaskDisable {
			glog.Debug("exit config server crontab")
			break
		}

		select {
		case defModel := <-DefModelChan:
			s.statReport.DefModel = append(s.statReport.DefModel, defModel)
			s.currentRemainNum++
		case counter := <-CounterChan:
			s.statReport.Counter = append(s.statReport.Counter, counter)
			s.currentRemainNum++
		case histogram := <-HistogramChan:
			s.statReport.Histogram = append(s.statReport.Histogram, histogram)
			s.currentRemainNum++

		case <-ticker.C:
			s.PostStatLogToPlatform()
		case <-time.After(time.Millisecond * time.Duration(s.maxInterMSec/2)):
			if s.currentRemainNum > s.maxRemainNum {
				s.PostStatLogToPlatform()
			}
		}
	}
}

func (s *LogCollector) PostStatLogToPlatform() {
	if s.currentRemainNum == 0 {
		return
	}
	defer func() {
		s.statReport.DefModel = s.statReport.DefModel[:0]
		s.statReport.Counter = s.statReport.Counter[:0]
		s.statReport.Histogram = s.statReport.Histogram[:0]
		s.currentRemainNum = 0
	}()

	mJson, err := json.Marshal(s.statReport)
	if err != nil {
		glog.Error("marshal failed. err:%v", err)
		return
	}
	glog.Debug("post stat log num:%d", s.currentRemainNum)
	go func() {
		//glog.Debug("json:%s", mJson)
		req := bytes.NewBuffer(mJson)
		resp, err := http.Post(AgentAddr, BodyType, req)
		if err != nil {
			glog.Error("report stat failed. err:%v", err)
			return
		}
		//glog.Debug("resp:%+v", resp)
		//注意：resp如果不读取并关闭body的话，端口就不会重用，最终导致端口用光。
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
}

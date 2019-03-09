package base

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/sumaig/glog"
	"github.com/twitchtv/twirp"
	"github.com/twitchtv/twirp/hooks/statsd"
)

func NewStatsdServerHooks(stats statsd.Statter) *twirp.ServerHooks {
	hook := statsd.NewStatsdServerHooks(GetLogStater())
	hook.Error = func(ctx context.Context, err twirp.Error) context.Context {
		glog.Error("req failed. ctx:%v err:%v", ctx, err)
		return ctx
	}
	return hook
}

//TODO send log to stats dir
type LogStater struct {
}

type StatReport struct {
	AppName       string       `json:"app_name"`
	AppVer        string       `json:"app_ver"`
	ServiceName   string       `json:"service_name"`
	Step          int          `json:"stp"`
	Ver           string       `json:"ver"`
	Skip1stPeriod string       `json:"skip_1st_period"`
	DefModel      []*DefModel  `json:"defmodel,omitempty"`
	Counter       []*Counter   `json:"counter,omitempty"`
	Histogram     []*Histogram `json:"histo,omitempty"`
}

type DefModel struct {
	Topic     string   `json:"topic"`
	Uri       string   `json:"uri"`
	UriTag    string   `json:"uri_tag"`
	Duration  int64    `json:"duration"`
	Code      string   `json:"code"`
	IsSuccess string   `json:"isSuccess"`
	Scale     []uint32 `json:"scale"`
}

type Counter struct {
	Topic string `json:"topic"`
	Uri   string `json:"uri"`
	Val   int64  `json:"val"`
}

type Histogram struct {
	Topic string   `json:"topic"`
	Uri   string   `json:"uri"`
	Val   int64    `json:"val"`
	Scale []uint32 `json:"scale"`
}

func GetLogStater() LogStater {
	logStater := LogStater{}
	return logStater
}

func (ls LogStater) Inc(metric string, val int64, rate float32) error {
	glog.Debug("incr %s: %d @ %f", metric, val, rate)
	counter := &Counter{"req_resp_counter", metric, val}
	CounterChan <- counter
	return nil
}

func (ls LogStater) IncInternalError(counterName, metric string, val int64, rate float32) error {
	glog.Debug("incr %s err %s: %d @ %f", counterName, metric, val, rate)
	counter := &Counter{counterName, metric, val}
	CounterChan <- counter
	return nil
}

func (ls LogStater) TimingDuration(metric string, val time.Duration, rate float32) error {
	glog.Debug("time %s: %v @ %f", metric, val, rate)
	arr := strings.Split(metric, ".")
	if len(arr) == 4 {
		len := len(arr)
		code := arr[len-1]
		codeNum, _ := strconv.Atoi(code)
		isSuccess := "y"
		if codeNum > 300 {
			isSuccess = "n"
		}
		defModel := &DefModel{"", metric, "s", val.Nanoseconds() / 1000000, code, isSuccess, GlobalLogCollectorConf.ReportFormat.Scale}
		DefModelChan <- defModel
	} else if len(arr) == 3 {
		histogram := &Histogram{"histogram", metric, val.Nanoseconds() / 1000000, GlobalLogCollectorConf.ReportFormat.Scale}
		HistogramChan <- histogram
	}
	return nil
}

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"base"
	"cmd/spiderInput/router"
	"common"
	"db"
	"flag"
	"github.com/sumaig/glog"
)

type Skel struct {
	listener net.Listener
}

func (s *Skel) Start() (err error) {
	//init cfg
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version == true {
		fmt.Println(common.VERSION)
		os.Exit(0)
	}

	common.ParseConfig(*cfg)
	base.InitLogger(common.GetConfig().LoggerConf)
	runtime.GOMAXPROCS(runtime.NumCPU())

	db.InitMysqlClient(common.GetConfig().DB.MySQL)

	base.InitLogCollector(common.GetConfig().LogCollectorConf)
	//start server
	router.Registermulroutes()
	srv := &http.Server{Handler: router.GetMux()}
	s.listener, err = net.Listen("tcp", common.GetConfig().HTTPAddr)
	if err != nil {
		panic(err)
	}
	go func() {
		err = srv.Serve(s.listener)
		if err != nil {
			fmt.Fprintf(os.Stderr, "server exit. err:%v", err)
		}
	}()

	return
}

func (s *Skel) Stop() {
	s.listener.Close()
}

func (s *Skel) WaitSignal(wg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	go func() {
		defer close(c)
		for {
			sig := <-c
			if sig == syscall.SIGHUP || sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				glog.Debug("receiv SIGQUIT,graceful stop")
				s.Stop()
				glog.Debug("graceful stop ok ")
				wg.Done()
				break
			}
			glog.Debug("receiv ignore signal %v", s)
		}
	}()
}

func (s *Skel) Wait() {
	var wg sync.WaitGroup
	wg.Add(1)
	go s.WaitSignal(&wg)
	wg.Wait()
}

func main() {
	skel := &Skel{}
	if err := skel.Start(); err != nil {
		panic(err)
	}
	skel.Wait()
}

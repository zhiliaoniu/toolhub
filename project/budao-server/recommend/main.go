package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	_ "runtime/pprof"
	"sync"
	"syscall"

	"base"
	"common"
	"db"
	"router"

	"github.com/sumaig/glog"
)

// Skel listen
type Skel struct {
	listener net.Listener
}

// Start basic service
func (s *Skel) Start() (err error) {
	//init cfg
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version == true {
		fmt.Println("3.0")
		os.Exit(0)
	}

	common.ParseConfig(*cfg)
	base.InitLogger(common.GetConfig().LoggerConf)

	runtime.GOMAXPROCS(runtime.NumCPU())

	//init redis client
	db.InitRedisConnPool(common.GetConfig().DB.Redis)

	//init log collector
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
		glog.Debug("recommend-server exit success")
	}()

	return
}

// Stop close listen and db
func (s *Skel) Stop() {
	s.listener.Close()
}

// WaitSignal handle signal
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

// Wait defining the signal processing mechanism
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

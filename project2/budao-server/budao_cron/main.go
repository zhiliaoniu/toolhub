package main

import (
	"flag"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	_ "runtime/pprof"
	"sync"
	"syscall"

	"github.com/sumaig/glog"

	"base"
	"common"
	"db"
	"router"
)

// Skel listen
type Skel struct {
}

// Start basic service
func (s *Skel) Start() (err error) {
	//init cfg
	cfg := flag.String("c", "cfg.json", "configuration file")
	flag.Parse()

	common.ParseConfig(*cfg)
	base.InitLogger(common.GetConfig().LoggerConf)
	runtime.GOMAXPROCS(runtime.NumCPU())

	//init db client
	db.InitMysqlClient(common.GetConfig().DB.MySQL)
	db.InitRedisConnPool(common.GetConfig().DB.Redis)

	//init log collector
	base.InitLogCollector(common.GetConfig().LogCollectorConf)

	//start server
	router.Registermulroutes()

	return
}

// Stop close listen and db
func (s *Skel) Stop() {
	//db.QueryDB.Close()
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
				glog.Debug("graceful stop ok")
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

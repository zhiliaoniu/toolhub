package skel

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"service"

	"github.com/sumaig/glog"
)

type Skel struct {
	Service *service.ServiceContext
	modules map[string]Module
}

func New() (sk *Skel) {
	sk = &Skel{}
	sk.modules = map[string]Module{}
	return
}

func (sk *Skel) Start() (err error) {
	s := service.New()
	sk.Service = s
	err = s.Init()
	if err != nil {
		panic(err.Error())
	}
	for name, m := range sk.modules {
		if err = m.Init(s); err != nil {
			return errors.New(name + " init fail " + err.Error())
		}
	}
	for name, m := range sk.modules {
		if err = m.Start(); err != nil {
			return errors.New(name + " start fail " + err.Error())
		}
	}
	go s.Run()
	return
}

func (sk *Skel) WaitSignal(wg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT)
	go func() {
		defer close(c)
		for {
			sig := <-c
			if sig == syscall.SIGHUP || sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				glog.Debug("receiv SIGQUIT,graceful stop")
				sk.Service.Stop()
				glog.Debug("graceful stop ok ")
				wg.Done()
				break
			}
			glog.Debug("receiv ignore signal %v", sig)
		}
	}()
}

func (sk *Skel) Wait() {
	var wg sync.WaitGroup
	wg.Add(1)
	go sk.WaitSignal(&wg)
	wg.Wait()
}

package master

import (
	"sync"
	"time"

	"master/task"
)

type Cluster struct {
	sync.RWMutex
	isRunning   bool
	wg          sync.WaitGroup
	master      *Master
	taskManager task.Manager
	loader      task.Loader
	co          *Coordinator
}

func newScheudle(m *Master) (c *Cluster) {
	c = &Cluster{
		master: m,
		co:     newCoordinator(c),
	}

	return
}

func (s *Cluster) Init() (err error) {
	s.taskManager = newPbTaskMansger(s.master)
	s.loader = newEtcdLoader(s.master)
	return
}

func (s *Cluster) Start() (err error) {
	s.Lock()
	defer s.Unlock()
	if s.isRunning {
		return
	}
	err = s.LoadTaskInfos()
	if err != nil {
		return
	}
	s.wg.Add(1)
	go s.Run()
	//run coodinator
	s.wg.Add(1)
	go s.co.Run(s.wg)
	s.isRunning = true

	return
}

func (s *Cluster) LoadTaskInfos() (err error) {
	//TODO get data from etcd
	//data = get data from etcd
	//history = get data from etcd
	//s.taskManager.LoadData(loader)
	return
}

func (s *Cluster) Run() {
	defer s.wg.Done()
	for {
		if !s.isRunning {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func (s *Cluster) Stop() {
	s.Lock()
	defer s.Unlock()
	s.isRunning = false
	s.co.Stop()
	s.wg.Wait()
}

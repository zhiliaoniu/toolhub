package master

import (
	"sync"

	"github.com/juju/errors"
)

var (
	errStreamExist = errors.New("stream exist")
)

type Coordinator struct {
	cluster *Cluster
	stopped bool

	streamsLock sync.Mutex
	streams     map[string]*Stream
}

func newCoordinator(c *Cluster) (co *Coordinator) {
	co = &Coordinator{
		cluster: c,
	}
	return
}

func (co *Coordinator) Run(wg sync.WaitGroup) {
	defer wg.Done()
	co.stopped = false
	for {
		if co.stopped {
			break
		}
		//workers := taskManager.NextWorkers()
		//if len(workers) == 0 {
		//	continue
		//}
		//co.Dispatch(workers)
	}
}

// 选择一个stream，把worker分发出去
//func (co *Coordinator) Dispatch(workers []*task.Worker) {
//}

func (co *Coordinator) Stop() {
	co.stopped = true
}

func (co *Coordinator) BindStream(s *Stream) (err error) {
	id := s.Id
	co.streamsLock.Lock()
	defer co.streamsLock.Unlock()
	if _, ok := co.streams[id]; ok {
		return errStreamExist
	}
	co.streams[id] = s
	s.Run()
	return
}

func (co *Coordinator) RemoveStream(id string) {
	co.streamsLock.Lock()
	defer co.streamsLock.Unlock()
	if old, ok := co.streams[id]; ok {
		old.Cancel()
		delete(co.streams, id)
	}
}

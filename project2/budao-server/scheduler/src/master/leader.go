package master

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"etcdop"
	pb "proto/scheduler"

	"github.com/coreos/etcd/clientv3"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"github.com/sumaig/glog"
)

var (
	errNoLeader = errors.New("no leader")
)

//实现了leader选举及leader loop的逻辑

func (m *Master) IsLeader() bool {
	return atomic.LoadInt64(&m.isLeader) == 1
}

func (m *Master) enableLeader(b bool) {
	value := int64(0)
	if b {
		value = 1
	}

	atomic.StoreInt64(&m.isLeader, value)
}

//path to write leader value
func (m *Master) LeaderKey() string {
	return m.cfg.rootPath + "/leader"
}

//when this node become leader
//write this value to LeaderKey
//if err,must panic
func (m *Master) LeaderValue() string {
	leader := &pb.Member{
		Name:       m.Name(),
		MemberId:   m.ID(),
		ClientUrls: strings.Split(m.cfg.Etcd.AdvertiseClientUrls, ","),
		PeerUrls:   strings.Split(m.cfg.Etcd.AdvertisePeerUrls, ","),
	}
	v, err := proto.Marshal(leader)
	if err != nil {
		panic(err)
	}
	return string(v)
}

func (m *Master) StartLeaderLoop() (err error) {
	m.leaderLoopWg.Add(1)
	m.leaderLoopCtx, m.leaderLoopCancel = context.WithCancel(context.Background())
	go m.LeaderLoop()
	glog.Debug("StartLeaderLoop end")
	return
}

func (m *Master) StopLeaderLoop() (err error) {
	glog.Debug("stop leader loop begin")
	m.leaderLoopCancel()
	m.leaderLoopWg.Wait()
	glog.Debug("stop leader loop ok")
	return
}

func (m *Master) LeaderLoop() {
	defer m.leaderLoopWg.Done()
	for {
		glog.Debug("LeaderLoop for begin. m.ID:%d", m.ID())
		if m.stopped {
			glog.Debug("master stopped, quit")
			return
		}
		leader, err := m.GetLeader()
		if err != nil && err != errNoLeader {
			time.Sleep(200 * time.Millisecond)
			glog.Error("get leader failed, err %v", err)
			continue
		}
		glog.Debug("leader is: %+v", leader)
		if leader != nil {
			if m.isSameLeader(leader) {
				// oh, we are already leader, we may meet something wrong
				// in previous campaignLeader. we can delete and campaign again.
				glog.Warn("leader is still %s, delete and campaign again", leader)
				if err = m.deleteLeaderKey(); err != nil {
					glog.Error("delete leader key err %s", err)
					time.Sleep(200 * time.Millisecond)
					continue
				}
			} else {
				glog.Debug("leader is %s, watch it", leader)
				m.watchLeader()
				glog.Debug("leader changed, try to campaign leader")
			}
		}

		etcdLeader := m.etcd.Server.Lead()
		if etcdLeader != m.ID() {
			glog.Error("id not equal. m.etcd.Server.Lead():%d, m.ID():%d", etcdLeader, m.ID())
			time.Sleep(200 * time.Millisecond)
			continue
		}
		glog.Debug("m.etcd.Server.Lead():%d, m.ID():%d", etcdLeader, m.ID())

		if err = m.campaignLeader(); err != nil {
			glog.Error("campaign leader err %s", errors.ErrorStack(err))
		}
		glog.Debug("after campaignLeader")
	}
}

func (m *Master) campaignLeader() (err error) {
	glog.Debug("begin to campaign leader %s", m.Name())

	//1.create lessor
	lessor := clientv3.NewLease(m.client)
	defer lessor.Close()

	start := time.Now()
	ctx, cancel := context.WithTimeout(m.client.Ctx(), requestTimeout)
	leaseResp, err := lessor.Grant(ctx, m.cfg.LeaderLease)
	cancel()

	cost := time.Since(start)
	if cost > slowRequestTime {
		glog.Warn("lessor grants too slow, cost %s", cost)
	}
	glog.Debug("lessor grant cost time:%v", cost)

	if err != nil {
		glog.Error("err1:%v", err)
		return errors.Trace(err)
	}

	//2.put leader key
	leaderKey := m.LeaderKey()
	//ID := int64(leaseResp.ID)
	//IDStr := strconv.FormatInt(ID, 10)
	//IDHex := hex.EncodeToString([]byte(IDStr))
	//glog.Debug("idhex:%s", IDHex)
	glog.Debug("leaderKey:%s, leaderValue:%+v, leaseResp.ID:%v", leaderKey, m.LeaderValue(), leaseResp.ID)
	// The leader key must not exist, so the CreateRevision is 0.
	resp, err := m.txn().
		If(clientv3.Compare(clientv3.CreateRevision(leaderKey), "=", 0)).
		Then(clientv3.OpPut(leaderKey, m.LeaderValue(), clientv3.WithLease(clientv3.LeaseID(leaseResp.ID)))).
		Commit()
	if err != nil {
		glog.Error("err2:%v", err)
		return errors.Trace(err)
	}
	if !resp.Succeeded {
		glog.Error("err3:%v", err)
		return errors.New("campaign leader failed, other server may campaign ok")
	}
	glog.Debug("put leaderkey resp:%v", resp)

	//3.keepalive leader key.
	ctx, cancel = context.WithCancel(m.leaderLoopCtx)
	defer cancel()

	ch, err := lessor.KeepAlive(ctx, clientv3.LeaseID(leaseResp.ID))
	if err != nil {
		glog.Error("err4:%v", err)
		return errors.Trace(err)
	}
	glog.Debug("campaign leader ok %s", m.Name())

	/*
		//4.begin leader biz
		m.enableLeader(true)
		defer m.enableLeader(false)

		glog.Debug("leader %s is ready to serve", m.Name())

		if err = m.cluster.Start(); err != nil {
			glog.Error("%v start cluster fail,give up leader", m.Name())
			return
		}

		glog.Debug("leader %s is start to schedule cluster", m.Name())
		defer m.cluster.Stop()
	*/

	//5.check keepalive
	tsTicker := time.NewTicker(leaderTimerStep)
	defer tsTicker.Stop()
	for {
		glog.Debug("begin for. ch.len:%d", len(ch))
		//if len(ch) == 1 {
		//	select {
		//	case r, ok := <-ch:
		//		if !ok {
		//			glog.Error("keep alive channel is closed")
		//			return nil
		//		} else {
		//			glog.Debug("ch return ok. r:%v", r)
		//		}
		//	}
		//}

		select {
		case r, ok := <-ch:
			if !ok {
				glog.Error("keep alive channel is closed")
				return nil
			} else {
				glog.Debug("ch return ok. r:%v", r)
			}
		case <-tsTicker.C:
			glog.Debug("on timer, now do nothing, name:%s", m.Name())
		case <-m.resignCh:
			glog.Debug("%s resigns leadership", m.Name())
			return nil
		case <-ctx.Done():
			glog.Debug("ctx done. %s resigns leadership", m.Name())
			return errors.New("server closed")
		}
	}
	glog.Debug("end compain leader")
	return
}

func (m *Master) leaderCmp() clientv3.Cmp {
	return clientv3.Compare(clientv3.Value(m.LeaderKey()), "=", m.LeaderValue())
}

func (m *Master) leaderTxn(cs ...clientv3.Cmp) clientv3.Txn {
	return m.txn().If(append(cs, m.leaderCmp())...)
}

func (m *Master) deleteLeaderKey() (err error) {
	leaderKey := m.LeaderKey()
	resp, err := m.leaderTxn().Then(clientv3.OpDelete(leaderKey)).Commit()
	if err != nil {
		glog.Error("delete leader key failed. err:%v", err)
		return errors.Trace(err)
	}
	if !resp.Succeeded {
		glog.Error("delete leader key failed2. err:%v", err)
		return errors.New("resign leader failed, we are not leader already")
	}
	glog.Debug("delete leader key success. leaderkey:%s", leaderKey)
	return
}

func (m *Master) watchLeader() {
	glog.Debug("begin watchLeader")
	watcher := clientv3.NewWatcher(m.client)
	defer watcher.Close()

	ctx, cancel := context.WithCancel(m.leaderLoopCtx)
	defer cancel()

	for {
		rch := watcher.Watch(ctx, m.LeaderKey())
		for wresp := range rch {
			if wresp.Canceled {
				glog.Debug("wresp.Canceled true")
				return
			}

			for _, ev := range wresp.Events {
				if ev.Type == clientv3.EventTypeDelete {
					glog.Debug("leader is deleted")
					return
				}
			}
		}

		select {
		case <-ctx.Done():
			// server closed, return
			return
		default:
		}
	}
}

func (m *Master) isSameLeader(leader *pb.Member) bool {
	return leader.GetMemberId() == m.ID()
}

func (m *Master) GetLeader() (*pb.Member, error) {
	leader := &pb.Member{}
	ok, err := etcdop.GetProtoMsg(m.client, m.LeaderKey(), leader)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if !ok {
		return nil, errNoLeader
	}
	return leader, nil
}

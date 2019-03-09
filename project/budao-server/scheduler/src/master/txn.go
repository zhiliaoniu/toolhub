package master

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/juju/errors"
	"github.com/pingcap/pd/pkg/etcdutil"
	"github.com/sumaig/glog"
)

const (
	slowRequestTime = etcdutil.DefaultSlowRequestTime
	requestTimeout  = etcdutil.DefaultRequestTimeout
	leaderTimerStep = 1000 * time.Millisecond
)

type slowLogTxn struct {
	clientv3.Txn
	cancel context.CancelFunc
}

func newSlowLogTxn(client *clientv3.Client) clientv3.Txn {
	ctx, cancel := context.WithTimeout(client.Ctx(), requestTimeout)
	return &slowLogTxn{
		Txn:    client.Txn(ctx),
		cancel: cancel,
	}
}

func (t *slowLogTxn) If(cs ...clientv3.Cmp) clientv3.Txn {
	return &slowLogTxn{
		Txn:    t.Txn.If(cs...),
		cancel: t.cancel,
	}
}

func (t *slowLogTxn) Then(ops ...clientv3.Op) clientv3.Txn {
	return &slowLogTxn{
		Txn:    t.Txn.Then(ops...),
		cancel: t.cancel,
	}
}

// Commit implements Txn Commit interface.
func (t *slowLogTxn) Commit() (*clientv3.TxnResponse, error) {
	start := time.Now()
	resp, err := t.Txn.Commit()
	t.cancel()

	cost := time.Since(start)
	if cost > slowRequestTime {
		glog.Warn("txn runs too slow, resp: %v, err: %v, cost: %s", resp, err, cost)
	}
	//TODO add some stat here
	return resp, errors.Trace(err)
}

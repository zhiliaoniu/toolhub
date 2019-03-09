package etcdop

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
	"github.com/sumaig/glog"
)

var logName = "CloudSafeLine.Quantum.Etcdop"

const (
	KvRequestTimeout  = time.Second * 10
	KvSlowRequestTime = time.Second * 1
)

func kvGet(c *clientv3.Client, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(c.Ctx(), KvRequestTimeout)
	defer cancel()

	start := time.Now()
	resp, err := clientv3.NewKV(c).Get(ctx, key, opts...)
	if err != nil {
		glog.Error("load from etcd error: %v", err)
	}
	if cost := time.Since(start); cost > KvSlowRequestTime {
		glog.Warn("kv gets too slow: key %v cost %v err %v", key, cost, err)
	}

	return resp, errors.Trace(err)
}

// A helper function to get value with key from etcd.
// TODO: return the value revision for outer use.
func getValue(c *clientv3.Client, key string, opts ...clientv3.OpOption) ([]byte, error) {
	resp, err := kvGet(c, key, opts...)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if n := len(resp.Kvs); n == 0 {
		return nil, nil
	} else if n > 1 {
		return nil, errors.Errorf("invalid get value resp %v, must only one", resp.Kvs)
	}

	return resp.Kvs[0].Value, nil
}

// Return boolean to indicate whether the key exists or not.
// TODO: return the value revision for outer use.
func GetProtoMsg(c *clientv3.Client, key string, msg proto.Message, opts ...clientv3.OpOption) (bool, error) {
	value, err := getValue(c, key, opts...)
	if err != nil {
		return false, errors.Trace(err)
	}
	if value == nil {
		return false, nil
	}

	if err = proto.Unmarshal(value, msg); err != nil {
		return false, errors.Trace(err)
	}

	return true, nil
}

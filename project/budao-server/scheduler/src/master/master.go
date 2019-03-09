package master

import (
	"context"
	"crypto/tls"
	"strings"
	"sync"
	"time"

	"service"
	"skel"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/embed"
	"github.com/coreos/etcd/pkg/types"
	"github.com/juju/errors"
	"github.com/pingcap/pd/pkg/etcdutil"
	"github.com/sumaig/glog"
)

const (
	etcdTimeout      = time.Second * 3
	etcdStartTimeout = time.Second * 180
)

type Master struct {
	skel.ModuleBase
	//the master's conf
	cfg      *Config
	etcd     *embed.Etcd
	client   *clientv3.Client
	id       uint64 // etcd server id.
	resignCh chan struct{}

	stopped  bool
	isLeader int64

	leaderLoopCtx    context.Context
	leaderLoopCancel func()
	leaderLoopWg     sync.WaitGroup

	//业务集群的信息
	cluster *Cluster
}

func (m *Master) ID() uint64 {
	return m.id
}

func (m *Master) Name() string {
	return m.cfg.Etcd.Name
}

func New() (m *Master) {
	m = &Master{
		resignCh: make(chan struct{}),
	}
	m.cluster = newScheudle(m)
	return
}

func (m *Master) Init(s *service.ServiceContext) (err error) {
	confFilePath, ok := s.Conf.SectionGet("common", "master_config_path")
	if !ok {
		return errors.New("not find master conf path")
	}
	glog.Debug("Init common.master_config_path confFilePath:%s", confFilePath)
	m.cfg, err = NewConfigFromeFile(confFilePath)
	if err != nil {
		return
	}
	glog.Debug("Init conf:%+v", m.cfg)
	err = m.cluster.Init()

	return
}

func (m *Master) Start() (err error) {
	if err = m.StartEtcd(); err != nil {
		glog.Error("start etcd failed. err:%v", err)
		return
	}
	glog.Debug("to start own leader")
	if err = m.StartLeaderLoop(); err != nil {
		glog.Error("start leader loop failed. err:%v", err)
		return
	}

	return
}

func (m *Master) Stop() {
	m.stopped = true
	m.StopLeaderLoop()
	m.cluster.Stop()

	glog.Debug("master stop success")

	return
}

func (m *Master) StartEtcd() (err error) {
	//1.start etcd
	glog.Debug("start etcd begin")

	embedCfg, err := m.cfg.Etcd.GenEmbedEtcdConfig()
	if err != nil {
		return errors.Trace(err)
	}
	etcd, err := embed.StartEtcd(embedCfg)
	if err != nil {
		glog.Error("start etcd fail err %v", err)
		return errors.Trace(err)
	}
	glog.Debug("start etcd end. ok\n")

	//2.check cluster ID
	glog.Debug("check cluster id begin.")
	urlmap, err := types.NewURLsMap(embedCfg.InitialCluster)
	if err != nil {
		return errors.Trace(err)
	}
	glog.Debug("urlMap.len:%d", len(urlmap))
	var tlsConfig *tls.Config
	if err = etcdutil.CheckClusterID(etcd.Server.Cluster().ID(), urlmap, tlsConfig); err != nil {
		return errors.Trace(err)
	}
	glog.Debug("wait etcd cluster ready")

	select {
	// Wait etcd until it is ready to use
	case <-etcd.Server.ReadyNotify():
	case <-time.After(etcdStartTimeout):
		glog.Error("after %v seconds  waiting embed etcd to be ready,error", etcdStartTimeout)
		return errors.Errorf("after %v seconds  waiting embed etcd to be ready,error", etcdStartTimeout)
	}
	m.etcd = etcd
	glog.Debug("etcd cluster ready")

	err = m.getStartedEtcdInfo(embedCfg, tlsConfig)

	glog.Debug("start etcd end")

	return
}

func (m *Master) getStartedEtcdInfo(embedCfg *embed.Config, tlsConfig *tls.Config) (err error) {
	endpoints := []string{embedCfg.ACUrls[0].String()}
	glog.Debug("create etcd v3 client with endpoints %v", endpoints)

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: etcdTimeout,
		TLS:         tlsConfig,
	})
	if err != nil {
		return errors.Trace(err)
	}
	etcdServerID := uint64(m.etcd.Server.ID())
	glog.Debug("m.etcd.Server.ID:%d", etcdServerID)

	// update advertise peer urls.
	etcdMembers, err := etcdutil.ListEtcdMembers(client)
	if err != nil {
		return errors.Trace(err)
	}
	glog.Debug("etcdMembers:%v", etcdMembers)
	for _, memb := range etcdMembers.Members {
		if etcdServerID == memb.ID {
			etcdPeerURLs := strings.Join(memb.PeerURLs, ",")
			if m.cfg.Etcd.AdvertisePeerUrls != etcdPeerURLs {
				glog.Debug("update advertise peer urls from %s to %s", m.cfg.Etcd.AdvertisePeerUrls, etcdPeerURLs)
				m.cfg.Etcd.AdvertisePeerUrls = etcdPeerURLs
			}
		}
	}

	m.client = client
	m.id = etcdServerID
	glog.Debug("master.id:%d", m.id)
	return

}

func (m *Master) txn() clientv3.Txn {
	return newSlowLogTxn(m.client)
}

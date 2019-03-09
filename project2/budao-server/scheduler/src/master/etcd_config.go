package master

import (
	"fmt"
	"github.com/sumaig/glog"
	"net/url"
	"os"
	"strings"

	"github.com/coreos/etcd/embed"
	"github.com/juju/errors"
	"master/valueutil"
)

func newEtcdConfig() (c *EtcdConfig) {
	c = &EtcdConfig{}
	return
}

type EtcdConfig struct {
	//name of this etcd master
	Name                       string `json:"name"`
	DataDir                    string `json:"data-dir"`
	InitialCluster             string `json:"initial-cluster"`
	InitialClusterState        string `json:"initial-cluster-state"`
	DisableStrictReconfigCheck bool   `json:"disable-strict-reconfig-check"`
	TickMs                     int64  `json:"tick-ms"`
	ElectionMs                 int64  `json:"election-ms"`

	// QuotaBackendBytes Raise alarms when backend size exceeds the given quota. 0 means use the default quota.
	// the default size is 2GB, the maximum is 8GB.
	QuotaBackendBytes int64 `json:"quota-backend-bytes"`
	// AutoCompactionRetention for mvcc key value store in hour. 0 means disable auto compaction.
	// the default retention is 1 hour ,compactor.ModePeriodic
	AutoCompactionRetention string `json:"auto-compaction-retention"`

	//本member侧使用，用于监听etcd客户发送信息的地址。ip为全0代表监听本member侧所有接口
	ClientUrls string `json:"client-urls"`
	//本member侧使用，用于监听其他member发送信息的地址。ip为全0代表监听本member侧所有接口,比如ip写0.0.0.0
	PeerUrls string `json:"peer-urls"`
	//其他member使用，其他member通过该地址与本member交互信息。一定要保证从其他member能可访问该地址,比如一个内网地址
	AdvertiseClientUrls string `json:"advertise-client-urls"`
	//其他member使用，其他member通过该地址与本member交互信息。一定要保证从其他member能可访问该地址
	AdvertisePeerUrls string `json:"advertise-peer-urls"`
}

func HostNameOrEmpty() string {
	host, err := os.Hostname()
	if err != nil {
		return ""
	}
	return host
}

var (
	defaultName = HostNameOrEmpty()
)

const (
	defaultClientUrls          = "http://127.0.0.1:21987"
	defaultPeerUrls            = "http://127.0.0.1:21989"
	defaultInitialClusterState = embed.ClusterStateFlagNew

	// etcd use 100ms for heartbeat and 1s for election timeout.
	// We can enlarge both a little to reduce the network aggression.
	// now embed etcd use TickMs for heartbeat, we will update
	// after embed etcd decouples tick and heartbeat.
	defaultTickInterval = 500 //ms
	// embed etcd has a check that `5 * tick > election`
	defaultElectionInterval = 3000 //ms
)

func (c *EtcdConfig) InitDefaultValues() {
	valueutil.CheckSetString(&c.Name, defaultName)
	valueutil.CheckSetString(&c.DataDir, fmt.Sprintf("/data/budao-server/scheduler/default.%s", c.Name))
	valueutil.CheckSetString(&c.ClientUrls, defaultClientUrls)
	valueutil.CheckSetString(&c.AdvertiseClientUrls, c.ClientUrls)
	valueutil.CheckSetString(&c.PeerUrls, defaultPeerUrls)
	valueutil.CheckSetString(&c.AdvertisePeerUrls, c.PeerUrls)
	if len(c.InitialCluster) == 0 {
		// The advertise peer urls may be http://127.0.0.1:2380,http://127.0.0.1:2381
		// so the initial cluster is pd=http://127.0.0.1:2380,pd=http://127.0.0.1:2381
		items := strings.Split(c.AdvertisePeerUrls, ",")

		sep := ""
		for _, item := range items {
			c.InitialCluster += fmt.Sprintf("%s%s=%s", sep, c.Name, item)
			sep = ","
		}
	}
	valueutil.CheckSetString(&c.InitialClusterState, defaultInitialClusterState)
	valueutil.CheckSetInt64(&c.TickMs, defaultTickInterval)
	valueutil.CheckSetInt64(&c.ElectionMs, defaultElectionInterval)
	return
}

func (c *EtcdConfig) Check() error {
	if c.Name == "" {
		return errors.New("empty etcd node name")
	}
	return nil
}
func (c *EtcdConfig) GenEmbedEtcdConfig() (*embed.Config, error) {
	cfg := embed.NewConfig()
	cfg.Name = c.Name
	cfg.Dir = c.DataDir
	cfg.WalDir = ""
	cfg.InitialCluster = c.InitialCluster
	cfg.ClusterState = c.InitialClusterState
	cfg.StrictReconfigCheck = !c.DisableStrictReconfigCheck
	cfg.TickMs = uint(c.TickMs)
	cfg.ElectionMs = uint(c.ElectionMs)
	cfg.AutoCompactionRetention = c.AutoCompactionRetention
	cfg.QuotaBackendBytes = int64(c.QuotaBackendBytes)

	var err error

	cfg.LPUrls, err = ParseUrls(c.PeerUrls)
	if err != nil {
		return nil, errors.Trace(err)
	}

	cfg.APUrls, err = ParseUrls(c.AdvertisePeerUrls)
	if err != nil {
		return nil, errors.Trace(err)
	}

	cfg.LCUrls, err = ParseUrls(c.ClientUrls)
	if err != nil {
		return nil, errors.Trace(err)
	}

	cfg.ACUrls, err = ParseUrls(c.AdvertiseClientUrls)
	if err != nil {
		return nil, errors.Trace(err)
	}
	glog.Debug("generate embedEtcdCfg ok %+v", *cfg)

	return cfg, nil
}

// ParseUrls parse a string into multiple urls.
// Export for api.
func ParseUrls(s string) ([]url.URL, error) {
	items := strings.Split(s, ",")
	urls := make([]url.URL, 0, len(items))
	for _, item := range items {
		u, err := url.Parse(item)
		if err != nil {
			return nil, errors.Trace(err)
		}

		urls = append(urls, *u)
	}

	return urls, nil
}

package master

import (
	"encoding/json"
	"io/ioutil"

	"master/valueutil"

	_ "github.com/coreos/etcd/embed"
	"github.com/sumaig/glog"
)

const (
	defaultRootPath    = "/scheduler"
	defaultLeaderLease = int64(3)
)

type Config struct {
	Etcd *EtcdConfig `json:"etcd"`
	//root path of etcd
	rootPath    string `json:"root-path"`
	LeaderLease int64  `json:"lease"`
}

func NewConfigFromeFile(path string) (c *Config, err error) {
	c = &Config{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		glog.Debug("unmarshal failed %v", path)
		return
	}
	c.InitDefaultValues()
	err = c.Init()
	return
}

func (c *Config) InitDefaultValues() {
	valueutil.CheckSetString(&c.rootPath, defaultRootPath)
	valueutil.CheckSetInt64(&c.LeaderLease, defaultLeaderLease)
	if c.Etcd == nil {
		c.Etcd = newEtcdConfig()
	}
	c.Etcd.InitDefaultValues()
	return
}

func (c *Config) Init() (err error) {
	err = c.Etcd.Check()
	//c.embedEtcdCfg, err = c.Etcd.GenEmbedEtcdConfig()
	return
}

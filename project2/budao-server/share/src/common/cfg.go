package common

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"

	"base"
	"db"
)

// GlobalConfig structure.
type GlobalConfig struct {
	APPName          string                 `json:"appName"`
	HTTPAddr         string                 `json:"httpAddr"`
	DB               *DBConfig              `json:"db"`
	LoggerConf       *base.LoggerConf       `json:"loggerconf"`
	LogCollectorConf *base.LogCollectorConf `json:"logCollectorConf"`
	VideoBottomPage  *VideoBottomPageConf   `json:"videoBottomPageConf"`
	TopicPage        *TopicPageConf         `json:"TopicPageConf"`
	DownLoadURL      string                 `json:"downLoadURL"`
}

// DBConfig structure.
type DBConfig struct {
	Redis *db.RedisConfig `json:"redis"`
}

type VideoBottomPageConf struct {
	CommentNUM int64 `json:"commentNUM"`
}

type TopicPageConf struct {
	VideoNUM int64 `json:"videoNUM"`
}

var (
	config *GlobalConfig
	lock   = new(sync.RWMutex)
)

// GetConfig function get global configuration structure variables.
func GetConfig() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()

	return config
}

// ParseConfig parse configure parameters.
func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file.")
	}
	if !file.IsExist(cfg) {
		log.Fatalln("configuration file:", cfg, "isn`t exist.")
	}

	// read configure file.
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "failed. error:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("json unmarshal config file:", cfg, "failed. error:", err)
	}

	lock.Lock()
	defer lock.Unlock()
	config = &c
}

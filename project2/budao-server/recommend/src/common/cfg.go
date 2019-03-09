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
	RecommendConf    *RecommendConf         `json:"recommendConf"`
}

// DBConfig structure.
type DBConfig struct {
	Redis *db.RedisConfig `json:"redis"`
}

//-----------------------recommend conf------------------
type RecommendConf struct {
	CronTaskDisable     bool `json:"cronTaskDisable"`
	CronTaskInternalSec int  `json:"cronTaskInteranlSec"`
	MaxRetVideoNum      int  `json:"maxRetVideoNum"`
	RetPerTopicViewNum  int  `json:"retPerTopicViewNum"`
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

func InitConfig(cfg string) {
	var c GlobalConfig
	err := json.Unmarshal([]byte(cfg), &c)
	if err != nil {
		log.Fatalln("json unmarshal config file:", cfg, "failed. error:", err)
		return
	}
	config = &c
}

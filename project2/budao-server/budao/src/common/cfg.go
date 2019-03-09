package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/sumaig.bak/glog"
	"github.com/toolkits/file"

	"base"
	"db"
)

// GlobalConfig structure.
type GlobalConfig struct {
	APPName             string                 `json:"appName"`
	HTTPAddr            string                 `json:"httpAddr"`
	DB                  *DBConfig              `json:"db"`
	Audit               map[string]*AuditConf  `json:"auditconf"`
	LoggerConf          *base.LoggerConf       `json:"loggerconf"`
	LogCollectorConf    *base.LogCollectorConf `json:"logCollectorConf"`
	CommentConf         *CommentConf           `json:"commentConf"`
	TimelineConf        *TimelineConf          `json:"timelineConf"`
	IOSAuditConf        *IOSAuditConf          `json:"iosAuditConf"`
	RecommendClientConf *RecommendClientConf   `json:"recommendClientConf"`
	TopicConf           *TopicConf             `json:"topicConf"`
	TestChannelShow     string                 `json:"testChannelShow"`
	ParseTypeConf       string                 `json:"parseTypeConf"`
}

// DBConfig structure.
type DBConfig struct {
	MySQL map[string]*db.MySQLConfig `json:"mySQL"`
	Redis *db.RedisConfig            `json:"redis"`
}

//审核相关
type AuditConf struct {
	AuditAddr  string      `json:"auditAddr"` //审核视频的服务地址
	Appid      string      `json:"appid"`     //业务唯一标示
	Chid       string      `json:"chid"`      //业务通道id
	SecretId   string      `json:"secretId"`  //业务凭证，秘钥ID
	SecretKey  string      `json:"secretKey"` //业务凭证，秘钥Key
	PunishConf *PunishConf `json:"punishConf"`
	IsDisable  bool        `json:"isDisable"` //是否关闭审核
}

type PunishConf struct {
	ListenAddr string `json:"listenAddr"`
}

//-----------------------comment conf------------------
type CommentConf struct {
	CronTaskDisable     bool `json:"cronTaskDisable"`
	CronTaskInternalSec int  `json:"cronTaskInteranlSec"`
	ReturnHotCommentNum int  `json:"returnHotCommentNum"`
	ReturnCommentNum    int  `json:"returnCommentNum"`
	ReturnReplyNum      int  `json:"returnReplyNum"`
	ExplicitReplyNum    int  `json:"explicitReplyNum"`
}

//-----------------------timeline conf------------------
type TimelineConf struct {
	ShowQuestion        bool `json:"showQuestion"`
	CronTaskDisable     bool `json:"cronTaskDisable"`
	CronTaskInternalSec int  `json:"cronTaskInteranlSec"`
	MaxRetVideoNum      int  `json:"maxRetVideoNum"`
	RetPerTopicViewNum  int  `json:"retPerTopicViewNum"`
	ExposeVideoMaxLen   int  `json:"exposeVideoMaxLen"`
}

//-----------------------ios audit conf------------------
type IOSAuditConf struct {
	VidsFileName     string `json:"vidsFileName"`
	TopicIdsFileName string `json:"topicIdsFileName"`
	AppVersion       string `json:"appVersion"`
}

//-----------------------recommendClient conf------------------
type RecommendClientConf struct {
	ReqTimeoutMs int    `json:"reqTimeoutMs"`
	Addr         string `json:"addr"`
}

//-----------------------Topic Conf------------------
type TopicConf struct {
	ShareDisabled bool `json:"shareDisabled"`
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

func InitConf(confPath string) (err error) {
	if err = ParseConfig(confPath); err != nil {
		panic(err)
	}
	if err = ParseABTestConfig("/data/budao-server/scheduler/conf/budao/budao_abtest.json"); err != nil {
		return
	}
	return
}

// ParseConfig parse configure parameters.
func ParseConfig(cfg string) (err error) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file.")
		return
	}
	if !file.IsExist(cfg) {
		log.Fatalln("configuration file:", cfg, "isn`t exist.")
		return
	}

	// read configure file.
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "failed. error:", err)
		return
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("json unmarshal config file:", cfg, "failed. error:", err)
		return
	}

	lock.Lock()
	defer lock.Unlock()
	config = &c

	return
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

//--------------------------ab test conf----------------------------

var (
	abTestConf *ABTestConf
	ablock     = new(sync.RWMutex)
)

type QuestionOpenTest struct {
	Percent int `json:"percent"`
}

type ABTestConf struct {
	QuestionOpenTestMap map[string]*QuestionOpenTest `json:"questionOpenTest"`
}

func GetABTestConfig() *ABTestConf {
	ablock.RLock()
	defer ablock.RUnlock()

	return abTestConf
}

func ParseABTestConfig(confPath string) (err error) {
	if confPath == "" {
		glog.Error("use -c to specify configuration file.")
		log.Fatalln("use -c to specify configuration file.")
		err = errors.New("confPath is empty")
		return
	}
	if !file.IsExist(confPath) {
		glog.Error("configuration file:%s isn`t exist.", confPath)
		log.Fatalf("configuration file:%s isn`t exist.", confPath)
		err = fmt.Errorf("confPath:%s not exist", confPath)
		return
	}

	// read configure file.
	configContent, err := file.ToTrimString(confPath)
	if err != nil {
		glog.Error("read config file:%s failed. error:%v", confPath, err)
		log.Fatalf("read config file:%s failed. error:%v", confPath, err)
		return
	}
	glog.Debug("configContent:%s", configContent)
	log.Printf("configContent:%s\n", configContent)

	var c ABTestConf
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		glog.Error("json unmarshal config file:%s failed. error:%v", confPath, err)
		log.Fatalf("json unmarshal config file:%s failed. error:%v", confPath, err)
		return
	}

	ablock.Lock()
	defer ablock.Unlock()
	abTestConf = &c

	glog.Debug("reload confPath:%s success. abTestConf:%v", confPath, abTestConf)
	log.Printf("reload confPath:%s success. abTestConf:%v", confPath, abTestConf)
	return
}

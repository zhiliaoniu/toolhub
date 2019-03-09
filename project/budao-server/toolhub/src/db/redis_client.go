package db

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/sumaig/glog"
)

// RedisConfig structure
type RedisConfig struct {
	Addr         string `json:"addr"`
	MaxIdle      int    `json:"maxIdle"`
	Db           int    `json:"db"`
	ConnTimeout  int    `json:"connTimeout"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
}

// RedisConnPool instance of redis connect tool
var RedisConnPool *redis.Pool

// InitRedisConnPool function init redis connect pool
func InitRedisConnPool(redisCfg *RedisConfig) {

	addr := redisCfg.Addr
	maxIdle := redisCfg.MaxIdle
	idleTimeout := 240 * time.Second
	connTimeout := time.Duration(redisCfg.ConnTimeout) * time.Millisecond
	readTimeout := time.Duration(redisCfg.ReadTimeout) * time.Millisecond
	writeTimeout := time.Duration(redisCfg.WriteTimeout) * time.Millisecond

	RedisConnPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", addr, connTimeout, readTimeout, writeTimeout)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}

	glog.Debug("redis connect pool init succeed")
}

// PingRedis function verify survival
func PingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		glog.Error("ping redis failed. error:", err)
	}

	return err
}

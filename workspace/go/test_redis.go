package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/sumaig/glog"
)

// RedisConnPool instance of redis connect tool
var RedisConnPool *redis.Pool

// InitRedisConnPool function init redis connect pool
func InitRedisConnPool() {

	addr := "221.228.106.9:4019"
	maxIdle := 10000
	idleTimeout := 240 * time.Second
	connTimeout := time.Duration(1000) * time.Millisecond
	readTimeout := time.Duration(1000) * time.Millisecond
	writeTimeout := time.Duration(1000) * time.Millisecond

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

func main() {
	InitRedisConnPool()

	//var wc sync.WaitGroup
	/*
		begin, end := 1000000, 3000000
		for {
			fmt.Println(begin, "      -------------\t:", end)

			begin1 := begin
			end1 := begin1 + 50000
			fmt.Println(begin1, "++:", end1)
			begin = end1
			wc.Add(1)
			go func(begin1, end1 int) {
				conn := RedisConnPool.Get()
				defer conn.Close()
				defer wc.Done()
				fmt.Println(begin1, "++2:", end1)
				for {
					if begin1 > end1 {
						break
					}
					begin1++
					fmt.Println("begin:", begin1, "end:", end1)
					r, err := redis.Int(conn.Do("zadd", "test1", begin1, begin1))
					if err != nil {
						fmt.Println(err)
						break
					}
					fmt.Println(r, err)
				}
			}(begin1, end1)
			if begin > end {
				break
			}
		}
	*/
	/*
			begin, end := 4000000, 4005000
			var wc sync.WaitGroup
			for {
				fmt.Println(begin, "      -------------\t:", end)

				begin1 := begin
				end1 := begin1 + 50000
				fmt.Println(begin1, "++:", end1)
				begin = end1
				wc.Add(1)
				go func(begin1, end1 int) {
					conn := RedisConnPool.Get()
					defer conn.Close()
					defer wc.Done()
					for {
						if begin1 > end1 {
							break
						}
						begin1++
						fmt.Println("begin:", begin1, "end:", end1)
						r, err := redis.Int(conn.Do("zadd", "test2", begin1, begin1))
						if err != nil {
							fmt.Println(err)
							break
						}
						fmt.Println(r, err)
					}
				}(begin1, end1)
				if begin > end {
					break
				}
			}
		wc.Wait()
	*/
	/*
		wc.Add(1)
		go func() {
			conn := RedisConnPool.Get()
			begin := time.Now()
			i := 0
			for {
				if i > 1000 {
					break
				}
				r, err := redis.Int(conn.Do("zrank", "test2", 4001001))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(i, "\t", r)
				i++
			}
			t := time.Now().Sub(begin)
			fmt.Println("read 1000 time from 5000 cap. cost time:", t)
			wc.Done()
		}()

		wc.Add(1)
		go func() {
			conn := RedisConnPool.Get()
			begin := time.Now()
			i := 0
			for {
				if i > 1000 {
					break
				}
				r, err := redis.Int(conn.Do("zrank", "test1", 1001001))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(i, "\t", r)
				i++
			}
			t := time.Now().Sub(begin)
			fmt.Println("read 1000 time from 2000000 cap. cost time:", t)
			wc.Done()
		}()
		wc.Wait()
	*/
	/*
		conn := RedisConnPool.Get()
		r, err := conn.Do("zrange", "test2", 1000, 10000)
		fmt.Printf("r:%v, err:%v\n", r, err)
		r, err = redis.Strings(conn.Do("zrange", "test2", 1000, 10000))
		fmt.Printf("r:%v, err:%v\n", r, err)
	*/
}

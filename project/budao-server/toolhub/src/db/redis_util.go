package db

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sumaig/glog"
)

//----------------------------operate string---------------------------------

// IsKeyExist check if the specified field exists in the string key
func IsKeyExist(key string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return value, err
	}

	return value, err
}

// GetString get string value of key
func GetString(key string) (string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do("get", key))
	if err != nil {
		return "", err
	}

	return value, nil
}

// SetString set string value of key
func SetString(key, value string) error {
	conn := RedisConnPool.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("set", key, value))
	if err != nil {
		return err
	}

	return nil
}

func DelKey(key string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do("del", key))
	if err != nil {
		return value, err
	}

	return value, nil
}

func DelMultiKey(keys []interface{}) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do("del", keys))
	if err != nil {
		return value, err
	}

	return value, nil
}

// SetNXWithExpire set the value of key only if key does not exist with expire
func SetNXWithExpire(key, value string, expire int) (string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("set", key, value, "EX", expire, "NX"))
	if err != nil {
		return r, err
	}

	return r, nil
}

// SetNX set the value of key only if key does not exist
func SetNX(key, value string) (string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("set", key, value, "NX"))
	if err != nil {
		return r, err
	}

	return r, nil
}

//----------------------------operate hash---------------------------------
func HKeys(key string) ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HSetNX set the value of the hash table field only if the field field does not exist
func HSetNX(key, field, value string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("HSETNX", key, field, value))
	if err != nil {
		return r, err
	}

	return r, nil
}

func HIncrBy(key, field string, increment int) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("HINCRBY", key, field, increment))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HMSet simultaneously set multiple field-value pairs to hash table key
func HMSet(key string, args []interface{}) (string, error) {
	//glog.Debug("start HMSet")
	conn := RedisConnPool.Get()
	defer conn.Close()
	keyArr := make([]interface{}, 0)
	keyArr = append(keyArr, key)
	args = append(keyArr, args...)
	r, err := redis.String(conn.Do("HMSET", args...))
	//glog.Debug("r:%v, err:%v", r, err)
	if err != nil {
		return r, err
	}

	return r, nil
}

// HMGet get Multiple member from redis hash
func HMGet(key string, args []interface{}) ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	keyArr := make([]interface{}, 0)
	keyArr = append(keyArr, key)
	args = append(keyArr, args...)
	//glog.Debug("hmget %v", args)
	r, err := redis.Strings(conn.Do("HMGET", args...))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HMGet get Multiple member from redis hash
func HMGetInt(key string, args []interface{}) ([]int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	keyArr := make([]interface{}, 0)
	keyArr = append(keyArr, key)
	args = append(keyArr, args...)
	//glog.Debug("hmget %v", args)
	r, err := redis.Ints(conn.Do("HMGET", args...))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HGetAll get all memver of hash
func HGetAll(key string) (map[string]int64, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.Int64Map(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// HSet set the value of the hash table field
//If the field already exists in the hash table, the old value will be overwritten
func HSet(key string, field string, value interface{}) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("HSET", key, field, value))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HGet get a member from redis hash
func HGet(key, field string) (string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("HGET", key, field))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HGetInt get a member from redis hash
func HGetInt(key, field string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("HGET", key, field))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HDelete delete the field in hash
func HDelete(key string, field string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("HDEL", key, field))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HDelete delete the field in hash
func HMDelete(key string, args []interface{}) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	keyArr := make([]interface{}, 0)
	keyArr = append(keyArr, key)
	args = append(keyArr, args...)
	r, err := redis.Int(conn.Do("HDEL", args...))
	if err != nil {
		return r, err
	}

	return r, nil
}

// HExists check if the specified field exists in the hash table key
func HExists(key, field string) (bool, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("HEXISTS", key, field))
	if err != nil {
		return false, err
	}

	return r == 1, err
}

// HLen returns the number of field in the hash
func HLen(key string) (uint64, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.Uint64(conn.Do("HLEN", key))
	if err != nil {
		return ret, err
	}

	return ret, nil
}

//----------------------------operate list---------------------------------

//----------------------------operate set---------------------------------

// Sismember determine if a member element is a member of a collection
func Sismember(key, member string) (bool, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("SISMEMBER", key, member))
	if err != nil {
		return false, err
	}

	return r == 1, err
}

// Smembers return all members of set
func Smembers(key string) ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// SAdd add a member into set
func SAdd(key, member string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("SADD", key, member))
	if err != nil {
		return r, err
	}

	return r, err
}

// SUnionStore union multi set to one set
func SUnionStore(destination, key string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("SUNIONSTORE", destination, key))
	if err != nil {
		return r, err
	}

	return r, err
}

// SPop pop a member from set randmon
func SPop(key string) (string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("SPOP", key))
	if err != nil {
		return r, err
	}

	return r, err
}

// SAddMulti add multiple members once
func SAddMulti(key string, args []interface{}) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	keyArr := make([]interface{}, 0)
	keyArr = append(keyArr, key)
	args = append(keyArr, args...)
	r, err := redis.Int(conn.Do("SADD", args...))
	if err != nil {
		return r, err
	}
	glog.Debug("ZAddMulti success")

	return r, err
}

// SCard returns the number of elements in the collection
func SCard(key string) (uint64, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.Uint64(conn.Do("SCARD", key))
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// SScan get all member from redis set multiple times
func SScan(key string, count int) ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()

	iter := 0
	var keys []string
	var allKeys []string
	for {
		if arr, err := redis.MultiBulk(conn.Do("SSCAN", key, iter, "COUNT", count)); err != nil {
			glog.Error("sscan redis set error:%v", err)
			return nil, err
		} else {
			iter, _ = redis.Int(arr[0], nil)
			keys, _ = redis.Strings(arr[1], nil)
		}
		allKeys = append(allKeys, keys...)
		if iter == 0 {
			break
		}
	}

	return allKeys, nil
}

// SRandMember returns multiple random elements
func SRandMember(key string) ([]string, error) {
	return SRandMemberByNum(key, 10)
}

func SRandMemberByNum(key string, num int) ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.Strings(conn.Do("SRANDMEMBER", key, num))
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// SDiff return difference
func SDiff(keyone, keytwo string) ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.Strings(conn.Do("SDIFF", keyone, keytwo))
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// SRandSingleMember returns a random element
func SRandSingleMember(key string) (string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.String(conn.Do("SRANDMEMBER", key))
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// SRem delete multiple member
func SRem(key string, args []interface{}) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	keyArr := make([]interface{}, 0)
	keyArr = append(keyArr, key)
	args = append(keyArr, args...)
	r, err := redis.Int(conn.Do("SREM", args...))
	glog.Debug("r:%v, err:%v", r, err)
	if err != nil {
		return r, err
	}

	return r, nil
}

// SRemOneMember delete one member
func SRemOneMember(key string, member string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("SREM", member))
	glog.Debug("r:%v, err:%v", r, err)
	if err != nil {
		return r, err
	}

	return r, nil
}

//----------------------------operate zset---------------------------------
func ZRevRank(key, member string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("ZREVRANK", key, member))
	if err != nil {
		return r, err
	}

	return r, err
}

// Zrange return an ordered set into members within a specified range by indexed intervals
func Zrange(key string, begin, end int) ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Strings(conn.Do("ZRANGE", key, begin, end))
	if err != nil {
		return r, err
	}

	return r, err
}

// ZRevRange return an ordered(from large to small) set into members within a specified range by indexed intervals
func ZRevRange(key string, begin, end int64) ([]string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Strings(conn.Do("ZREVRANGE", key, begin, end))
	if err != nil {
		return r, err
	}

	return r, err
}

// ZAddMulti add multiple members once
func ZAddMulti(key string, args []interface{}) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	keyArr := make([]interface{}, 0)
	keyArr = append(keyArr, key)
	args = append(keyArr, args...)
	r, err := redis.Int(conn.Do("ZADD", args...))
	if err != nil {
		return r, err
	}
	//glog.Debug("ZAddMulti success")

	return r, err
}

// ZRangeWithScores get all memver of sort set
func ZRangeWithScores(key string) (map[string]int64, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()

	ret, err := redis.Int64Map(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		glog.Error("zrange tithscores redis set error:%v", err)
		return nil, err
	}
	/*
		for value, score := range ret {
			fmt.Println(value)
			fmt.Println(score)
			fmt.Println("**************")
		}
	*/
	return ret, nil
}

// ZAdd add a member into redis
func ZAdd(key string, score int64, member string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("ZADD", key, score, member))
	if err != nil {
		return r, err
	}

	return r, err
}

// ZExists check if the specified field exists in sort set
func ZExists(key, field string) (bool, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.String(conn.Do("ZSCORE", key, field))
	if err != nil {
		return false, err
	}

	return ret != "", nil
}

// ZDelete delete the field in sort set
func ZDelete(key string, field string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("ZREM", key, field))
	if err != nil {
		return r, err
	}

	return r, nil
}

func ZMDelete(key string, members []interface{}) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("ZREM", key, members))
	if err != nil {
		return r, err
	}

	return r, nil
}

// ZCard returns the number of elements in the collection
func ZCard(key string) (uint64, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.Uint64(conn.Do("ZCARD", key))
	if err != nil {
		return ret, err
	}

	return ret, nil
}

func ZScore(key, member string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	ret, err := redis.Int(conn.Do("ZSCORE", key, member))
	if err != nil {
		return ret, err
	}

	return ret, nil
}

//----------------------------operate key---------------------------------

func KExpire(key string, seconds int) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("EXPIRE", key, seconds))
	if err != nil {
		return r, err
	}

	return r, nil
}

// KDelete delete the key of redis
func KDelete(key string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()

	r, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		return r, err
	}

	return r, nil
}

// KExists check exist
func KExists(key string) (bool, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return r == 1, err
}

// KRenameNX rename key
func KRenameNX(oldKey, newKey string) (int, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.Int(conn.Do("RENAMENX", oldKey, newKey))
	if err != nil {
		return r, err
	}

	return r, nil
}

// KRename rename key
func KRename(oldKey, newKey string) (string, error) {
	conn := RedisConnPool.Get()
	defer conn.Close()
	r, err := redis.String(conn.Do("RENAME", oldKey, newKey))
	if err != nil {
		return r, err
	}

	return r, nil
}

package redisz

import (
	"github.com/futurez/litego/logger"

	"github.com/garyburd/redigo/redis"
)

//set string
func (r *RedisPool) Set(key string, value interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		logger.Warn("SET ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

//get slice
func (r *RedisPool) Get(key string) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		logger.Warn("GET ", r.server, " ", r.name, " ", err.Error())
		return nil, err
	}
	return val, nil
}

//get int
func (r *RedisPool) GetInt64(key string) (int64, error) {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		logger.Warn("GET ", r.server, " ", r.name, " ", err.Error())
		return -1, err
	}
	return val, nil
}

func (r *RedisPool) Incr(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		logger.Warn("INCR ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) IncrBy(key string, value int64) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("INCRBY", key, value))
	if err != nil {
		logger.Warn("INCRBY ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) Decr(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(redis.Int64(conn.Do("DECR", key)))
	if err != nil {
		logger.Warn("DECR ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) DecrBy(key string, value int64) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("DECRBY", key, value))
	if err != nil {
		logger.Warn("DECRBY ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) DecrByFloat(key string, value float64) float64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Float64(conn.Do("DECRBYFLOAT", key, value))
	if err != nil {
		logger.Warn("DECRBYFLOAT ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) Append(key, value string) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("APPEND", key, value))
	if err != nil {
		logger.Warn("APPEND ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) StrLen(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("STRLEN", key))
	if err != nil {
		logger.Warn("STRLEN ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) Mget(keys ...string) []string {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}

	vals, err := redis.Strings(conn.Do("MGET", args...))
	if err != nil {
		logger.Warn("MGET ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

func (r *RedisPool) Mset(kv map[string]string) error {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(kv))
	for k, v := range kv {
		args = append(args, k, v)
	}

	_, err := conn.Do("MSET", args)
	if err != nil {
		logger.Warn("MSET ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return nil
}

//Returns the bit value at offset in the string value stored at key. (index 0)
func (r *RedisPool) GetBit(key string, index int64) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("GETBIT", key, index))
	if err != nil {
		logger.Warn("GETBIT ", r.server, " ", r.name, " ", err.Error())
		return 0
	}
	return val
}

//Sets or clears the bit at offset in the string value stored at key. return old value.
func (r *RedisPool) SetBit(key string, index, value int64) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("SETBIT", key, index, value))
	if err != nil {
		logger.Warn("GETBIT ", r.server, " ", r.name, " ", err.Error())
		return 0
	}
	return val
}

//BITCOUNT foo
//BITCOUNT foo 0 1
func (r *RedisPool) BitCount(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("BITCOUNT", key))
	if err != nil {
		logger.Warn("BITCOUNT ", r.server, " ", r.name, " ", err.Error())
		return 0
	}
	return val
}

func (r *RedisPool) GetRange(key string, start, end int) string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("GETRANGE", key, start, end))
	if err != nil {
		logger.Warn("GETRANGE ", r.server, " ", r.name, " ", err.Error())
		return ""
	}
	return val
}

//atomically set key to value and return the old value stored at key.
func (r *RedisPool) GetSet(key, value string) string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("GETSET", key, value))
	if err != nil {
		logger.Warn("GETSET ", r.server, " ", r.name, " ", err.Error())
		return ""
	}
	return val
}

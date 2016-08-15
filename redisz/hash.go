package redisz

import (
	"github.com/futurez/litego/logger"

	"github.com/garyburd/redigo/redis"
)

/*
hash use to save object:
	key = objecttype:ID
	field   = object attr
	value = object attr value

for example:
	car:2 = color -> white
			name  -> audi
			price -> 90W
*/
func (r *RedisPool) Hset(key, field string, value interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()

	//insert return 1; update return 0
	_, err := redis.Int64(conn.Do("HSET", key, field, value))
	if err != nil {
		logger.Warn("HSET ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

func (r *RedisPool) Hget(key, field string) string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("HGET", key, field))
	if err != nil {
		logger.Warn("HGET ", r.server, " ", r.name, " ", err.Error())
		return ""
	}
	return val
}

func (r *RedisPool) Hmset(key string, kv map[string]string) error {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(kv)*2+1)
	args = append(args, key)
	for k, v := range kv {
		args = append(args, k, v)
	}

	_, err := conn.Do("HMSET", args...)
	if err != nil {
		logger.Warn("HMSET ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

func (r *RedisPool) Hmget(key string, fields []string) []string {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}

	vals, err := redis.Strings(conn.Do("HMGET", args...))
	if err != nil {
		logger.Warn("HMGET ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

func (r *RedisPool) Hgetall(key string) (map[string]string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	valmap, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		logger.Warn("HGETALL ", r.server, " ", r.name, " ", err.Error())
		return nil, err
	}
	return valmap, nil
}

//HEXISTS if exist return 1
func (r *RedisPool) Hexists(key, field string) bool {
	conn := r.pool.Get()
	defer conn.Close()

	v, err := redis.Int64(conn.Do("HEXISTS", key, field))
	if err != nil {
		logger.Warn("HEXISTS ", r.server, " ", r.name, " ", err.Error())
		return false
	}
	return v == 1
}

//HSETNX Sets field in the hash stored at key to value, only if field does not yet exist
//1 if field is a new field in the hash and value was set.
//0 if filed already exists in the hash and no operation was performed.
func (r *RedisPool) Hsetnx(key, field, value string) int {
	conn := r.pool.Get()
	defer conn.Close()

	v, err := redis.Int(conn.Do("HSETNX", key, field, value))
	if err != nil {
		logger.Warn("HSETNX ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return v
}

func (r *RedisPool) Hincrby(key, field string, increment int64) (int64, error) {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("HINCRBY", key, field, increment))
	if err != nil {
		logger.Warn("HINCRBY ", r.server, " ", r.name, " ", err.Error())
		return -1, err
	}
	return val, nil
}

func (r *RedisPool) Hdel(key string, fields ...string) error {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}

	_, err := redis.Int64(conn.Do("HDEL", args...))
	if err != nil {
		logger.Warn("HDEL ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

func (r *RedisPool) Hkeys(key string) []string {
	conn := r.pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		logger.Warn("HKEYS ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return keys
}

func (r *RedisPool) Hvals(key string) []string {
	conn := r.pool.Get()
	defer conn.Close()

	vals, err := redis.Strings(conn.Do("HVALS", key))
	if err != nil {
		logger.Warn("HVALS ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

func (r *RedisPool) Hlen(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("HLEN", key))
	if err != nil {
		logger.Warn("HLEN ", r.server, " ", r.name, " ", err.Error())
		return 0
	}
	return val
}

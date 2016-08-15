package redisz

import (
	"github.com/futurez/litego/logger"

	"github.com/garyburd/redigo/redis"
)

func (r *RedisPool) Sadd(key string, values ...string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	for _, value := range values {
		args = append(args, value)
	}

	val, err := redis.Int(conn.Do("SADD", args...))
	if err != nil {
		logger.Warn("SADD ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) Srem(key string, values ...string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	for _, value := range values {
		args = append(args, value)
	}

	val, err := redis.Int(conn.Do("SREM", args...))
	if err != nil {
		logger.Warn("SREM ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) SmembersString(key string) ([]string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	vals, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		logger.Warn("SMEMBERS ", r.server, " ", r.name, " ", err.Error())
		return nil, err
	}
	return vals, nil
}

func (r *RedisPool) SmembersInt(key string) ([]int, error) {
	conn := r.pool.Get()
	defer conn.Close()

	vals, err := redis.Ints(conn.Do("SMEMBERS", key))
	if err != nil {
		logger.Warn("SMEMBERS ", r.server, " ", r.name, " ", err.Error())
		return nil, err
	}
	return vals, nil
}

func (r *RedisPool) SIsmember(key, member string) bool {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("SISMEMBER", key, member))
	if err != nil {
		logger.Warn("SMEMBERS ", r.server, " ", r.name, " ", err.Error())
		return false
	}
	return val == 1
}

func (r *RedisPool) Sdiff(keys ...string) []string {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}

	vals, err := redis.Strings(conn.Do("SDIFF", args...))
	if err != nil {
		logger.Warn("SDIFF ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

func (r *RedisPool) Sinter(keys ...string) []string {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}

	vals, err := redis.Strings(conn.Do("SINTER", args...))
	if err != nil {
		logger.Warn("SINTER ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

func (r *RedisPool) Sunion(keys ...string) []string {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}

	vals, err := redis.Strings(conn.Do("SUNION", args...))
	if err != nil {
		logger.Warn("SUNION ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

func (r *RedisPool) SdiffStore(dest string, keys ...string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, 1+len(keys))
	args = append(args, dest)
	for _, key := range keys {
		args = append(args, key)
	}

	val, err := redis.Int(conn.Do("SDIFFSTORE", args...))
	if err != nil {
		logger.Warn("SDIFFSTORE ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) SinterStore(dest string, keys ...string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, 1+len(keys))
	args = append(args, dest)
	for _, key := range keys {
		args = append(args, key)
	}

	val, err := redis.Int(conn.Do("SINTERSTORE", args...))
	if err != nil {
		logger.Warn("SINTERSTORE ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) SunionStore(dest string, keys ...string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, 1+len(keys))
	args = append(args, dest)
	for _, key := range keys {
		args = append(args, key)
	}

	val, err := redis.Int(conn.Do("SUNIONSTORE", args...))
	if err != nil {
		logger.Warn("SUNIONSTORE ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

//set size
func (r *RedisPool) Scard(key string) int {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("SCARD", key))
	if err != nil {
		logger.Warn("SCARD ", r.server, " ", r.name, " ", err.Error())
		return 0
	}
	return val
}

//When called with just the key argument, return a random element from the set value stored at key.
//when called with the additional count argument, return an array of count distinct elements if count is positive.
//If called with a negative count the behavior changes and the command is allowed to return the same element multiple times.
func (r *RedisPool) SrandMember(key string, count int) ([]string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	vals, err := redis.Strings(conn.Do("SRANDMEMBER", key, count))
	if err != nil {
		logger.Warn("SRANDMEMBER ", r.server, " ", r.name, " ", err.Error())
		return nil, err
	}
	return vals, nil
}

func (r *RedisPool) Spop(key string) string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("SPOP", key))
	if err != nil {
		logger.Warn("SPOP ", r.server, " ", r.name, " ", err.Error())
		return ""
	}
	return val
}

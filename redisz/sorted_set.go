package redisz

import (
	"github.com/futurez/litego/logger"

	"github.com/garyburd/redigo/redis"
)

func (r *RedisPool) Zadd(key string, score int, member string) int {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("ZADD", key, score, member))
	if err != nil {
		logger.Warn("ZADD ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) Zadds(key string, scoremap map[int]string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(scoremap)*2+1)
	args = append(args, key)
	for k, v := range scoremap {
		args = append(args, k, v)
	}

	val, err := redis.Int(conn.Do("ZADD", args...))
	if err != nil {
		logger.Warn("ZADD ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

//get member score
func (r *RedisPool) Zscore(key string, member string) int {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("ZSCORE", key, member))
	if err != nil {
		logger.Warn("ZSCORE ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

//get member
func (r *RedisPool) Zrange(key string, start, end int) []string {
	conn := r.pool.Get()
	defer conn.Close()

	vals, err := redis.Strings(conn.Do("ZRANGE", key, start, end))
	if err != nil {
		logger.Warn("ZRANGE ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

func (r *RedisPool) ZrangeMap(key string, start, end int) map[string]string {
	conn := r.pool.Get()
	defer conn.Close()

	valMap, err := redis.StringMap(conn.Do("ZRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		logger.Warn("ZRANGE ... [WITHSCORES] ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return valMap
}

//从大到小 get member
func (r *RedisPool) Zrevrrange(key string, start, end int) []string {
	conn := r.pool.Get()
	defer conn.Close()

	vals, err := redis.Strings(conn.Do("ZREVRRANGE", key, start, end))
	if err != nil {
		logger.Warn("ZREVRRANGE ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

func (r *RedisPool) ZrevrrangeMap(key string, start, end int) map[string]string {
	conn := r.pool.Get()
	defer conn.Close()

	valMap, err := redis.StringMap(conn.Do("ZREVRRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		logger.Warn("ZREVRRANGE ... [WITHSCORES] ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return valMap
}

func (r *RedisPool) ZrangebylexLimit(key string, offset, count int64) []string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Strings(conn.Do("ZRANGEBYSCORE", key, "-inf", "+inf", "limit", offset, count))
	if err != nil {
		logger.Warn("ZRANGEBYSCORE ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return val
}

//del mem
func (r *RedisPool) Zrem(key string, member ...string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(member)+1)
	args = append(args, key)
	for _, value := range member {
		args = append(args, value)
	}

	val, err := redis.Int(conn.Do("ZREM", args...))
	if err != nil {
		logger.Warn("ZREM ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) Zrems(key string, members []string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(members)+1)
	args = append(args, key)
	for _, value := range members {
		args = append(args, value)
	}

	val, err := redis.Int(conn.Do("ZREM", args...))
	if err != nil {
		logger.Warn("ZREM ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) Zcard(key string) int {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("ZCARD", key))
	if err != nil {
		logger.Warn("ZCARD ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

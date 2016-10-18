package redisz

import (
	"strings"
	"time"

	"github.com/futurez/litego/logger"

	"github.com/garyburd/redigo/redis"
)

type RedisPool struct {
	name   string
	server string
	pool   *redis.Pool
}

// server  : 192.168.1.141:6379
// password: ""
// maxIdle : 1
func NewRedisPool(name, server, password string, maxIdle int) *RedisPool {
	if maxIdle < 1 {
		maxIdle = 1
	}
	if strings.Count(server, ":") == 0 {
		server += ":6379"
	}

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			//logger.Info("new redis connection...")
			c, err := redis.Dial("tcp", server)
			if err != nil {
				logger.Error("NewPool Dail err=", err.Error())
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					logger.Error("NewPool Dail Do `AUTH` err=", err.Error())
					return nil, err
				}
			}
			return c, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},

		MaxIdle:     maxIdle,
		MaxActive:   maxIdle * 500, //active = 2 * idle
		IdleTimeout: 240 * time.Second,
		Wait:        true,
	}

	return &RedisPool{name, server, pool}
}

func NewRedisSelectPool(name, server, password string, maxIdle int, id int) *RedisPool {
	if maxIdle < 1 {
		maxIdle = 1
	}
	if strings.Count(server, ":") == 0 {
		server += ":6379"
	}

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			//logger.Info("new redis connection...")
			c, err := redis.Dial("tcp", server)
			if err != nil {
				logger.Error("NewPool Dail err=", err.Error())
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					logger.Error("NewPool Dail Do `AUTH` err=", err.Error())
					return nil, err
				}
			}
			c.Do("select", id)
			return c, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},

		MaxIdle:     maxIdle,
		MaxActive:   maxIdle * 2, //active = 2 * idle
		IdleTimeout: 240 * time.Second,
		Wait:        true,
	}

	return &RedisPool{name, server, pool}
}

func (r *RedisPool) Close() {
	r.pool.Close()
}

func (r *RedisPool) Dbsize() int {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("DBSIZE"))
	if err != nil {
		logger.Warn("DBSIZE ", r.server, " ", r.name, " ", err.Error())
		return 0
	}
	return val
}

func (r *RedisPool) Del(keys ...string) bool {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}

	val, err := redis.Int64(conn.Do("DEL", args...))
	if err != nil {
		logger.Warn("DEL ", r.server, " ", r.name, " ", err.Error())
		return false
	}
	return val != 0
}

//return if key exists.
func (r *RedisPool) Exists(key string) bool {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("EXISTS", key))
	if err != nil {
		logger.Warn("EXISTS ", r.server, " ", r.name, " ", err.Error())
		return false
	}
	return val == 1
}

//set a timeout on key
func (r *RedisPool) Exprie(key string, time int) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := redis.Int64(conn.Do("EXPIRE", key, time))
	if err != nil {
		logger.Warn("EXPIRE ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

//get all keys
func (r *RedisPool) Keys(pattern string) []string {
	if pattern == "" {
		pattern = "*"
	}

	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Strings(conn.Do("KEYS", pattern))
	if err != nil {
		logger.Warn("KEYS ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return val
}

//remove the existing timeout on key
func (r *RedisPool) Persist(key string) bool {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("PERSIST", key))
	if err != nil {
		logger.Warn("PERSIST ", r.server, " ", r.name, " ", err.Error())
		return false
	}
	return val == 1
}

//return the remaining time to live of a key that has a timeout.
func (r *RedisPool) Ttl(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("TTL", key))
	if err != nil {
		logger.Warn("TTL ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

// return value type
func (r *RedisPool) Type(key string) string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("TYPE", key))
	if err != nil {
		logger.Warn("TYPE ", r.server, " ", r.name, " ", err.Error())
		return "nil"
	}
	return val
}

// rename oldkey newkey
func (r *RedisPool) Rename(oldkey, newkey string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("RENAME", oldkey, newkey)
	if err != nil {
		logger.Warn("RENAME ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

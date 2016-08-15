package redisz

import (
	"github.com/futurez/litego/logger"

	"github.com/garyburd/redigo/redis"
)

//list存储一个有序的字符串列表,常用的操作是向列表两端添加元素,或者获取列表的某一个片段
//是双链表实现的,所以向列表两端添加元素的时间复杂度为O(1),获取越接近两端的元素速度就越快

//-1 if error, other is len(list)
func (r *RedisPool) Lpush(key string, values ...string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	for _, value := range values {
		args = append(args, value)
	}

	val, err := redis.Int(conn.Do("LPUSH", args...))
	if err != nil {
		logger.Warn("LPUSH ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

//-1 if error, other is len(list)
func (r *RedisPool) Rpush(key string, values ...string) int {
	conn := r.pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(values)+1)
	args = append(args, key)
	for _, value := range values {
		args = append(args, value)
	}

	val, err := redis.Int(conn.Do("RPUSH", args...))
	if err != nil {
		logger.Warn("RPUSH ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

//Removes and returns the first element of the list stored at key.
func (r *RedisPool) Lpop(key string) string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("LPOP", key))
	if err != nil {
		logger.Warn("LPOP ", r.server, " ", r.name, " ", err.Error())
		return ""
	}
	return val
}

//Removes and returns the last element of the list stored at key.
func (r *RedisPool) Rpop(key string) string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("RPOP", key))
	if err != nil {
		logger.Warn("RPOP ", r.server, " ", r.name, " ", err.Error())
		return ""
	}
	return val
}

//Returns the length of the list stored at key.
func (r *RedisPool) Llen(key string) int {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("LLEN", key))
	if err != nil {
		logger.Warn("LLEN ", r.server, " ", r.name, " ", err.Error())
		return 0
	}
	return val
}

//Returns the specified elements of the list stored at key.
//The return value includes the rightmost element.
func (r *RedisPool) LrangeString(key string, start, stop int) []string {
	conn := r.pool.Get()
	defer conn.Close()

	vals, err := redis.Strings(conn.Do("LRANGE", key))
	if err != nil {
		logger.Warn("LRANGE ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

//Returns the specified elements of the list stored at key.
//The return value includes the rightmost element.
func (r *RedisPool) LrangeInt(key string, start, stop int) []int {
	conn := r.pool.Get()
	defer conn.Close()

	vals, err := redis.Ints(conn.Do("LRANGE", key))
	if err != nil {
		logger.Warn("LRANGE ", r.server, " ", r.name, " ", err.Error())
		return nil
	}
	return vals
}

//Removes the first count occurrences of elements equal to value from the list stored at key.
//The count argument influences the operation in the following ways:
//
//count > 0: Remove elements equal to value moving from head to tail.
//count < 0: Remove elements equal to value moving from tail to head.
//count = 0: Remove all elements equal to value.
func (r *RedisPool) Lrem(key string, count int, value string) int {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.Int(conn.Do("LREM", key, count, value))
	if err != nil {
		logger.Warn("LREM ", r.server, " ", r.name, " ", err.Error())
		return 0
	}
	return val
}

//
func (r *RedisPool) Lindex(key string, index int) string {
	conn := r.pool.Get()
	defer conn.Close()

	val, err := redis.String(conn.Do("LINDEX", key, index))
	if err != nil {
		logger.Warn("LINDEX ", r.server, " ", r.name, " ", err.Error())
		return ""
	}
	return val
}

func (r *RedisPool) Lset(key string, index int, value string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("LSET", key, index)
	if err != nil {
		logger.Warn("LSET ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

func (r *RedisPool) Ltrim(key string, start, end int) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("LTRIM", key, start, end)
	if err != nil {
		logger.Warn("LTRIM ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

func (r *RedisPool) Linsert(key string, before bool, pivot, value string) int {
	conn := r.pool.Get()
	defer conn.Close()

	var (
		val int
		err error
	)

	if before {
		val, err = redis.Int(conn.Do("LINSERT", key, "BEFORE", pivot, value))
	} else {
		val, err = redis.Int(conn.Do("LINSERT", key, "AFTER", pivot, value))
	}

	if err != nil {
		logger.Warn("LINSERT ", r.server, " ", r.name, " ", err.Error())
		return -1
	}
	return val
}

func (r *RedisPool) RpopLpush(key string, srcList, destList string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("RPOPLPUSH", key, srcList, destList)
	if err != nil {
		logger.Warn("RPOPLPUSH ", r.server, " ", r.name, " ", err.Error())
		return err
	}
	return nil
}

package RedisPool

import (
	"Common/redisz"
	"errors"
	"fmt"

	_ "github.com/garyburd/redigo/redis"
)

var redisPool *redisz.RedisPool

func InitRedis(ip string, runCount int) {
	redisPool = redisz.NewRedisPool("common", ip, "", runCount)
}

func RedisDBSize() {
	size := redisPool.Dbsize()
	fmt.Printf("size is %d \n", size)
}

func RedisReadBytes(key string) ([]byte, error) {
	return redisPool.Get(key)
}

func RedisReadStringMap(key string) (map[string]string, error) {
	return redisPool.Hgetall(key)
}

func RedisReadString(key string, timesava int) (string, error) {
	data, err := redisPool.Get(key)
	if timesava > 0 {
		err = redisPool.Exprie(key, timesava*3600)
		if err != nil {
			return "", err
		}
	}
	return string(data), nil
}

func RedisReadInt64(key string, timesava int) (int64, error) {
	val, err := redisPool.GetInt64(key)
	if err != nil {
		return 0, err
	}
	if timesava > 0 {
		err = redisPool.Exprie(key, timesava*3600)
		if err != nil {
			return 0, err
		}
	}
	return val, nil
}

func RedisWrite(key string, data interface{}, timesava int) error {
	err := redisPool.Set(key, data)
	if err != nil {
		return err
	}
	if timesava > 0 {
		err = redisPool.Exprie(key, timesava*3600)
		if err != nil {
			return err
		}
	}
	return nil
}

func RedisWriteMapField(key, field string, data interface{}, timesava int) error {
	err := redisPool.Hset(key, field, data)
	if err != nil {
		return err
	}
	if timesava > 0 {
		err = redisPool.Exprie(key, timesava*3600)
		if err != nil {
			return err
		}
	}
	return nil
}

func RedisListRPush(key string, data string, timesava int) error {
	ret := redisPool.Rpush(key, data)
	if ret < 0 {
		return errors.New("Redis RPush Fail")
	}
	if timesava > 0 {
		err := redisPool.Exprie(key, timesava*3600)
		if err != nil {
			return err
		}
	}
	return nil
}

func RedisListRPop(key string) (string, error) {
	retData := redisPool.Rpop(key)
	if retData == string("") {
		return retData, errors.New("No Data")
	}

	return retData, nil
}

//var myPool chan redis.Conn
//var add string
//var count int

//func InitRedis(ip string, runCount int) {
//	logger.Info("InitRedis", ip)
//	add = ip
//	count = runCount
//}

//func getRedis() redis.Conn {
//	if myPool == nil {
//		myPool = make(chan redis.Conn, count)
//	}

//	createFunc := func() {
//		for i := 0; i < count/2; i++ {
//			for {
//				conn, err := redis.DialTimeout("tcp", add, 10*time.Second, 1*time.Second, 1*time.Second)
//				if err != nil {
//					logger.Error("connect to redis error", err)
//					time.Sleep(3 * time.Second)
//				} else {
//					putRedis(conn)
//					break
//				}
//			}
//		}
//	}

//	if len(myPool) == 0 {
//		go createFunc()
//	}

//	for {
//		select {
//		case p := <-myPool:
//			s, err := p.Do("PING")
//			if s == "PONG" && err == nil {
//				return p
//			} else {
//				p.Close()
//			}
//		case <-time.After(5 * time.Second):
//			return nil
//		}
//	}
//}

//func putRedis(conn redis.Conn) {
//	if myPool == nil {
//		myPool = make(chan redis.Conn, count)
//	}
//	if len(myPool) == count {
//		conn.Close()
//		return
//	}
//	myPool <- conn
//}

//func closeConnector(conn redis.Conn) {
//	conn.Close()
//}

//func RedisDBSize() {
//	conn := getRedis()
//	if conn == nil {
//		return
//	}
//	defer putRedis(conn)
//	size, _ := conn.Do("DBSIZE")
//	fmt.Printf("size is %d \n", size)
//}

//func RedisReadBytes(key string) ([]byte, error) {
//	conn := getRedis()
//	if conn == nil {
//		return nil, errors.New("get redis error")
//	}
//	defer putRedis(conn)
//	data, err := redis.Bytes(conn.Do("GET", key))
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}

//func RedisReadStringMap(key string) (map[string]string, error) {
//	conn := getRedis()
//	if conn == nil {
//		return nil, errors.New("get redis error")
//	}
//	defer putRedis(conn)
//	data, err := redis.StringMap(conn.Do("HGETALL", key))
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}

//func RedisReadString(key string, timesava int) (string, error) {
//	conn := getRedis()
//	if conn == nil {
//		return "", errors.New("get redis error")
//	}
//	defer putRedis(conn)
//	data, err := redis.String(conn.Do("GET", key))
//	if err != nil {
//		return "", err
//	}
//	if timesava > 0 {
//		_, err = conn.Do("EXPIRE", key, timesava*3600)
//		if err != nil {
//			logger.Error("set ket time error", err)
//			return "", err
//		}
//	}
//	return data, nil
//}

//func RedisReadInt64(key string, timesava int) (int64, error) {
//	conn := getRedis()
//	if conn == nil {
//		return 0, errors.New("get redis error")
//	}
//	defer putRedis(conn)
//	data, err := redis.Int64(conn.Do("GET", key))
//	if err != nil {
//		return 0, err
//	}
//	if timesava > 0 {
//		_, err = conn.Do("EXPIRE", key, timesava*3600)
//		if err != nil {
//			logger.Error("set ket time error", err)
//			return 0, err
//		}
//	}
//	return data, nil
//}

//func RedisWrite(key string, data interface{}, timesava int) error {
//	conn := getRedis()
//	if conn == nil {
//		return errors.New("get redis error")
//	}
//	defer putRedis(conn)
//	_, err := conn.Do("SET", key, data)
//	if err != nil {
//		logger.Error("save redis error", err)
//		return err
//	}
//	_, err = conn.Do("EXPIRE", key, timesava*3600)
//	if err != nil {
//		logger.Error("set ket time error", err)
//		return err
//	}
//	return nil

//	//	logger.Debug(n)
//	//	if n == int64(1) {
//	//		n, err := conn.Do("EXPIRE", key, timesava*3600)
//	//		if n == int64(1) {
//	//			return nil
//	//		} else {
//	//			logger.Error("set ket time error", err)
//	//			return err
//	//		}
//	//	} else if n == int64(0) {
//	//		return errors.New("the key has already existed")
//	//	}
//	//	return errors.New(fmt.Sprintf("error is %lld", n))
//}

//func RedisWriteMapField(key, field string, data interface{}, timesava int) error {
//	conn := getRedis()
//	if conn == nil {
//		return errors.New("get redis error")
//	}
//	defer putRedis(conn)
//	_, err := conn.Do("HSET", key, field, data)
//	if err != nil {
//		logger.Error("save redis error", err)
//		return err
//	}
//	if timesava > 0 {
//		_, err = conn.Do("EXPIRE", key, timesava*3600)
//		if err != nil {
//			logger.Error("set ket time error", err)
//			return err
//		}
//	}
//	return nil
//}

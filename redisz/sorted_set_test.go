package redisz

import (
	"testing"
)

func TestZadd(t *testing.T) {
	redisPool := NewRedisSelectPool("hash", "192.168.1.10:6379", "", 10, 9)

	testKey := "soretedset"
	redisPool.Del(testKey)

	testMap := make(map[int]string)
	testMap[1] = "USA"
	testMap[2] = "China"
	testMap[3] = "English"
	testMap[4] = "Japan"
	testMap[5] = "SSS"
	testMap[6] = "SSSA"
	ret := redisPool.Zadds(testKey, testMap)
	if ret < 0 {
		t.Error("Zadd failed.")
		return
	}

	score := redisPool.Zscore(testKey, "English")
	if score != 3 {
		t.Error("Zscore failed. score=", score)
		return
	}
	t.Log("return line=", ret, ", score=", score)

	if rets := redisPool.Zrange(testKey, 0, -1); rets != nil {
		t.Log("rets=", rets)
	}

	if rets := redisPool.ZrangeMap(testKey, 0, -1); rets != nil {
		t.Log("rets=", rets)
	}

	t.Log(redisPool.Zscore(testKey, "呵呵呵"))

	ret = redisPool.Zadd(testKey, 100, "呵呵呵")
	if rets := redisPool.Zrange(testKey, 0, -1); rets != nil {
		t.Log("rets=", rets)
	}

	t.Log(redisPool.Zscore(testKey, "呵呵呵"))
	t.Log(redisPool.Zcard(testKey))

	if rets := redisPool.ZrangebylexLimit(testKey, 0, 2); rets != nil {
		t.Log("rets=", rets)
	}
	if rets := redisPool.ZrangebylexLimit(testKey, 1, 2); rets != nil {
		t.Log("rets=", rets)
	}

	ret = redisPool.Zrem(testKey, "SSS")
	t.Log(ret)
	if rets := redisPool.Zrange(testKey, 0, -1); rets != nil {
		t.Log("rets=", rets)
	}

	ret = redisPool.Zrem(testKey, "SSS", "English", "China")
	t.Log(ret)
	if rets := redisPool.Zrange(testKey, 0, -1); rets != nil {
		t.Log("rets=", rets)
	}

	var mems []string
	mems = append(mems, "呵呵呵")
	mems = append(mems, "呵呵呵")
	ret = redisPool.Zrems(testKey, mems)
	t.Log(ret)
	if rets := redisPool.Zrange(testKey, 0, -1); rets != nil {
		t.Log("rets=", rets)
	}

	if rets := redisPool.ZrangebylexLimit(testKey, 0, 3); rets != nil {
		t.Log("rets=", rets)
	}
	if rets := redisPool.ZrangebylexLimit(testKey, 1, 3); rets != nil {
		t.Log("rets=", rets)
	}
	t.Log(redisPool.Zcard(testKey))
	//	redisPool.Exprie(testKey, 100)
}

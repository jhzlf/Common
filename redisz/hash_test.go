package redisz

//import (
//	"fmt"
//	"strconv"
//	"testing"

//	"github.com/futurez/litego/util"
//)

////  并发测试
////	var wait sync.WaitGroup
////	for i := 0; i < 20; i++ {
////		wait.Add(1)
////		go func(n int) {
////			defer wait.Done()
////			key := "car" + strconv.FormatInt(int64(n), 10)
////			redisPool.Hset(key, "price", strconv.FormatInt(500, 10))
////			redisPool.Hset(key, "name", "BMW")
////		}(i)
////	}
////	wait.Wait()

//func TestHsetAndHget(t *testing.T) {
//	redisPool := NewRedisPool("hash", "192.168.1.141:6379", "", 3)
//	if err := redisPool.Hset("car", "price", strconv.FormatInt(500, 10)); err != nil {
//		t.Error("HSET car price 500 failed!")
//		return
//	}

//	val, err := strconv.Atoi(redisPool.Hget("car", "price"))
//	if err != nil {
//		t.Error("HGET car price failed!")
//		return
//	}

//	if val != 500 {
//		t.Error("HGET car price not empty 500")
//		return
//	}
//	t.Log("Test HSET, HGET success, val=", val)

//	redisPool.Close()
//}

//func TestHmsetAndHmget(t *testing.T) {
//	redisPool := NewRedisPool("hash", "192.168.1.141:6379", "", 3)

//	kv := make(map[string]string)
//	kv["price"] = "300"
//	kv["name"] = "BMW"

//	if redisPool.Hmset("carm", kv) != nil {
//		t.Error("HMSET car  failed!")
//		return
//	}

//	var fields []string
//	fields = append(fields, "price", "name")
//	vals := redisPool.Hmget("carm", fields)
//	if vals == nil {
//		t.Error("HMGET failed")
//		return
//	}
//	t.Log("vals = ", vals)
//	if vals[0] != "300" {
//		t.Error("HGET car price not empty 500")
//		return
//	}

//	if vals[1] != "BMW" {
//		t.Error("HGET car name not BMW")
//		return
//	}
//	t.Log("Test HMSET, HMGET success")

//	redisPool.Close()
//}

//func TestHgetAll(t *testing.T) {
//	redisPool := NewRedisPool("hash", "192.168.1.141:6379", "", 3)

//	kv := make(map[string]string)
//	kv["price"] = "300"
//	kv["name"] = "BMW"

//	if redisPool.Hmset("cara", kv) != nil {
//		t.Error("HMSET car  failed!")
//		return
//	}

//	valMap, err := redisPool.Hgetall("carm")
//	if err != nil {
//		t.Error("HGETALL failed")
//		return
//	}
//	t.Log("valMap = ", valMap)
//	v1, ok := valMap["price"]
//	if !ok || v1 != "300" {
//		t.Error("HGETALL price not equal 300")
//		return
//	}

//	v2, ok := valMap["name"]
//	if !ok || v2 != "BMW" {
//		t.Error("HGETALL name not equal 'BMW'")
//		return
//	}

//	t.Log("Test HGETALL success")
//	redisPool.Close()
//}

//func TestCommon(t *testing.T) {
//	redisPool := NewRedisPool("hash", "192.168.1.141:6379", "", 1)

//	postId := redisPool.Incr("posts:count")

//	slugName := util.UUID()
//	t.Log("postId=", postId, " slugName=", slugName)
//	retSlug := redisPool.Hsetnx("slug.to.id", slugName, strconv.FormatInt(postId, 10))
//	if retSlug <= 0 {
//		t.Error("Hsetnx slug.to.id error")
//		return
//	}

//	postId_01, _ := strconv.ParseInt(redisPool.Hget("slug.to.id", slugName), 10, 0)
//	if postId != postId_01 {
//		t.Error(postId, " != ", postId_01)
//		return
//	}

//	postKey := fmt.Sprintf("post:%d", postId)
//	kv := make(map[string]string)
//	kv["title"] = "hello"
//	kv["content"] = "hello world"
//	kv["slug"] = slugName
//	if err := redisPool.Hmset(postKey, kv); err != nil {
//		t.Error(err.Error())
//		return
//	}

//	postMap, err := redisPool.Hgetall(postKey)
//	if err != nil {
//		t.Error("hgetall ", postKey, "error")
//		return
//	}
//	t.Log("title = ", postMap["title"])

//	newSlugName := util.UUID()
//	retSlug = redisPool.Hsetnx("slug.to.id", newSlugName, strconv.FormatInt(postId, 10))
//	if retSlug <= 0 {
//		t.Error("Hsetnx slug.to.id error, ", newSlugName)
//		return
//	}

//	oldSlug := redisPool.Hget(postKey, "slug")
//	if oldSlug != slugName {
//		t.Error(oldSlug, " != ", slugName)
//		return
//	}

//	t.Log("old =", oldSlug, ", new=", newSlugName)
//	if err := redisPool.Hset(postKey, "slug", newSlugName); err != nil {
//		t.Error(err.Error())
//		return
//	}

//	if err := redisPool.Hdel("slug.to.id", slugName); err != nil {
//		t.Error(err.Error())
//		return
//	}
//	redisPool.Close()
//}

//func TestHkeysAndHvals(t *testing.T) {
//	redisPool := NewRedisPool("hash", "192.168.1.141:6379", "", 1)

//	keys := redisPool.Hkeys("slug.to.id")
//	if keys == nil {
//		t.Error("HKEYS error")
//		return
//	}
//	t.Log("keys = ", keys)

//	vals := redisPool.Hvals("slug.to.id")
//	if vals == nil {
//		t.Error("HKEYS error")
//		return
//	}
//	t.Log("vals = ", vals)

//	l := redisPool.Hlen("slug.to.id")
//	if int(l) != len(keys) {
//		t.Error("slug.to.id len error")
//		return
//	}
//	t.Log("len = ", l)
//}

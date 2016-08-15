package redisz

//import (
//	"testing"
//)

//func TestSaddAndSrem(t *testing.T) {
//	//	redisPool := NewRedisPool("hash", "192.168.1.141:6379", "", 1)

//	//	redisPool.Sadd("POST:42:tags", "闲言碎语", "技术文章", "Java")
//	//	redisPool.Srem("POST:42:tags", "闲言碎语")
//	//	tags, _ := redisPool.SmembersString("POST:42:tags")
//	//	if tags[0] != "技术文章" && tags[0] != "Java" {
//	//		t.Error(tags[0], "error")
//	//		return
//	//	}

//	//	if tags[1] != "技术文章" && tags[1] != "Java" {
//	//		t.Error(tags[1], "error")
//	//		return
//	//	}
//}

//func TestSInter(t *testing.T) {
//	redisPool := NewRedisPool("hash", "192.168.1.141:6379", "", 1)

//	redisPool.Sadd("post:1:tags", "Java")
//	redisPool.Sadd("post:2:tags", "Java", "MySQL")
//	redisPool.Sadd("post:3:tags", "Java", "MySQL", "Redis")

//	redisPool.Sadd("tag:Redis:posts", "3")
//	redisPool.Sadd("tag:MySQL:posts", "2", "3")
//	redisPool.Sadd("tag:Java:posts", "1", "2", "3")

//	vals := redisPool.Sinter("tag:Redis:posts", "tag:MySQL:posts", "tag:Java:posts")
//	if len(vals) != 1 {
//		t.Error("values = ", vals)
//		return
//	}

//	key := "post:" + vals[0] + ":tags"
//	members, _ := redisPool.SmembersString(key)
//	t.Log("key", key, "members : ", members)
//}

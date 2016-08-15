package redisz

//	"fmt"
//	"testing"

//	"github.com/futurez/litego/util"

const (
	PostsPerPage = 10
)

//func TestListCommon(t *testing.T) {
//	redisPool := NewRedisPool("hash", "192.168.1.141:6379", "", 1)

//	curPage := int(util.RandRange(0, 10))
//	start := (curPage - 1) * PostsPerPage
//	end := curPage*PostsPerPage - 1
//	postsIds := redisPool.LrangeInt("posts:list", start, end)

//	for _, postsId := range postsIds {
//		key := fmt.Sprintf("post:%d", postsId)
//		post, _ := redisPool.Hgetall(key)
//		t.Log("doc titel = ", post["title"])
//	}

//}

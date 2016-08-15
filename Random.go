package Common

import (
	"math/rand"
	"time"
)

func GetRandInt(scope int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(scope)
}

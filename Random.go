package Common

import (
	"math/rand"
	"time"
)

func GetRandInt(scope int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(scope)
}

func MakeRandCode(index, kind int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var s string
	for i := 0; i < index; i++ {
		n := r.Intn(kind)
		switch n {
		case 0:
			m := r.Intn(9)
			b := 0x30 + m
			s += string(b)
		case 1:
			m := r.Intn(25)
			b := 0x61 + m
			s += string(b)
		case 2:
			m := r.Intn(25)
			b := 0x41 + m
			s += string(b)
		}
	}
	return s
}

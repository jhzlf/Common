package Common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

const (
	MAX_UTF8_BITS = 6
)

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func TrimUtf8String(str []byte) {
	size := len(str)
	for i := size - 1; i > size-MAX_UTF8_BITS; i-- {
		if ((str[i])&0x80) == 0x00 || ((str[i])&0x40) == 0x40 {
			str = str[:i]
			break
		}
	}
}

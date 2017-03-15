package Common

import (
	"Common/logger"
	"bytes"
	"crypto/des"
	"errors"
)

func Decrypt_DES(key []byte, crypted []byte) []byte {
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		logger.Debug(err)
		return nil
	}
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil
	}
	out := make([]byte, len(crypted))
	dst := out
	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}
	out, err = PKCS5UnPadding(out)
	if err != nil {
		return nil
	}
	return out
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length > unpadding {
		return origData[:(length - unpadding)], nil
	}
	return nil, errors.New("slice bounds out of range")
}

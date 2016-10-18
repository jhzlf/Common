package Common

import (
	"bytes"
	"crypto/des"
	"fmt"
	"log"
)

func Decrypt(des3Key []byte, crypted []byte) []byte {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	block, err := des.NewTripleDESCipher(des3Key)
	if err != nil {
		log.Println(err)
		return nil
	}

	// blockMode := cipher.NewCBCDecrypter(block, des3Key[:8])

	origData := make([]byte, len(crypted))
	blockSize := block.BlockSize()
	for i := 0; i < len(origData)/blockSize; i++ {
		// log.Println(i)
		block.Decrypt(origData[blockSize*i:blockSize*(i+1)], crypted[blockSize*i:blockSize*(i+1)])
	}

	// blockMode.CryptBlocks(origData, crypted)

	origData = pKCS5UnPadding(origData)

	return origData
}

func Encrypt(des3Key []byte, origData []byte) []byte {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	block, err := des.NewTripleDESCipher(des3Key)
	if err != nil {
		fmt.Println("NewTripleDESCipher error")
		return nil

	}
	// fmt.Printf("origData: %v \n", origData)
	// fmt.Printf("des3_key: %v \n", des3_key)

	// fmt.Printf("block size: %d", block.BlockSize())
	origData = pKCS5Padding(origData, block.BlockSize())

	// blockMode := cipher.NewCBCEncrypter(block, des3Key[:8])

	crypted := make([]byte, len(origData))

	// blockMode.CryptBlocks(crypted, origData)

	blockSize := block.BlockSize()
	for i := 0; i < len(origData)/blockSize; i++ {

		block.Encrypt(crypted[blockSize*i:blockSize*(i+1)], origData[blockSize*i:blockSize*(i+1)])
	}
	// block.Encrypt(crypted, origData)
	// fmt.Printf("encrypt:%v \n", crypted)
	// fmt.Printf("enc tmpbody:%v \n", crypted)
	return crypted
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	// log.Println(len(ciphertext), blockSize, padding)

	// if padding == blockSize {
	// 	return ciphertext
	// }

	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	// log.Println(padtext)
	return append(ciphertext, padtext...)
}

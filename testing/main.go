package main

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"fmt"
)

func encrypt(key, plaintext []byte) ([]byte, error) {
	cipherBlock, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	paddedMsg := padMsg(plaintext, des.BlockSize)
	cipherText := make([]byte, len(paddedMsg)+des.BlockSize)
	iv := cipherText[:des.BlockSize]
	rand.Read(iv)
	fmt.Println(len(paddedMsg))
	blockMode := cipher.NewCBCEncrypter(cipherBlock, iv)
	blockMode.CryptBlocks(cipherText[des.BlockSize:], paddedMsg)
	return cipherText, nil
}

func padMsg(plaintext []byte, blockSize int) []byte {

	if len(plaintext)%blockSize == 0 {
		return plaintext
	}
	padLength := len(plaintext) % blockSize
	return append(plaintext, make([]byte, blockSize-padLength)...)
}

func main() {
	encrypt([]byte("12344321"), []byte("Today I met my crush, what a hunk"))
}

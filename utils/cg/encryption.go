package cg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// 对配置文件选项 进行加密解密

// Encrypt 加密函数
func Encrypt(text, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, []byte(key)[:aes.BlockSize])
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密函数
func Decrypt(encryptedText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	ciphertext, _ := base64.StdEncoding.DecodeString(encryptedText)
	cfb := cipher.NewCFBDecrypter(block, []byte(key)[:aes.BlockSize])
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}

// GenerateAESKey 随机生成 AES key
func GenerateAESKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

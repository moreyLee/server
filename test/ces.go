package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// 对配置文件选项 进行加密解密

// 加密函数
func encrypt(text, key string) (string, error) {
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

// 解密函数
func decrypt(encryptedText, key string) (string, error) {
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

// 随机生成 AES key
func generateAESKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}
func main() {
	//key, _ := generateAESKey(32) // 256 位密钥
	key := "61afa6ddca1285f640b03c8c64d02c31"
	//fmt.Printf("Generated AES key: %s\n", hex.EncodeToString(key))
	encryptedString := "11700ee17be3621da8bb4443e073763a69"
	//botToken := "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
	//webhook := "https://a904-217-165-23-20.ngrok-free.app/jenkins/telegram-webhook"
	//url := "http://jenkins1.3333d.vip/"
	//webhook := "https://2156-8-218-67-135.ngrok-free.app/jenkins/telegram-webhook"
	// 加密
	encrypted, err := encrypt(encryptedString, key)
	if err != nil {
		fmt.Println("加密失败:", err)
		return
	}
	fmt.Println("加密后的字符串:", encrypted)

	// 解密
	decrypted, err := decrypt(encrypted, key)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}
	fmt.Println("解密后的字符串:", decrypted)
}

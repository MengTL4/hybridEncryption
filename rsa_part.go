package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

// EncryptWithPublicKey 使用公钥加密消息
func EncryptWithPublicKey(publicKeyPath string, message []byte) ([]byte, error) {
	// 读取公钥文件
	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取公钥文件: %v", err)
	}

	// 解析公钥PEM数据
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	if publicKeyBlock == nil {
		return nil, fmt.Errorf("无法解析公钥PEM数据")
	}

	// 解析公钥
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("无法解析公钥: %v", err)
	}

	// 使用公钥进行加密
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), message)
	if err != nil {
		return nil, fmt.Errorf("加密错误: %v", err)
	}

	return encryptedMessage, nil
}

// DecryptWithPrivateKey 使用私钥解密消息
func DecryptWithPrivateKey(privateKeyPath string, encryptedMessageBase64 string) ([]byte, error) {
	// 读取私钥文件
	privateKeyPEM, _ := os.ReadFile(privateKeyPath)
	// 解析私钥PEM数据
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	// 解析私钥
	privateKey, _ := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	// 解码Base64编码的密文
	encryptedMessage, _ := base64.StdEncoding.DecodeString(encryptedMessageBase64)
	// 使用私钥进行解密
	decryptedMessage, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedMessage)

	return decryptedMessage, nil
}

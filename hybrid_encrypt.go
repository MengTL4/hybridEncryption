package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	aes2 "hybrid_encryption/aes"
	"hybrid_encryption/utils"
)

type Args struct {
	IV   string `json:"iv"`
	Key  string `json:"key"`
	Data string `json:"data"`
}

func HybridEncrypt(key []byte, plaintext string) ([]byte, map[int][]byte, map[int][]byte, map[int][]byte, map[int][]byte) {
	// 加密
	//ciphertext, err := AesEncrypt(key, plaintext)
	// 随机生成128bit的IV
	iv := make([]byte, 16) // 128bit = 16byte
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	aes, _ := aes2.NewAES(key)
	ciphertext, roundStates, roundStates2, roundStates3, roundStates4 := aes.EncryptCBC([]byte(plaintext), iv, utils.PKCS7Padding)
	//base64只为方便输出
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)
	encryptedKey, err := EncryptWithPublicKey("public.pem", key)
	encryptedKeyBase64 := base64.StdEncoding.EncodeToString(encryptedKey)
	ivBase64 := base64.StdEncoding.EncodeToString(iv)
	args := Args{
		IV:   ivBase64,
		Key:  encryptedKeyBase64,
		Data: ciphertextBase64,
	}
	args2, err := json.Marshal(args)
	if err != nil {
		fmt.Println("JSON序列化错误", err)
		return nil, nil, nil, nil, nil
	}
	return args2, roundStates, roundStates2, roundStates3, roundStates4
}

func HybridDecrypt(jsonData []byte) string {
	var newArgs Args
	err := json.Unmarshal(jsonData, &newArgs)
	if err != nil {
		fmt.Println("JSON反序列化错误:", err)
		return ""
	}
	decryptedKey, err := DecryptWithPrivateKey("private.pem", newArgs.Key)
	if err != nil {
		fmt.Println("RSA解密发生错误:", err)
		return ""
	}
	decryptedTextBase64, _ := base64.StdEncoding.DecodeString(newArgs.Data)
	decryptedIv, _ := base64.StdEncoding.DecodeString(newArgs.IV)
	//decryptedText, err := AesDecrypt(decryptedKey, decryptedTextBase64)
	aes, _ := aes2.NewAES(decryptedKey)
	decryptedText := aes.DecryptCBC(decryptedTextBase64, decryptedIv, utils.PKCS7Unpadding)
	if err != nil {
		fmt.Println("AES解密错误:", err)
		return ""
	}
	return string(decryptedText)
}

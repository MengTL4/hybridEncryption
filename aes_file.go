package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"hybrid_encryption/utils"
	"io"
	"log"
	"os"
	"path/filepath"
)

func encryptFile(filepathString, fileName string) string {
	// Reading plaintext file
	plainText, err := os.ReadFile(filepathString)
	if err != nil {
		log.Fatalf("read file err: %v", err.Error())
	}

	key := []byte("SpP5j4DACmh5uEWR")
	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}

	// Generating random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
	}

	// Decrypt file
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	filePathString := filepath.Join(utils.Mkdir("encrypt_upload"), "/", "encrypt"+fileName)
	// Writing ciphertext file
	err = os.WriteFile(filePathString, cipherText, 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
	}
	return filePathString

}

func decryptFile(filepathString, fileName string) string {
	// Reading ciphertext file
	cipherText, err := os.ReadFile(filepathString)
	if err != nil {
		log.Fatal(err)
	}

	// Reading key
	key := []byte("SpP5j4DACmh5uEWR")

	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("cipher err: %v", err.Error())
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("cipher GCM err: %v", err.Error())
	}

	// Deattached nonce and decrypt
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Fatalf("decrypt file err: %v", err.Error())
	}

	// Writing decryption content
	filePathString := filepath.Join(utils.Mkdir("decrypt_upload"), "/", "decrypt"+fileName)
	err = os.WriteFile(filePathString, plainText, 0777)
	if err != nil {
		log.Fatalf("write file err: %v", err.Error())
	}
	return filePathString
}

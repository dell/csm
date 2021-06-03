package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

// TODO: approach should be sanitized
var cipherKey = []byte("thisshouldbeparameterized32bytes")

func GetEnv(envName string, defaultValue string) string {
	val := os.Getenv(envName)
	if val == "" {
		return defaultValue
	}
	return val
}

func EncryptPassword(plainText []byte) ([]byte, error) {
	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	return gcm.Seal(nonce, nonce, plainText, nil), nil
}

func DecryptPassword(cipherText []byte) ([]byte, error) {
	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return []byte{}, err
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return []byte{}, err
	}
	return plainText, nil
}

// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
	"strings"

	"github.com/dell/csm-deployment/utils/constants"

	goYAML "github.com/go-yaml/yaml"
)

const (
	// RequiredCipherKeyLength is the required length, in bytes, of the cipher key
	RequiredCipherKeyLength = 32
)

// CipherKey is used to encrypt/decrypt passwords
var CipherKey []byte

// GetEnv - Utility method to get Environment variable
var GetEnv = func(envName string, defaultValue string) string {
	val := os.Getenv(envName)
	if val == "" {
		return defaultValue
	}
	return val
}

// EncryptPassword - Method to encrypts password
func EncryptPassword(plainText []byte) ([]byte, error) {
	c, err := aes.NewCipher(CipherKey)
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

// DecryptPassword - Method to decipher password
var DecryptPassword = func(cipherText []byte) ([]byte, error) {
	c, err := aes.NewCipher(CipherKey)
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

// SplitYAML divides a big bytes of yaml files in individual yaml files.
func SplitYAML(gaintYAML []byte) ([][]byte, error) {
	decoder := goYAML.NewDecoder(bytes.NewReader(gaintYAML))
	nullByte := []byte{110, 117, 108, 108, 10} /* byte returned by  goYAML when yaml is empty*/

	var res [][]byte
	for {
		var value interface{}
		if err := decoder.Decode(&value); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		valueBytes, err := goYAML.Marshal(value)
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(valueBytes, nullByte) {
			res = append(res, valueBytes)
		}
	}
	return res, nil
}

// Find returns true is val is in slice
func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// GetValueFromMetadataKey returns value corresponding to the key in the given metadata string
func GetValueFromMetadataKey(metaData, key string) string {

	for _, data := range strings.Split(metaData, constants.MetadataDelimeter) {
		dataSplit := strings.Split(data, constants.MetadataSeparator)
		if len(dataSplit) > 1 && key == dataSplit[0] {
			return dataSplit[1]
		}
	}
	return ""
}

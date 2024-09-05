/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package crypto provides encryption and decryption functionality using AES-GCM encryption mode.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const noneLen = 12

// Encrypter is an interface that defines methods for encrypting and decrypting text.
type Encrypter interface {
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
}

type encryption struct {
	key []byte
}

// NewEncryption creates a new instance of an Encrypter with the provided encryption key.
func NewEncryption(key []byte) Encrypter {
	return &encryption{key: key}
}

// Encrypt is used to encrypt text
func (e *encryption) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, noneLen)
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(text), nil)
	ciphertext = append(ciphertext, nonce...)

	return hex.EncodeToString(ciphertext), nil
}

// Decrypt is used to decrypt text
func (e *encryption) Decrypt(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	plain, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}

	if len(plain) < noneLen {
		return "", fmt.Errorf("index is negative:%d", len(plain)-noneLen)
	}

	nonce := plain[len(plain)-noneLen:]
	plain = plain[:len(plain)-noneLen]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, plain, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

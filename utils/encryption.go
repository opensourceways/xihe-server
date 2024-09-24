package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

type SymmetricEncryption interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

func NewSymmetricEncryption(key, nonce string) (SymmetricEncryption, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	se := symmetricEncryption{aead: gcm}

	if nonce != "" {
		nonce1, err := hex.DecodeString(nonce)
		if err != nil {
			return nil, err
		}
		if len(nonce1) != gcm.NonceSize() {
			return nil, errors.New(
				"the length of nonce for symmetric encryption is unmatched",
			)
		}
		se.nonce = nonce1
	}

	return se, nil
}

type symmetricEncryption struct {
	aead  cipher.AEAD
	nonce []byte
}

func (se symmetricEncryption) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := se.nonce
	if nonce == nil {
		nonce = make([]byte, se.aead.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}
	}

	return se.aead.Seal(nonce, nonce, plaintext, nil), nil
}

func (se symmetricEncryption) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := se.aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return se.aead.Open(nil, nonce, ciphertext, nil) // #nosec G407
}

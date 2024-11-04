package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"github.com/xdg-go/pbkdf2"
)

const (
	pbkdf2Iterations = 10000
	pbkdf2KeyLength  = 128
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

// EncodeToken encodes the given token using the provided salt and returns the encoded token as a string.
func EncodeToken(token string, salt string) (string, error) {
	saltByte, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	encBytes := pbkdf2.Key([]byte(token), saltByte, pbkdf2Iterations, pbkdf2KeyLength, sha256.New)
	enc := base64.RawStdEncoding.EncodeToString(encBytes)

	return enc, nil
}

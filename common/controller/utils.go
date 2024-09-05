/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides utility functions for handling HTTP errors and error codes.
package controller

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
)

const (
	userAgent        = "User-Agent"
	pbkdf2Iterations = 10000
	pbkdf2KeyLength  = 32
)

// GetIp returns the client IP address from the given gin context.
func GetIp(ctx *gin.Context) (string, error) {
	return ctx.ClientIP(), nil
}

// GetUserAgent returns a primitive.UserAgent object based on the given gin context.
func GetUserAgent(ctx *gin.Context) (primitive.UserAgent, error) {
	return primitive.CreateUserAgent(""), nil
	// return primitive.NewUserAgent(ctx.GetHeader(userAgent))
}

// SetCookie sets a cookie with the given key, value, httpOnly flag, and expiry time in the given gin context.
func SetCookie(ctx *gin.Context, key, val, domain string, httpOnly bool, expiry *time.Time, sameSite http.SameSite) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    val,
		Path:     "/",
		HttpOnly: httpOnly,
		Secure:   true,
		SameSite: sameSite,
		Domain:   domain,
	}

	if expiry != nil {
		cookie.Expires = *expiry
	}

	http.SetCookie(ctx.Writer, cookie)
}

// GetCookie retrieves the value of the cookie with the given key from the given gin context.
func GetCookie(ctx *gin.Context, key string) (string, error) {
	cookie, err := ctx.Request.Cookie(key)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
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

func GetSaltHash(s string) string {
	h := sha256.New()

	bs := []byte(s)
	salt := generateRandomSalt(4)
	bs = append(bs, salt...)

	h.Write(bs)

	return hex.EncodeToString(h.Sum(nil))
}

func generateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		logrus.Errorf("generate sailt error: %v", err)
	}

	return salt
}

/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type RandomId interface {
	RandomId() string
}

func NewRandomId() (RandomId, error) {
	str, err := genRandomId()

	return randomId(str), err
}

func CreateRandomId(v string) RandomId {
	return randomId(v)
}

func ToRandomId(v string) (RandomId, error) {
	bytes, err := hex.DecodeString(v)
	if err != nil || len(bytes) != randomIdLength {
		return nil, errors.New("invalid id")
	}

	return randomId(v), nil
}

type randomId string

func (r randomId) RandomId() string {
	return string(r)
}

func genRandomId() (string, error) {
	k := make([]byte, randomIdLength)
	if _, err := rand.Read(k); err != nil {
		return "", err
	}

	return hex.EncodeToString(k), nil
}

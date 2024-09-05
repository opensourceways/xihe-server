/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"crypto/rand"
)

const (
	charNum        = 26
	digitalNum     = 10
	firstCharOfPw  = 0x21
	firstLowercase = 'a'
	lastLowercase  = 'z'
	firstUppercase = 'A'
	lastUppercase  = 'Z'
	firstDigital   = '0'
	lastDigital    = '9'
)

var (
	pwRange = byte(0x7E - 0x20)
)

type Password interface {
	Password() string
	Clear()
}

// NewPassword creates a new Password instance.
func NewPassword() (Password, error) {
	bt, err := passwordInstance.gen()
	if err != nil {
		return nil, err
	}

	return password(bt), nil
}

type password []byte

func (p password) Password() string {
	return string(p)
}

func (p password) Clear() {
	for i := range p {
		p[i] = 0
	}
}

func newPasswordImpl(cfg PasswordConfig) *passwordImpl {
	return &passwordImpl{
		cfg: cfg,
	}
}

type passwordImpl struct {
	cfg PasswordConfig
}

func (impl *passwordImpl) gen() ([]byte, error) {
	var bytes = make([]byte, impl.cfg.MinLength)

	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	items := make([]byte, len(bytes))
	for k, v := range bytes {
		items[k] = firstCharOfPw + v%pwRange
	}

	if impl.goodFormat(items) {
		return items, nil
	}

	for i := range bytes {
		items[i] = impl.genChar(i, bytes[i])
	}

	return items, nil
}

func (impl *passwordImpl) genChar(i int, v byte) byte {
	switch i % 3 {
	case 0:
		return byte(firstLowercase) + v%charNum
	case 1:
		return byte(firstUppercase) + v%charNum
	default:
		return byte(firstDigital) + v%digitalNum
	}
}

func (impl *passwordImpl) goodFormat(s []byte) bool {
	return impl.hasMultiChars(s) && !impl.hasConsecutive(s)
}

func (impl *passwordImpl) hasMultiChars(s []byte) bool {
	part := make([]bool, 4)

	for _, c := range s {
		if c >= firstLowercase && c <= lastLowercase {
			part[0] = true
		} else if c >= firstUppercase && c <= lastUppercase {
			part[1] = true
		} else if c >= firstDigital && c <= lastDigital {
			part[2] = true
		} else {
			part[3] = true
		}
	}

	i := 0
	for _, b := range part {
		if b {
			i++
		}
	}

	return i >= impl.cfg.MinNumOfCharKind
}

func (impl *passwordImpl) hasConsecutive(str []byte) bool {
	count := 1
	for i := 1; i < len(str); i++ {
		if str[i] == str[i-1] {
			if count++; count > impl.cfg.MinNumOfConsecutiveChars {
				return true
			}
		} else {
			count = 1
		}
	}

	return false
}

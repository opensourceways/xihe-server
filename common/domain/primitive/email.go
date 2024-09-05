/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"fmt"
	"net/mail"
)

// Email is an interface that represents an email address.
type Email interface {
	Email() string
}

// NewUserEmail creates a new Email instance with the given value.
func NewUserEmail(v string) (Email, error) {
	if v == "" {
		return dpEmail(v), nil
	}
	if len(v) > emailConfig.MaxLength {
		return nil, fmt.Errorf("invalid email address length")
	}

	if !emailConfig.regexp.MatchString(v) {
		return nil, fmt.Errorf("invalid email address match")
	}

	if v[0] == '-' {
		return nil, fmt.Errorf("invalid email address, first character can't be -")
	}

	if _, err := mail.ParseAddress(v); err != nil {
		return nil, fmt.Errorf("invalid  RFC 5322 email address")
	}

	return dpEmail(v), nil
}

// CreateUserEmail creates a new Email instance without validating the email address.
func CreateUserEmail(v string) Email {
	return dpEmail(v)
}

type dpEmail string

// Email returns the email address as a string.
func (r dpEmail) Email() string {
	return string(r)
}

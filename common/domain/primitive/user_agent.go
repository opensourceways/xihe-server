/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"strings"
)

const (
	merlin = "merlin"
)

// UserAgent is an interface that represents a user agent.
type UserAgent interface {
	UserAgent() string
}

// NewUserAgent creates a new user agent with the given value.
func NewUserAgent(v string) (UserAgent, error) {
	v = strings.ToLower(v)

	if v != merlin {
		return nil, errors.New("unknown user agent")
	}

	return userAgent(v), nil
}

// CreateUserAgent creates a new user agent with the given value.
func CreateUserAgent(v string) UserAgent {
	return userAgent(v)
}

type userAgent string

// UserAgent returns the string representation of the user agent.
func (r userAgent) UserAgent() string {
	return string(r)
}

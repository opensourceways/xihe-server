/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"regexp"
)

var usccRegexp *regexp.Regexp

// Init initializes the USCC package with the provided configuration.
func Init(cfg Config) (err error) {
	usccRegexp, err = regexp.Compile(cfg.USCCRegexp)

	return
}

// USCC is an interface representing unified social credit code
type USCC interface {
	USCC() string
}

// NewUSCC create a new USCC instance from a string value.
func NewUSCC(v string) (USCC, error) {
	if !usccRegexp.MatchString(v) {
		return nil, errors.New("invalid unified social credit code")
	}

	return uscc(v), nil
}

// CreateUSCC create a new USCC instance from a string value.
func CreateUSCC(v string) USCC {
	return uscc(v)
}

// uscc represents unified social credit code
type uscc string

// USCC represents unified social credit code
func (u uscc) USCC() string {
	return string(u)
}

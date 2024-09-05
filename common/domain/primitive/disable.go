/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"strings"
)

const (
	IllegalContent = "illegal_content"
)

var (
	ReasonIllegalContent = disablereason(IllegalContent)
)

// DisableReason is an interface that defines the disable reason of an object.
type DisableReason interface {
	DisableReason() string
}

// NewDisableReason creates a new DisableReason instance based on the given string.
func NewDisableReason(v string) (DisableReason, error) {
	v = strings.ToLower(v)
	if v != IllegalContent {
		return nil, errors.New("unknown reason")
	}

	return disablereason(v), nil
}

// CreateDisableReason creates a new DisableReason instance based on the given string.
func CreateDisableReason(v string) DisableReason {
	return disablereason(v)
}

type disablereason string

// DisableReason returns the disablereason as a string.
func (r disablereason) DisableReason() string {
	return string(r)
}

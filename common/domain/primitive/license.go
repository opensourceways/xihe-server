/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"strings"
)

// License is an interface representing a license.
type License interface {
	License() []string
}

// NewLicense creates a new License instance from a string value.
func NewLicense(v string) (License, error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if v == "" || !allLicenses[v] {
		return nil, errors.New("unsupported license")
	}

	return license([]string{v}), nil
}

// CreateLicense creates a new License instance directly from a string value.
func CreateLicense(v []string) License {
	return license(v)
}

type license []string

// License returns the string representation of the license.
func (r license) License() []string {
	return []string(r)
}

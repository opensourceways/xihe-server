/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"fmt"
)

// MSDName is an interface representing a name.
type MSDName interface {
	MSDName() string
	FirstLetter() byte
}

// NewMSDName creates a new MSDName instance from a string value.
func NewMSDName(v string) (MSDName, error) {
	n := len(v)
	if n > msdConfig.MaxNameLength || n < msdConfig.MinNameLength {
		return nil, fmt.Errorf("invalid name length, should between %d and %d",
			msdConfig.MinNameLength, msdConfig.MaxNameLength)
	}

	if !msdConfig.nameRegexp.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return msdName(v), nil
}

// CreateMSDName creates a new MSDName instance directly from a string value.
func CreateMSDName(v string) MSDName {
	return msdName(v)
}

type msdName string

// MSDName returns the string representation of the name.
func (r msdName) MSDName() string {
	return string(r)
}

// FirstLetter returns the first letter of the name as a byte.
func (r msdName) FirstLetter() byte {
	return string(r)[0]
}

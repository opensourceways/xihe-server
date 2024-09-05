/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import "strconv"

// Identity is an interface that represents an identity with both a string and integer representation.
type Identity interface {
	Identity() string
	Integer() int64
}

// NewIdentity creates a new Identity instance from a string value.
func NewIdentity(v string) (Identity, error) {
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, err
	}

	return identity(n), nil
}

// CreateIdentity creates a new Identity instance from an integer value.
func CreateIdentity(v int64) Identity {
	return identity(v)
}

type identity int64

// Identity returns the string representation of the identity.
func (r identity) Identity() string {
	return strconv.FormatInt(int64(r), 10)
}

// Integer returns the integer representation of the identity.
func (r identity) Integer() int64 {
	return int64(r)
}

/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import "errors"

const (
	legalRepresentative = "法定代表人"
	authorizedPerson    = "被授权人"
)

// Identity is an interface that represents an identity.
type Identity interface {
	Identity() string
}

// ValidateIdentity checks if the provided string value is a valid identity.
func ValidateIdentity(v string) bool {
	return v == legalRepresentative || v == authorizedPerson
}

// NewIdentity creates a new instance of the Identity type with the provided value.
func NewIdentity(v string) (Identity, error) {
	if !ValidateIdentity(v) {
		return nil, errors.New("invalid identity")
	}

	return identity(v), nil
}

// CreateIdentity creates an instance of the Identity type with the provided value.
func CreateIdentity(v string) Identity {
	return identity(v)
}

type identity string

// Identity returns the string representation of the identity.
func (i identity) Identity() string {
	return string(i)
}

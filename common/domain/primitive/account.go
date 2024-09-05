/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/xerrors"
)

// Account is an interface that represents an account.
type Account interface {
	Account() string
}

// NewAccount creates a new account with the given name.
func NewAccount(v string) (Account, error) {
	if v == "" {
		return nil, errors.New("empty name")
	}

	if accountConfig.reservedAccounts.Has(v) {
		return nil, errors.New("name is reserved")
	}

	n := len(v)
	if n > accountConfig.MaxNameLength || n < accountConfig.MinNameLength {
		return nil, fmt.Errorf("invalid name length, should between %d and %d",
			accountConfig.MinNameLength, accountConfig.MaxNameLength)
	}

	if !accountConfig.nameRegexp.MatchString(v) {
		return nil, errors.New("name can only contain alphabet, integer, _ and -")
	}

	return dpAccount(v), nil
}

// CreateAccount is usually called internally, such as repository.
func CreateAccount(v string) Account {
	return dpAccount(v)
}

type dpAccount string

// Account returns the account name.
func (r dpAccount) Account() string {
	return string(r)
}

// TokenName is an interface that represents a token name.
type TokenName interface {
	TokenName() string
}

// NewTokenName creates a new token name with the given name.
func NewTokenName(v string) (TokenName, error) {
	if v == "" {
		return nil, errors.New("empty token name")
	}

	n := len(v)
	if n > tokenConfig.MaxNameLength || n < tokenConfig.MinNameLength {
		return nil, fmt.Errorf("invalid token name length, should between %d and %d",
			tokenConfig.MinNameLength, tokenConfig.MaxNameLength)
	}

	if !tokenConfig.regexp.MatchString(v) {
		return nil, errors.New("token name can only contain alphabet, integer, _ and -")
	}

	if _, err := strconv.ParseInt(v, 0, 64); err == nil {
		return nil, xerrors.Errorf("token name %s can not only contain integer", v)
	}

	return dpTokenName(v), nil
}

// CreateTokenName is usually called internally, such as repository.
func CreateTokenName(v string) TokenName {
	return dpTokenName(v)
}

type dpTokenName string

// TokenName returns the token name.
func (r dpTokenName) TokenName() string {
	return string(r)
}

/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/utils"
)

// AccountFullname is an interface representing a full name.
type AccountFullname interface {
	AccountFullname() string
}

// NewAccountFullname creates a new AccountFullname instance from a string value.
func NewAccountFullname(v string) (AccountFullname, error) {
	if v == "" {
		return accountFullname(v), nil
	}

	if utils.StrLen(v) > accountConfig.MaxFullnameLength {
		return nil, errors.New("fullname is too long")
	}

	v = utils.XSSEscapeString(v)
	if utils.StrLen(v) > accountConfig.MaxFullnameLength {
		return nil, errors.New("fullname is too long")
	}

	return accountFullname(v), nil
}

// NewOrgFullname creates a new OrgFullname instance from a string value.TODO: add orgFullname
func NewOrgFullname(v string) (AccountFullname, error) {
	if v == "" {
		return nil, errors.New("org fullname can't be empty")
	}

	if utils.StrLen(v) > accountConfig.MaxFullnameLength {
		return nil, errors.New("org fullname is too long")
	}

	v = utils.XSSEscapeString(v)
	if utils.StrLen(v) > accountConfig.MaxFullnameLength || utils.StrLen(v) < accountConfig.MinFullnameLength {
		return nil, fmt.Errorf("invalid org fullname length, should between %d and %d",
			accountConfig.MinFullnameLength, accountConfig.MaxFullnameLength)
	}

	return accountFullname(v), nil
}

// CreateAccountFullname creates a new AccountFullname instance directly from a string value.
func CreateAccountFullname(v string) AccountFullname {
	return accountFullname(v)
}

type accountFullname string

// AccountFullname returns the string representation of the full name.
func (r accountFullname) AccountFullname() string {
	return string(r)
}

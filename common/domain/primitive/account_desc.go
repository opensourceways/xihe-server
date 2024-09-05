/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"

	"github.com/opensourceways/xihe-server/utils"
)

// AccountDesc is an interface representing a user signature or org description.
type AccountDesc interface {
	AccountDesc() string
}

// NewAccountDesc creates a new AccountDesc instance from a string value.
func NewAccountDesc(v string) (AccountDesc, error) {
	if v == "" {
		return accountDesc(v), nil
	}

	if utils.StrLen(v) > accountConfig.MaxDescLength {
		return nil, errors.New("desc is too long")
	}

	v = utils.XSSEscapeString(v)
	if utils.StrLen(v) > accountConfig.MaxDescLength {
		return nil, errors.New("desc is too long")
	}

	return accountDesc(v), nil
}

// CreateAccountDesc creates a new AccountDesc instance directly from a string value.
func CreateAccountDesc(v string) AccountDesc {
	return accountDesc(v)
}

type accountDesc string

// AccountDesc returns the string representation of the description.
func (r accountDesc) AccountDesc() string {
	return string(r)
}

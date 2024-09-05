/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"

	"github.com/opensourceways/xihe-server/utils"
)

// MSDDesc is an interface representing a description.
type MSDDesc interface {
	MSDDesc() string
}

// NewMSDDesc creates a new MSDDesc instance from a string value.
func NewMSDDesc(v string) (MSDDesc, error) {
	if v == "" {
		return msdDesc(v), nil
	}

	if utils.StrLen(v) > msdConfig.MaxDescLength {
		return nil, errors.New("desc is too long")
	}

	v = utils.XSSEscapeString(v)
	if utils.StrLen(v) > msdConfig.MaxDescLength {
		return nil, errors.New("desc is too long")
	}

	return msdDesc(v), nil
}

// CreateMSDDesc creates a new MSDDesc instance directly from a string value.
func CreateMSDDesc(v string) MSDDesc {
	return msdDesc(v)
}

type msdDesc string

// MSDDesc returns the string representation of the description.
func (r msdDesc) MSDDesc() string {
	return string(r)
}

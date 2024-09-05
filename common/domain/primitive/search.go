/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"fmt"

	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

type SearchKey interface {
	SearchKey() string
}

const (
	SearchKeyMaxLength = 100
	SearchKeyMinLength = 0

	SearchTypeMaxLength = 6
	SearchTypeMinLength = 0

	SizeMaxLength = 100
	SizeMinLength = 0

	TypeModelResultNum = 6

	SearchTypeModel   = "model"
	SearchTypeDataset = "dataset"
	SearchTypeSpace   = "space"
	SearchTypeUser    = "user"
	SearchTypeOrg     = "org"
)

func getTypeLimit() []string {
	return []string{SearchTypeModel, SearchTypeDataset, SearchTypeSpace, SearchTypeUser, SearchTypeOrg}
}

func NewSearchKey(v string) (SearchKey, error) {
	n := len(v)
	if n >= SearchKeyMaxLength || n <= SearchKeyMinLength {
		return nil, fmt.Errorf("invalid searchKey length, should between %d and %d", 0, 100)
	}

	return searchKey(v), nil
}

func CreateSearchKey(v string) SearchKey {
	return searchKey(v)
}

func (sk searchKey) SearchKey() string {
	return string(sk)
}

type searchKey string

type SearchType interface {
	SearchType() []string
}

func NewSearchType(v []string) (SearchType, error) {
	n := len(v)
	if n >= SearchTypeMaxLength || n <= SearchTypeMinLength {
		return nil, fmt.Errorf("invalid searchType length, should between %d and %d", 0, 5)
	}

	for _, s := range v {
		if !utils.Contains(getTypeLimit(), s) {
			return nil, fmt.Errorf("invalid searchType, should be one of %v", getTypeLimit())
		}
	}

	return searchType(v), nil
}

func CreateSearchType(v []string) SearchType {
	return searchType(v)
}

func (st searchType) SearchType() []string {
	return []string(st)
}

type searchType []string

type Size interface {
	Size() int
}

func NewSize(v int) (Size, error) {
	if v == 0 {
		logrus.Infof("size is 0, set to default value %d", TypeModelResultNum)
		return size(TypeModelResultNum), nil
	}

	if v >= SizeMaxLength || v <= SizeMinLength {
		return nil, fmt.Errorf("invalid size, should between %d and %d", 0, 100)
	}

	return size(v), nil
}

func CreateSize(v int) Size {
	return size(v)
}

func (s size) Size() int {
	return int(s)
}

type size int

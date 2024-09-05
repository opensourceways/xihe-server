/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package primitive

import (
	"errors"
	"strings"
)

const (
	computilityTypeNpu = "npu"
	computilityTypeCpu = "cpu"
)

// ComputilityType is an interface that defines computility hardware.
type ComputilityType interface {
	ComputilityType() string
	IsCpu() bool
	IsNpu() bool
}

// NewComputilityType creates a new ComputilityType instance decided by sdk based on the given string.
func NewComputilityType(v string) (ComputilityType, error) {
	v = strings.ToLower(v)

	switch v {
	case computilityTypeNpu:
	case computilityTypeCpu:
	default:
		return nil, errors.New("unknown computility type")
	}

	return computilityType(v), nil
}

// CreateComputilityType creates a new ComputilityType instance based on the given string.
func CreateComputilityType(v string) ComputilityType {
	return computilityType(v)
}

type computilityType string

// computilityType returns the string representation of the ComputilityType.
func (r computilityType) ComputilityType() string {
	return string(r)
}

func (r computilityType) IsCpu() bool {
	return string(r) == computilityTypeCpu
}

func (r computilityType) IsNpu() bool {
	return string(r) == computilityTypeNpu
}

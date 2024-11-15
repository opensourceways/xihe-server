/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import "errors"

const (
	itemBlock  = "block"
	itemPass   = "pass"
	itemReview = "review"
)

// SensitiveItemResult is an interface that defines the result of a sensitive item check.
type SensitiveItemResult interface {
	SensitiveItemResult() string
	IsEmpty() bool
	IsPass() bool
}

// SensitiveItemValidate returns true if the sensitive item is valid.
func SensitiveItemValidate(v string) bool {
	return v == "" || v == itemBlock || v == itemReview || v == itemPass
}

// NewSensitiveItemResult returns a new SensitiveItemResult.
func NewSensitiveItemResult(v string) (SensitiveItemResult, error) {
	if !SensitiveItemValidate(v) {
		return nil, errors.New("invalid sensitive item")
	}

	return sensitiveItemResult(v), nil
}

// EmptySensitiveItemResult returns an empty SensitiveItemResult.
func EmptySensitiveItemResult() SensitiveItemResult {
	return sensitiveItemResult("")
}

// CreateSensitiveItemResult returns a new SensitiveItemResult.
func CreateSensitiveItemResult(v string) SensitiveItemResult {
	return sensitiveItemResult(v)
}

// sensitiveItemResult is the result of a sensitive item check.
type sensitiveItemResult string

// SensitiveItemResult returns the result of the sensitive item check.
func (s sensitiveItemResult) SensitiveItemResult() string {
	return string(s)
}

// IsEmpty returns true if the sensitive item result is empty.
func (s sensitiveItemResult) IsEmpty() bool {
	return s.SensitiveItemResult() == ""
}

// IsPass returns true if the sensitive item result is pass.
func (s sensitiveItemResult) IsPass() bool {
	return s.SensitiveItemResult() == itemPass
}

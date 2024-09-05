/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
)

// Phone number, china mainland supported only for now
type Phone interface {
	PhoneNumber() string
}

// NewPhone creates a new Phone instance with the given string value.
func NewPhone(v string) (Phone, error) {
	if v == "" {
		return nil, errors.New("empty phone number")
	}

	if len(v) > phoneConfig.MaxLength {
		return nil, errors.New("phone is too laong")
	}

	if !phoneConfig.regexp.MatchString(v) {
		return nil, errors.New("invalid phone")
	}

	return phoneNumber(v), nil
}

// CreatePhoneNumber creates a new Phone instance with the given string value without validation.
func CreatePhoneNumber(v string) Phone {
	return phoneNumber(v)
}

type phoneNumber string

// PhoneNumber returns the string representation of the phone number.
func (r phoneNumber) PhoneNumber() string {
	return string(r)
}

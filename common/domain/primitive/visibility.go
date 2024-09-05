/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"strings"
)

const (
	Public  = "public"
	Private = "private"
)

var (
	VisibilityPublic  = visibility(Public)
	VisibilityPrivate = visibility(Private)
)

// Visibility is an interface that defines the visibility of an object.
type Visibility interface {
	IsPublic() bool
	IsPrivate() bool
	Visibility() string
}

// NewVisibility creates a new Visibility instance based on the given string.
func NewVisibility(v string) (Visibility, error) {
	v = strings.ToLower(v)
	if v != Public && v != Private {
		return nil, errors.New("unknown visibility")
	}

	return visibility(v), nil
}

// CreateVisibility creates a new Visibility instance based on the given string.
func CreateVisibility(v string) Visibility {
	return visibility(v)
}

type visibility string

// Visibility returns the visibility as a string.
func (r visibility) Visibility() string {
	return string(r)
}

// IsPrivate checks if the visibility is private.
func (r visibility) IsPrivate() bool {
	return string(r) == Private
}

// IsPublic checks if the visibility is public.
func (r visibility) IsPublic() bool {
	return string(r) == Public
}

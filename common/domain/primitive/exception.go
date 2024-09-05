/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"strings"
)

const (
	RelatedModelDisabled = "related_model_disabled"
	RelatedModelNotFound = "related_model_notfound"
	NoApplicationFile    = "no_application_file"
	NoCompQuotaException = "no_comp_quota_exception"
)

var ExceptionMap = map[string]string{
	// RelatedModelDisabled reason
	RelatedModelDisabled: "the related model of space is disabled",
	// RelatedModelNotFound reason
	RelatedModelNotFound: "the related model of space is not found",
	// NoApplicationFile reason
	NoApplicationFile: "space no application file",
	// NoCompQuota reason
	NoCompQuotaException: "space no comp quota",
}

var (
	ExceptionRelatedModelDisabled = exception(RelatedModelDisabled)
	ExceptionRelatedModelNotFound = exception(RelatedModelNotFound)
	ExceptionNoApplicationFile    = exception(NoApplicationFile)
)

// Exception is an interface that defines the exception of an object.
type Exception interface {
	Exception() string
}

// NewException creates a new Exception instance based on the given string.
func NewException(v string) (Exception, error) {
	v = strings.ToLower(v)
	if v != RelatedModelDisabled {
		return nil, errors.New("unknown exception")
	}

	return exception(v), nil
}

// CreateException creates a new Exception instance based on the given string.
func CreateException(v string) Exception {
	return exception(v)
}

type exception string

// Exception returns the exception as a string.
func (r exception) Exception() string {
	return string(r)
}

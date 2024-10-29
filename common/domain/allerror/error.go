/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package allerror provides a set of error codes and error types used in the application.
package allerror

import (
	"errors"
	"strings"
)

const (
	// ErrorCodeComputilityAccountFindError find computility account error
	ErrorCodeComputilityAccountFindError = "computility_account_find_error"

	// ErrorCodeComputilityOrgFindError find computility org error
	ErrorCodeComputilityOrgFindError = "computility_org_find_error"

	// ErrorCodeComputilityOrgUpdateError update computility org error
	ErrorCodeComputilityOrgUpdateError = "computility_org_update_error"

	// ErrorComputilityOrgQuotaLowerBoundError quota count lower bound error
	ErrorComputilityOrgQuotaLowerBoundError = "computility_org_quota_lower_bound_error"

	// ErrorComputilityOrgQuotaMultipleError quota count not a multiple of default quota
	ErrorComputilityOrgQuotaMultipleError = "computility_org_quota_multiple_error"

	// ErrorCodeInsufficientQuota user has insufficient quota balance
	ErrorCodeInsufficientQuota = "insufficient_quota"

	// ErrorCodeNoNpuPermission user has no npu permission
	ErrorCodeNoNpuPermission = "no_npu_permission"

	// ErrorCodeAccessTokenInvalid This error code is for restful api
	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	// ErrorCodeSpaceAppUnmatchedStatus is const
	ErrorCodeSpaceAppUnmatchedStatus = "space_app_unmatched_status"

	// ErrorCodeSpaceAppNotFound space app
	ErrorCodeSpaceAppNotFound = "space_app_not_found"

	// ErrorCodeSpaceNotFound is const
	ErrorCodeSpaceNotFound = "space_not_found"

	// ErrorCodeSpaceCommitConflict is const
	ErrorCodeSpaceCommitConflict = "space_commit_conflict"

	// ErrorCodeSpaceAppCreateFailed
	ErrorCodeSpaceAppCreateFailed = "space_app_create_failed"
)

// errorImpl
type errorImpl struct {
	code     string
	msg      string
	innerErr error // error info for diagnostic
}

// Error returns the error message.
//
// This function returns the error message of the errorImpl struct.
//
// No parameters.
// Returns a string representing the error message.
func (e errorImpl) Error() string {
	return e.msg
}

// ErrorCode returns the error code.
//
// This function returns the error code of the errorImpl struct.
// The error code is a string representing the type of the error, it could be used for error handling and diagnostic.
//
// No parameters.
// Returns a string representing the error code.
func (e errorImpl) ErrorCode() string {
	return e.code
}

// InnerError returns the inner error.
type InnerError interface {
	InnerError() error
}

// InnerErr returns the inner error.
func (e errorImpl) InnerError() error {
	return e.innerErr
}

// InnerErr returns the inner error.
func InnerErr(err error) error {
	var v InnerError
	if ok := errors.As(err, &v); ok {
		return v.InnerError()
	}

	return err
}

// New creates a new error with the specified code and message.
//
// This function creates a new errorImpl struct with the specified code, message and error info
// for diagnostic. If the message is empty, the function will replace all "_" in the code with
// " " as the message.
//
// Parameters:
//
//	code: a string representing the type of the error
//	msg: a string representing the error message, which is returned to client or end user
//	err: error info for diagnostic, which is used for diagnostic by developers
//
// Returns an errorImpl struct.
func New(code string, msg string, err error) errorImpl {
	v := errorImpl{
		code:     code,
		innerErr: err,
	}

	if msg == "" {
		v.msg = strings.ReplaceAll(code, "_", " ")
	} else {
		v.msg = msg
	}

	return v
}

// notfoudError
type notfoudError struct {
	errorImpl
}

// NotFound is a marker method for a not found error.
func (e notfoudError) NotFound() {}

// NewNotFound creates a new not found error with the specified code and message.
func NewNotFound(code string, msg string, err error) notfoudError {
	return notfoudError{errorImpl: New(code, msg, err)}
}

// IsNotFound checks if the given error is a not found error.
func IsNotFound(err error) (notfoudError, bool) {
	if err == nil {
		return notfoudError{}, false
	}
	var v notfoudError
	ok := errors.As(err, &v)

	return v, ok
}

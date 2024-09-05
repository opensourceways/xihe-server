/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides utility functions for handling HTTP errors and error codes.
package controller

import (
	"errors"
	"net/http"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
)

const (
	errorSystemError      = "system_error"
	errorBadRequestBody   = "bad_request_body"
	errorModerationFailed = "moderation_failed"
	errorBadRequestParam  = "bad_request_param"
)

type errorCode interface {
	ErrorCode() string
}

type errorNotFound interface {
	errorCode

	NotFound()
}

type errorNoPermission interface {
	errorCode

	NoPermission()
}

func httpError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	sc := http.StatusInternalServerError
	code := errorSystemError
	var t errorCode

	if ok := errors.As(err, &t); ok {
		code = t.ErrorCode()

		var n errorNotFound
		var p errorNoPermission

		if ok := errors.As(err, &n); ok {
			sc = http.StatusNotFound

		} else if ok := errors.As(err, &p); ok {
			sc = http.StatusForbidden

		} else {
			switch code {
			case allerror.ErrorCodeAccessTokenInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeSessionIdMissing:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeSessionIdInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeSessionNotFound:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeSessionInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenMissing:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenNotFound:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeOrgExistResource:
				sc = http.StatusBadRequest

			case allerror.ErrorRateLimitOver:
				sc = http.StatusTooManyRequests

			default:
				sc = http.StatusBadRequest
			}
		}
	}

	return sc, code
}

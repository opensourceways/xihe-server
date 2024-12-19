package controller

import (
	"errors"
	"net/http"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
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
			default:
				sc = http.StatusBadRequest
			}
		}
	}

	return sc, code
}

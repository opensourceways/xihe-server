package bigmodel

import "errors"

// errorSensitiveInfo
type errorSensitiveInfo struct {
	error
}

func NewErrorSensitiveInfo(err error) errorSensitiveInfo {
	return errorSensitiveInfo{err}
}

// helper
func IsErrorSensitiveInfo(err error) bool {
	return errors.As(err, &errorSensitiveInfo{})

}

// errorBusySource
type errorBusySource struct {
	error
}

func NewErrorBusySource(err error) errorBusySource {
	return errorBusySource{err}
}

func IsErrorBusySource(err error) bool {
	_, ok := err.(errorBusySource)

	return ok
}

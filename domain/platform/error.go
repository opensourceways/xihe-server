package platform

import "errors"

// errorTooManyFilesToDelete
type errorTooManyFilesToDelete struct {
	error
}

func NewErrorTooManyFilesToDelete(err error) errorTooManyFilesToDelete {
	return errorTooManyFilesToDelete{err}
}

// helper
func IsErrorTooManyFilesToDelete(err error) bool {
	return errors.As(err, &errorTooManyFilesToDelete{})
}

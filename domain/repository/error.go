package repository

import "errors"

// ErrorDuplicateCreating
type ErrorDuplicateCreating struct {
	error
}

func NewErrorDuplicateCreating(err error) ErrorDuplicateCreating {
	return ErrorDuplicateCreating{err}
}

// ErrorResourceNotExists
type ErrorResourceNotExists struct {
	error
}

func NewErrorResourceNotExists(err error) ErrorResourceNotExists {
	return ErrorResourceNotExists{err}
}

// ErrorConcurrentUpdating
type ErrorConcurrentUpdating struct {
	error
}

func NewErrorConcurrentUpdating(err error) ErrorConcurrentUpdating {
	return ErrorConcurrentUpdating{err}
}

// helper
func IsErrorResourceNotExists(err error) bool {
	return errors.As(err, &ErrorResourceNotExists{})
}

func IsErrorDuplicateCreating(err error) bool {
	return errors.As(err, &ErrorDuplicateCreating{})
}

func IsErrorConcurrentUpdating(err error) bool {
	return errors.As(err, &ErrorConcurrentUpdating{})
}

type ExcendMaximumPageNumError struct {
	error
}

func NewExcendMaximumPageNumError(err error) ExcendMaximumPageNumError {
	return ExcendMaximumPageNumError{err}
}

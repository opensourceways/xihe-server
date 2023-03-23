package repository

// ErrorDuplicateCreating
type ErrorDuplicateCreating struct {
	error
}

func NewErrorDuplicateCreating(err error) ErrorDuplicateCreating {
	return ErrorDuplicateCreating{err}
}

// ErrorResourceNotFound
type ErrorResourceNotFound struct {
	error
}

func NewErrorResourceNotFound(err error) ErrorResourceNotFound {
	return ErrorResourceNotFound{err}
}

// ErrorConcurrentUpdating
type ErrorConcurrentUpdating struct {
	error
}

func NewErrorConcurrentUpdating(err error) ErrorConcurrentUpdating {
	return ErrorConcurrentUpdating{err}
}

// helper

func IsErrorResourceNotFound(err error) bool {
	_, ok := err.(ErrorResourceNotFound)

	return ok
}

func IsErrorDuplicateCreating(err error) bool {
	_, ok := err.(ErrorDuplicateCreating)

	return ok
}

func IsErrorConcurrentUpdating(err error) bool {
	_, ok := err.(ErrorConcurrentUpdating)

	return ok
}

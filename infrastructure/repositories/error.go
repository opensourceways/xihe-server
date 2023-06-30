package repositories

import "github.com/opensourceways/xihe-server/domain/repository"

// errorDuplicateCreating
type errorDuplicateCreating struct {
	error
}

func NewErrorDuplicateCreating(err error) errorDuplicateCreating {
	return errorDuplicateCreating{err}
}

// errorDataNotExists
type errorDataNotExists struct {
	error
}

func NewErrorDataNotExists(err error) errorDataNotExists {
	return errorDataNotExists{err}
}

func isErrorDataNotExists(err error) bool {
	_, ok := err.(errorDataNotExists)

	return ok
}

// errorConcurrentUpdating
type errorConcurrentUpdating struct {
	error
}

func NewErrorConcurrentUpdating(err error) errorConcurrentUpdating {
	return errorConcurrentUpdating{err}
}

// convertError
func ConvertError(err error) (out error) {
	return convertError(err)
}

func convertError(err error) (out error) {
	switch err.(type) {
	case errorDuplicateCreating:
		out = repository.NewErrorDuplicateCreating(err)

	case errorDataNotExists:
		out = repository.NewErrorResourceNotExists(err)

	case errorConcurrentUpdating:
		out = repository.NewErrorConcurrentUpdating(err)

	default:
		out = err
	}

	return
}

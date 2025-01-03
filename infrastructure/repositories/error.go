package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain/repository"
)

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
	return errors.As(err, &errorDataNotExists{})
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

func convertError(err error) error {
	if errors.As(err, &errorDuplicateCreating{}) {
		return repository.NewErrorDuplicateCreating(err)
	}
	if errors.As(err, &errorDataNotExists{}) {
		return repository.NewErrorResourceNotExists(err)
	}
	if errors.As(err, &errorConcurrentUpdating{}) {
		return repository.NewErrorConcurrentUpdating(err)
	}

	return err
}

package bigmodels

type errorConcurrentRequest struct {
	error
}

func NewErrorConcurrentRequest(err error) errorConcurrentRequest {
	return errorConcurrentRequest{err}
}

func IsErrorConcurrentRequest(err error) bool {
	_, ok := err.(errorConcurrentRequest)

	return ok
}

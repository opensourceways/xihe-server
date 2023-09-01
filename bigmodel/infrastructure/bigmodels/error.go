package bigmodels

const (
	CodeInputTextAuditError     = "code_input_text_audit_error"
	CodeOutputTextAuditError    = "code_output_text_audit_error"
	CodeBaiChuanGenerationError = "code_baichuan_generation_error"
)

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

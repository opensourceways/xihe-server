package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/repository"
)

const (
	errorNotAllowed          = "not_allowed"
	errorInvalidToken        = "invalid_token"
	errorSystemError         = "system_error"
	errorBadRequestBody      = "bad_request_body"
	errorBadRequestHeader    = "bad_request_header"
	errorBadRequestParam     = "bad_request_param"
	errorDuplicateCreating   = "duplicate_creating"
	errorResourceNotExists   = "resource_not_exists"
	errorConcurrentUpdating  = "concurrent_updateing"
	errorExccedMaxNum        = "exceed_max_num"
	errorUpdateLFSFile       = "update_lfs_file"
	errorPreviewLFSFile      = "preview_lfs_file"
	errorUnavailableRepoFile = "unavailable_repo_file"
)

var (
	respBadRequestBody = newResponseCodeMsg(
		errorBadRequestBody, "can't fetch request body",
	)
)

// responseData is the response data to client
type responseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func isErrorOfAccessingPrivateRepo(err error) bool {
	_, ok := err.(app.ErrorPrivateRepo)

	return ok
}

func newResponseError(err error) responseData {
	code := errorSystemError

	switch err.(type) {
	case repository.ErrorDuplicateCreating:
		code = errorDuplicateCreating

	case repository.ErrorResourceNotExists:
		code = errorResourceNotExists

	case repository.ErrorConcurrentUpdating:
		code = errorConcurrentUpdating

	case app.ErrorExceedMaxRelatedResourceNum:
		code = errorExccedMaxNum

	case app.ErrorUpdateLFSFile:
		code = errorUpdateLFSFile

	case app.ErrorUnavailableRepoFile:
		code = errorUnavailableRepoFile

	case app.ErrorPreviewLFSFile:
		code = errorPreviewLFSFile
	
	default:
		
	}

	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseData(data interface{}) responseData {
	return responseData{
		Data: data,
	}
}

func newResponseCodeError(code string, err error) responseData {
	return responseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseCodeMsg(code, msg string) responseData {
	return responseData{
		Code: code,
		Msg:  msg,
	}
}

func respBadRequestParam(err error) responseData {
	return newResponseCodeError(
		errorBadRequestParam, err,
	)
}

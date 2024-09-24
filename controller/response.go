package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/repository"
)

const (
	errorNotAllowed            = "not_allowed"
	errorInvalidToken          = "invalid_token"
	errorSystemError           = "system_error"
	errorBadRequestBody        = "bad_request_body"
	errorBadRequestHeader      = "bad_request_header"
	errorBadRequestParam       = "bad_request_param"
	errorDuplicateCreating     = "duplicate_creating"
	errorResourceNotExists     = "resource_not_exists"
	errorConcurrentUpdating    = "concurrent_updating"
	errorExccedMaxNum          = "exceed_max_num"
	errorUpdateLFSFile         = "update_lfs_file"
	errorPreviewLFSFile        = "preview_lfs_file"
	errorUnavailableRepoFile   = "unavailable_repo_file"
	errorDuplicateTrainingName = "duplicate_training_name"
	errorExccedMaxiumPageNum   = "excend_maximum_page_num"
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
	return errors.As(err, &app.ErrorPrivateRepo{})
}

func newResponseError(err error) responseData {
	code := errorSystemError

	if errors.As(err, &repository.ErrorDuplicateCreating{}) {
		code = errorDuplicateCreating
	} else if errors.As(err, &repository.ErrorResourceNotExists{}) {
		code = errorResourceNotExists
	} else if errors.As(err, &repository.ErrorConcurrentUpdating{}) {
		code = errorConcurrentUpdating
	} else if errors.As(err, &app.ErrorExceedMaxRelatedResourceNum{}) {
		code = errorExccedMaxNum
	} else if errors.As(err, &app.ErrorUpdateLFSFile{}) {
		code = errorUpdateLFSFile
	} else if errors.As(err, &app.ErrorUnavailableRepoFile{}) {
		code = errorUnavailableRepoFile
	} else if errors.As(err, &app.ErrorPreviewLFSFile{}) {
		code = errorPreviewLFSFile
	} else if errors.As(err, &app.ErrorDuplicateTrainingName{}) {
		code = errorDuplicateTrainingName
	} else if errors.As(err, &repository.ExcendMaxiumPageNumError{}) {
		code = errorExccedMaxiumPageNum
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

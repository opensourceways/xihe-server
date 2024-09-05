/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package controller provides utility functions for handling HTTP errors and error codes.
package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
)

// ResponseData is a struct that holds the response data for an API request.
type ResponseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newResponseData(data interface{}) ResponseData {
	return ResponseData{
		Data: data,
	}
}

// nolint:golint,unused
func newResponseCodeError(code string, err error) ResponseData {
	return ResponseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseCodeMsg(code, msg string) ResponseData {
	return ResponseData{
		Code: code,
		Msg:  msg,
	}
}

// SendBadRequestBody sends a bad request body error response.
func SendBadRequestBody(ctx *gin.Context, err error) {
	if _, ok := err.(errorCode); ok {
		SendError(ctx, err)
	} else {
		_ = ctx.Error(err)
		resp := newResponseCodeMsg(errorBadRequestBody, err.Error())
		if errors.As(err, new(validator.ValidationErrors)) {
			resp = newResponseCodeMsg(errorModerationFailed, "moderation failed")
		}
		ctx.JSON(http.StatusBadRequest, resp)
	}
}

// SendBadRequestParam sends a bad request parameter error response.
func SendBadRequestParam(ctx *gin.Context, err error) {
	if _, ok := err.(errorCode); ok {
		SendError(ctx, err)
	} else {
		_ = ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest,
			newResponseCodeMsg(errorBadRequestParam, err.Error()),
		)
	}
}

// SendRespOfPut sends a successful PUT response with data if provided.
func SendRespOfPut(ctx *gin.Context, data interface{}) {
	if data == nil {
		ctx.JSON(http.StatusAccepted, newResponseCodeMsg("", "success"))
	} else {
		ctx.JSON(http.StatusAccepted, newResponseData(data))
	}
}

// SendRespOfGet sends a successful GET response with data.
func SendRespOfGet(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, newResponseData(data))
}

// SendRespOfPost sends a successful POST response with data if provided.
func SendRespOfPost(ctx *gin.Context, data interface{}) {
	if data == nil {
		ctx.JSON(http.StatusCreated, newResponseCodeMsg("", "success"))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(data))
	}
}

// SendRespOfDelete sends a successful DELETE response.
func SendRespOfDelete(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, newResponseCodeMsg("", "success"))
}

// SendError sends an error response based on the given error.
func SendError(ctx *gin.Context, err error) {
	sc, code := httpError(err)

	_ = ctx.AbortWithError(sc, allerror.InnerErr(err))

	ctx.JSON(sc, newResponseCodeMsg(code, err.Error()))
}

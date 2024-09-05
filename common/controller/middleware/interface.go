/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package middleware provides a set of middleware functions for Gin framework
// that handle user authentication and authorization.
package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
)

// UserMiddleWare is an interface that defines methods for user authentication and authorization.
type UserMiddleWare interface {
	// 1. The token must be exist and valid and has write role,
	// otherwise abort directly and send allerror.ErrorCodeAccessTokenInvalid
	// 2. If token is valid, then parse the user bound to the token and save it
	Write(*gin.Context)

	// 1. The token must be exist and valid and has read role,
	// otherwise abort directly and send allerror.ErrorCodeAccessTokenInvalid
	// 2. If token is valid, then parse the user bound to the token and save it
	Read(*gin.Context)

	// Get user parsed from the token.
	GetUser(*gin.Context) primitive.Account

	// Get user parsed from the token and send error if failed.
	GetUserAndExitIfFailed(ctx *gin.Context) primitive.Account

	// 1. If token is not passed, ignore it.
	// 2. If token exists, it must be valid, otherwise abort directly and
	// send allerror.ErrorCodeAccessTokenInvalid
	// 3. If token is valid, then parse the user bound to the token and save it
	Optional(*gin.Context)
}

// TokenMiddleWare is an interface that defines methods for user authentication and authorization.
type TokenMiddleWare interface {
	// 1. The session must be exist and valid and has read role,
	// otherwise abort directly and send allerror.ErrorCodeAccessTokenInvalid
	// 2. If session is valid, then parse the user bound to the token and save it
	CheckSession(*gin.Context)
}

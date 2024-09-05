/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package internalservice provides functionality for internal services.
package internalservice

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	commonctl "github.com/opensourceways/xihe-server/common/controller"
	"github.com/opensourceways/xihe-server/common/controller/middleware"
	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/common/domain/primitive"
)

const tokenHeader = "TOKEN" // #nosec G101

var noUserError = errors.New("no user")

// NewAPIMiddleware creates a new instance of internalServiceAPIMiddleware.
func NewAPIMiddleware(securityLog middleware.SecurityLog) *internalServiceAPIMiddleware {
	return &internalServiceAPIMiddleware{
		securityLog: securityLog,
	}
}

// internalServiceAPIMiddleware
type internalServiceAPIMiddleware struct {
	securityLog middleware.SecurityLog
}

// Write method for internalServiceAPIMiddleware.
func (m *internalServiceAPIMiddleware) Write(ctx *gin.Context) {
	m.must(ctx)
}

// Read method for internalServiceAPIMiddleware.
func (m *internalServiceAPIMiddleware) Read(ctx *gin.Context) {
	m.must(ctx)
}

// Optional method for internalServiceAPIMiddleware.
func (m *internalServiceAPIMiddleware) Optional(ctx *gin.Context) {
	if v := ctx.GetHeader(tokenHeader); v == "" {
		ctx.Next()
	} else {
		m.must(ctx)
	}
}

func (m *internalServiceAPIMiddleware) must(ctx *gin.Context) {
	if err := m.checkToken(ctx); err != nil {
		commonctl.SendError(ctx, err)
		m.securityLog.Warn(ctx, err.Error())

		ctx.Abort()
	} else {
		ctx.Next()
	}
}

// GetUser method for internalServiceAPIMiddleware.
func (m *internalServiceAPIMiddleware) GetUser(ctx *gin.Context) primitive.Account {
	return nil
}

// GetUserAndExitIfFailed method for internalServiceAPIMiddleware.
func (m *internalServiceAPIMiddleware) GetUserAndExitIfFailed(ctx *gin.Context) primitive.Account {
	commonctl.SendError(ctx, noUserError)

	return nil
}

func (m *internalServiceAPIMiddleware) checkToken(ctx *gin.Context) error {
	rawToken := ctx.GetHeader(tokenHeader)
	calcTokenHash, err := commonctl.EncodeToken(rawToken, config.Salt)
	if err != nil {
		return allerror.New(
			allerror.ErrorCodeAccessTokenInvalid, "check token failed", err,
		)
	}

	if calcTokenHash != config.TokenHash {
		return allerror.New(
			allerror.ErrorCodeAccessTokenInvalid, "invalid token", fmt.Errorf("token mismatch"),
		)
	}

	return nil
}

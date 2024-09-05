/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package middleware provides a set of middleware functions for Gin framework
package middleware

import "github.com/gin-gonic/gin"

const action = "OPERATION_LOG_ACTION"

// OperationLog is an interface that defines the Write method for writing operation logs.
type OperationLog interface {
	Write(*gin.Context)
}

// SetAction sets the action value in the given context.
func SetAction(ctx *gin.Context, v string) {
	ctx.Set(action, v)
}

// GetAction retrieves the action value from the given context.
func GetAction(ctx *gin.Context) string {
	v, ok := ctx.Get(action)
	if !ok {
		return ""
	}

	if str, ok := v.(string); ok {
		return str
	}

	return ""
}

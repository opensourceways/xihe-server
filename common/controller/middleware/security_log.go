/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package middleware provides a set of middleware functions for Gin framework
package middleware

import "github.com/gin-gonic/gin"

// SecurityLog is an interface that defines the Info and Warn methods for logging security-related information.
type SecurityLog interface {
	Info(*gin.Context, ...interface{})
	Warn(*gin.Context, ...interface{})
}

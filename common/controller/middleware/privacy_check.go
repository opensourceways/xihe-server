/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package middleware provides a set of middleware functions for Gin framework
package middleware

import "github.com/gin-gonic/gin"

// PrivacyCheck is an interface that defines the methods for checking user privacy agreement.
type PrivacyCheck interface {
	Check(*gin.Context)
	CheckOwner(ctx *gin.Context)
	CheckName(ctx *gin.Context)
}

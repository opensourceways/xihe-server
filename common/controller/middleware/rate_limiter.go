/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package middleware is limiter for ensure safe
package middleware

import "github.com/gin-gonic/gin"

// RateLimiter is an interface that defines the CheckLimit method for CheckLimit safe.
type RateLimiter interface {
	CheckLimit(*gin.Context)
}

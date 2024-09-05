/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package securitylog provides functionality for logging security-related information.
package securitylog

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	prefix = "MF_SECURITY_LOG"
)

// SecurityLog creates a new instance of the securityLog struct.
func SecurityLog() *securityLog {
	return &securityLog{}
}

type securityLog struct {
}

// Info logs an informational security message.
func (log *securityLog) Info(ctx *gin.Context, msg ...interface{}) {

	temp := fmt.Sprintf("%v | Operation record", prefix)

	clientIp := fmt.Sprintf(" | client ip: [%s]", ctx.ClientIP())

	requestUrl := fmt.Sprintf(" | request url: [%s]", ctx.Request.URL.String())

	method := fmt.Sprintf(" | method: [%s]", ctx.Request.Method)

	state := fmt.Sprintf(" | state: [%d]", ctx.Writer.Status())

	message := fmt.Sprintf(" | message: [%v]", fmt.Sprint(msg...))

	logrus.Info(temp, clientIp, requestUrl, method, state, message)

}

// Warn logs a warning message for intercepted illegal requests.
func (log *securityLog) Warn(ctx *gin.Context, msg ...interface{}) {

	temp := fmt.Sprintf("%v | Illegal requests are intercepted", prefix)

	clientIp := fmt.Sprintf(" | client ip: [%s]", ctx.ClientIP())

	requestUrl := fmt.Sprintf(" | request url: [%s]", ctx.Request.URL.String())

	method := fmt.Sprintf(" | method: [%s]", ctx.Request.Method)

	state := fmt.Sprintf(" | state: [%d]", ctx.Writer.Status())

	message := fmt.Sprintf(" | message: [%v]", fmt.Sprint(msg...))

	logrus.Warn(temp, clientIp, requestUrl, method, state, message)

}

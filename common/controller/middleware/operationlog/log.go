/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package operationlog provides functionality for logging operation-related information.
package operationlog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/common/controller"
	"github.com/opensourceways/xihe-server/common/controller/middleware"
	"github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/utils"
)

// OperationLog creates a new instance of the operationLog struct.
func OperationLog(u middleware.UserMiddleWare) *operationLog {
	return &operationLog{user: u}
}

type operationLog struct {
	user middleware.UserMiddleWare
}

// Write logs the operation details to the log file.
func (log *operationLog) Write(ctx *gin.Context) {
	record := domain.OperationLogRecord{}

	record.Time = utils.Time()

	ctx.Next()

	record.Action = middleware.GetAction(ctx)
	if record.Action == "" {
		// It is meaningless to record operation log if action is missing.
		return
	}

	record.User = log.user.GetUser(ctx)

	record.IP, _ = controller.GetIp(ctx)

	v := ctx.Writer.Status()
	record.Success = v >= http.StatusOK && v < http.StatusMultipleChoices

	record.Method = ctx.Request.Method

	logrus.Info(record.String())
}

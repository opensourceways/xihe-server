/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and types.
package domain

import (
	"fmt"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
)

const prefix = "MF_OPERATION_LOG"

// OperationLogRecord represents a operation log record.
type OperationLogRecord struct {
	IP      string
	User    primitive.Account
	Time    string
	Method  string
	Action  string
	Success bool
}

// String returns the description of record.
func (r *OperationLogRecord) String() string {
	user := ""
	if r.User != nil {
		user = r.User.Account()
	}

	result := "success"
	if !r.Success {
		result = "failed"
	}

	return fmt.Sprintf(
		"%s | %s | %s | %s | %s | %v | %s",
		prefix,
		r.Time,
		user,
		r.IP,
		r.Method,
		r.Action,
		result,
	)
}

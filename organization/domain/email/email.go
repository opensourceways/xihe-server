/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package email provides functionality for sending emails.
package email

import (
	"github.com/opensourceways/xihe-server/organization/domain"
	"github.com/opensourceways/xihe-server/organization/domain/primitive"
)

// Email is an interface for sending organization certificates.
type Email interface {
	Send(domain.OrgCertificate, primitive.Image) error
}

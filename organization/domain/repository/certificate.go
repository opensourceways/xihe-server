/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repository

import (
	"context"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain"
	orgprimitive "github.com/opensourceways/xihe-server/organization/domain/primitive"
)

// FindOption represents the options for finding an organization certificate.
type FindOption struct {
	Phone                   primitive.Phone
	OrgName                 primitive.Account
	CertificateOrgName      primitive.AccountFullname
	UnifiedSocialCreditCode orgprimitive.USCC
	OrgType                 orgprimitive.CertificateOrgType
}

// Certificate represents the repository for organization certificates.
type Certificate interface {
	Save(domain.OrgCertificate) error
	Find(context.Context, FindOption) (domain.OrgCertificate, error)
	GetAllCount() (int64, error)
	FinAllName(context.Context) ([]string, error)
	FindList(context.Context, int, int, orgprimitive.CertificateOrgType) ([]domain.OrgCertificate, int, error)
	DuplicateCheck(ctx context.Context, option FindOption) (domain.OrgCertificate, error)
	DeleteByOrgName(orgName primitive.Account) error
}

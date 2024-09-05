/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for managing approvals in an organization.
package repository

import (
	"context"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain"
)

// OrgMember interface defines the methods for managing organization members.
type OrgMember interface {
	Add(*domain.OrgMember) (domain.OrgMember, error)
	Save(*domain.OrgMember) (domain.OrgMember, error)
	Delete(context.Context, *domain.OrgMember) error
	DeleteByOrg(primitive.Account) error
	GetByOrg(*domain.OrgListMemberCmd) ([]domain.OrgMember, error)
	GetByOrgAndRole(string, primitive.Role) ([]domain.OrgMember, error)
	GetByOrgAndUser(context.Context, string, string) (domain.OrgMember, error)
	GetByUser(string) ([]domain.OrgMember, error)
	GetByUserAndRoles(primitive.Account, []primitive.Role) ([]domain.OrgMember, error)
}

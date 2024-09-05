/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repository provides interfaces for managing approvals in an organization.
package repository

import (
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain"
)

// Approve is an interface that defines the methods for handling approval-related operations.
type Approve interface {
	AddInvite(*domain.Approve) (domain.Approve, bool, error)
	SaveInvite(*domain.Approve) (domain.Approve, error)
	AddRequest(*domain.MemberRequest) (domain.MemberRequest, bool, error)
	SaveRequest(*domain.MemberRequest) (domain.MemberRequest, error)
	DeleteInviteAndReqByOrg(primitive.Account) error
	Count(primitive.Account) (int64, error)
	GetOneApply(string, string) ([]domain.MemberRequest, error)
	GetInvite(string, string) (domain.MemberRequest, error)
	ListInvitation(*domain.OrgInvitationListCmd) ([]domain.Approve, error)
	ListRequests(*domain.OrgMemberReqListCmd) ([]domain.MemberRequest, error)
	ListPagnation(*domain.OrgMemberReqListCmd) ([]domain.MemberRequest, int, error)
	UpdateAllApproveStatus(primitive.Account, primitive.Account, domain.ApproveStatus) error
}

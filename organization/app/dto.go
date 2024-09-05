/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package app provides functionality for handling organization-related operations.
package app

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain"
	orgprimitive "github.com/opensourceways/xihe-server/organization/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain/repository"
	userapp "github.com/opensourceways/xihe-server/user/app"
)

// OrganizationDTO represents the data transfer object for an organization.
type OrganizationDTO struct {
	Id           string `json:"id"`
	Name         string `json:"account"`
	FullName     string `json:"fullname"`
	AvatarId     string `json:"avatar_id"`
	Description  string `json:"description"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
	Website      string `json:"website"`
	Owner        string `json:"owner"`
	OwnerId      string `json:"owner_id"`
	DefaultRole  string `json:"default_role"`
	AllowRequest bool   `json:"allow_request"`
}

type ResBaseORgInfoDTO struct {
	BaseOrgInfo []domain.OrgBaseInfo `json:"org_info"`
	Total       int                  `json:"total"`
}

// ApproveDTO represents the data transfer object for an approval request.
type ApproveDTO struct {
	Id          string `json:"id"`
	OrgName     string `json:"org_name"`
	OrgFullName string `json:"org_full_name"`
	OrgId       string `json:"org_id"`
	UserName    string `json:"user_name"`
	UserId      string `json:"user_id"`
	Role        string `json:"role"`
	ExpiresAt   int64  `json:"expires_at"`
	Fullname    string `json:"fullname"`
	Inviter     string `json:"inviter"`
	Status      string `json:"status"`
	By          string `json:"by"`
	Msg         string `json:"msg"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// ToApproveDTO converts a domain.Approve object to an ApproveDTO object.
func ToApproveDTO(ctx context.Context, m *domain.Approve, user userapp.UserService) ApproveDTO {
	var fullname string
	u, err := user.GetByAccount(ctx, nil, m.Username)
	if err != nil {
		logrus.Warnf("failed to get fullname for %s, err:%s", m.Username, err)
	} else {
		fullname = u.Fullname
	}

	var orgFullName string
	orgFullName, err = user.GetUserFullname(ctx, m.OrgName)
	if err != nil {
		logrus.Warnf("failed to get org fullname for %s, err:%s", m.OrgName, err)
	}

	return ApproveDTO{
		Id:          m.Id.Identity(),
		OrgName:     m.OrgName.Account(),
		OrgFullName: orgFullName,
		OrgId:       m.OrgId.Identity(),
		UserName:    m.Username.Account(),
		UserId:      m.UserId.Identity(),
		Role:        m.Role.Role(),
		ExpiresAt:   m.ExpireAt, // will expire in 14 days
		Inviter:     m.Inviter.Account(),
		Status:      string(m.Status),
		Fullname:    fullname,
		Msg:         m.Msg,
		By:          m.By,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// MemberRequestDTO represents the data transfer object for a member request.
type MemberRequestDTO struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	UserId    string `json:"user_id"`
	Fullname  string `json:"fullname"`
	Role      string `json:"role"`
	OrgName   string `json:"org_name"`
	OrgId     string `json:"org_id"`
	Status    string `json:"status"`
	By        string `json:"by"`
	Msg       string `json:"msg"`
	AvatarId  string `json:"avatar_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type MemberPagnationDTO struct {
	Total   int                `json:"total"`
	Members []MemberRequestDTO `json:"members"`
}

// ToMemberRequestDTO converts a domain.MemberRequest object to a MemberRequestDTO object.
func ToMemberRequestDTO(ctx context.Context, m *domain.MemberRequest, user userapp.UserService) MemberRequestDTO {
	var fullname string
	u, err := user.GetByAccount(ctx, nil, m.Username)
	if err != nil {
		logrus.Warnf("failed to get fullname for %s, err:%s", m.Username, err)
	} else {
		fullname = u.Fullname
	}

	return MemberRequestDTO{
		Id:        m.Id.Identity(),
		Username:  m.Username.Account(),
		UserId:    m.UserId.Identity(),
		Role:      m.Role.Role(),
		OrgName:   m.OrgName.Account(),
		OrgId:     m.OrgId.Identity(),
		Status:    string(m.Status),
		Fullname:  fullname,
		By:        m.By,
		Msg:       m.Msg,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// MemberDTO represents the data transfer object for a member of an organization.
type MemberDTO struct {
	Id          string `json:"id"`
	OrgName     string `json:"org_name"`
	OrgId       string `json:"org_id"`
	OrgFullName string `json:"org_fullname"`
	UserName    string `json:"user_name"`
	FullName    string `json:"full_name"`
	UserId      string `json:"user_id"`
	Role        string `json:"role"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	AvatarId    string `json:"avatar_id"`
}

// OrgListOptions represents the options for listing organizations.
type OrgListOptions struct {
	Page     int
	PageSize int
	Owner    primitive.Account
	Member   primitive.Account
	Roles    []primitive.Role
	Search   string
	OrgType  orgprimitive.CertificateOrgType
}

type OrgListNewOptions struct {
	OrgListOptions
	Search  string
	OrgType orgprimitive.CertificateOrgType
}

// ToDTO converts a domain.Organization object to a userapp.UserDTO object.
func ToDTO(org *domain.Organization) userapp.UserDTO {
	return userapp.NewUserDTO(org, nil)
}

func ToBaseDTO(org *domain.OrgBaseInfo) userapp.BaseUserDTO {
	return userapp.NewBaseDTO(org, nil)
}

// ToMemberDTO converts a domain.OrgMember object to a MemberDTO object.
func ToMemberDTO(member *domain.OrgMember) MemberDTO {
	return MemberDTO{
		Id:        member.Id.Identity(),
		UserName:  member.Username.Account(),
		FullName:  member.FullName.AccountFullname(),
		UserId:    member.UserId.Identity(),
		Role:      member.Role.Role(),
		OrgName:   member.OrgName.Account(),
		OrgId:     member.OrgId.Identity(),
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
		AvatarId:  member.AvatarId,
	}
}

// OrgCertificateCmd represent a command to create certificate
type OrgCertificateCmd struct {
	domain.OrgCertificate
	Actor              primitive.Account
	ImageOfCertificate orgprimitive.Image
}

// OrgCertificateDuplicateCheckCmd represent repository.FindOption
type OrgCertificateDuplicateCheckCmd = repository.FindOption

// OrgCertificateDTO represent certificate information
type OrgCertificateDTO struct {
	Status                  string `json:"status"`
	Reason                  string `json:"reason"`
	CertificateOrgType      string `json:"certificate_org_type"`
	CertificateOrgName      string `json:"certificate_org_name"`
	UnifiedSocialCreditCode string `json:"unified_social_credit_code"`
	Phone                   string `json:"phone"`
}

// Masked mask sensitive information
func (d *OrgCertificateDTO) Masked() {
	d.UnifiedSocialCreditCode = ""
	d.Phone = ""
}

func toCertificationDTO(cert domain.OrgCertificate) OrgCertificateDTO {
	return OrgCertificateDTO{
		Status:                  cert.Status.CertificateStatus(),
		Reason:                  cert.Reason,
		CertificateOrgType:      cert.CertificateOrgType.CertificateOrgType(),
		CertificateOrgName:      cert.CertificateOrgName.AccountFullname(),
		UnifiedSocialCreditCode: cert.UnifiedSocialCreditCode.USCC(),
		Phone:                   cert.Phone.PhoneNumber(),
	}
}

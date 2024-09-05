/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repositoryimpl provides implementations of repository interfaces for the organization domain.
package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/common/domain/crypto"
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain"
	orgprimitive "github.com/opensourceways/xihe-server/organization/domain/primitive"
)

func toMemberDoc(o *domain.OrgMember) Member {
	do := Member{
		Username: o.Username.Account(),
		FullName: o.FullName.AccountFullname(),
		UserId:   o.UserId.Integer(),
		Orgname:  o.OrgName.Account(),
		OrgId:    o.OrgId.Integer(),
		Role:     o.Role.Role(),
		Version:  o.Version,
		Type:     o.Type,
	}

	do.ID = o.Id.Integer()
	return do
}

func toOrgMember(doc *Member) domain.OrgMember {
	return domain.OrgMember{
		Id:        primitive.CreateIdentity(doc.ID),
		OrgName:   primitive.CreateAccount(doc.Orgname),
		OrgId:     primitive.CreateIdentity(doc.OrgId),
		Role:      primitive.CreateRole(doc.Role),
		Username:  primitive.CreateAccount(doc.Username),
		FullName:  primitive.CreateAccountFullname(doc.FullName),
		UserId:    primitive.CreateIdentity(doc.UserId),
		UpdatedAt: doc.CreatedAt.Unix(),
		CreatedAt: doc.CreatedAt.Unix(),
		Type:      doc.Type,
		Version:   doc.Version,
	}
}

func toApproveDoc(o *domain.Approve) Approve {
	do := Approve{
		Username: o.Username.Account(),
		UserId:   o.UserId.Integer(),
		Orgname:  o.OrgName.Account(),
		OrgId:    o.OrgId.Integer(),
		Role:     o.Role.Role(),
		Expire:   o.ExpireAt,
		Inviter:  o.Inviter.Account(),
		Status:   o.Status,
		By:       o.By,
		Msg:      o.Msg,
		Version:  o.Version,
		Type:     domain.InviteTypeInvite,
	}

	do.ID = o.Id.Integer()
	return do
}

func toRequestDoc(o *domain.MemberRequest) Approve {
	do := Approve{
		Username: o.Username.Account(),
		UserId:   o.UserId.Integer(),
		Orgname:  o.OrgName.Account(),
		OrgId:    o.OrgId.Integer(),
		Role:     o.Role.Role(),
		Status:   o.Status,
		By:       o.By,
		Msg:      o.Msg,
		Version:  o.Version,
		Type:     domain.InviteTypeRequest,
	}

	do.ID = o.Id.Integer()
	return do
}

func toApprove(doc *Approve) domain.Approve {
	return domain.Approve{
		Id:        primitive.CreateIdentity(doc.ID),
		Username:  primitive.CreateAccount(doc.Username),
		UserId:    primitive.CreateIdentity(doc.UserId),
		OrgName:   primitive.CreateAccount(doc.Orgname),
		OrgId:     primitive.CreateIdentity(doc.OrgId),
		Role:      primitive.CreateRole(doc.Role),
		ExpireAt:  doc.Expire,
		Inviter:   primitive.CreateAccount(doc.Inviter),
		Version:   doc.Version,
		By:        doc.By,
		Status:    domain.ApproveStatus(doc.Status),
		Msg:       doc.Msg,
		CreatedAt: doc.CreatedAt.Unix(),
		UpdatedAt: doc.UpdatedAt.Unix(),
	}
}

func toMemberRequest(doc *Approve) domain.MemberRequest {
	return domain.MemberRequest{
		Id:        primitive.CreateIdentity(doc.ID),
		OrgName:   primitive.CreateAccount(doc.Orgname),
		OrgId:     primitive.CreateIdentity(doc.OrgId),
		Username:  primitive.CreateAccount(doc.Username),
		UserId:    primitive.CreateIdentity(doc.UserId),
		Role:      primitive.CreateRole(doc.Role),
		Version:   doc.Version,
		By:        doc.By,
		Status:    domain.ApproveStatus(doc.Status),
		CreatedAt: doc.CreatedAt.Unix(),
		UpdatedAt: doc.UpdatedAt.Unix(),
		Msg:       doc.Msg,
	}
}

func toCertificateDo(cert domain.OrgCertificate, e crypto.Encrypter) (CertificateDO, error) {
	encryptPhone, err := e.Encrypt(cert.Phone.PhoneNumber())
	if err != nil {
		return CertificateDO{}, err
	}

	return CertificateDO{
		OrgName:            cert.OrgName.Account(),
		CertificateOrgName: cert.CertificateOrgName.AccountFullname(),
		CertificateOrgType: cert.CertificateOrgType.CertificateOrgType(),
		USCC:               cert.UnifiedSocialCreditCode.USCC(),
		Status:             cert.Status.CertificateStatus(),
		Phone:              encryptPhone,
		Identity:           cert.Identity.Identity(),
	}, nil
}

func (do CertificateDO) toCertificate(e crypto.Encrypter) (domain.OrgCertificate, error) {
	phone, err := e.Decrypt(do.Phone)
	if err != nil {
		return domain.OrgCertificate{}, err
	}

	return domain.OrgCertificate{
		Status:                  orgprimitive.CreateCertificateStatus(do.Status),
		Reason:                  do.Reason,
		CertificateOrgType:      orgprimitive.CreateCertificateOrgType(do.CertificateOrgType),
		CertificateOrgName:      primitive.CreateAccountFullname(do.CertificateOrgName),
		UnifiedSocialCreditCode: orgprimitive.CreateUSCC(do.USCC),
		Phone:                   primitive.CreatePhoneNumber(phone),
		Identity:                orgprimitive.CreateIdentity(do.Identity),
		OrgName:                 primitive.CreateAccount(do.OrgName),
	}, nil
}

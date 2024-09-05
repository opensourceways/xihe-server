/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"context"
	"errors"

	"github.com/opensourceways/xihe-server/common/domain/primitive"
	commonrepository "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/organization/domain/email"
	"github.com/opensourceways/xihe-server/organization/domain/repository"
)

// OrgCertificateService is the interface definition for organization certificate service.
type OrgCertificateService interface {
	Certificate(context.Context, *OrgCertificateCmd) error
	GetCertification(ctx context.Context, orgName, actor primitive.Account) (OrgCertificateDTO, error)
	DuplicateCheck(ctx context.Context, cmd OrgCertificateDuplicateCheckCmd) (bool, error)
}

// NewOrgCertificateService creates a new instance of the organization certificate service.
func NewOrgCertificateService(
	perm *PermService,
	email email.Email,
	cert repository.Certificate,
) OrgCertificateService {
	return &orgCertificateService{
		perm:        perm,
		email:       email,
		certificate: cert,
	}
}

type orgCertificateService struct {
	perm        *PermService
	email       email.Email
	certificate repository.Certificate
}

// Certificate is a method of the orgCertificateService that handles the certificate operation.
func (org *orgCertificateService) Certificate(ctx context.Context, cmd *OrgCertificateCmd) error {
	err := org.perm.Check(ctx, cmd.Actor, cmd.OrgName, primitive.ObjTypeOrg, primitive.ActionWrite)
	if err != nil {
		return err
	}

	option := repository.FindOption{
		Phone:                   cmd.Phone,
		CertificateOrgName:      cmd.CertificateOrgName,
		UnifiedSocialCreditCode: cmd.UnifiedSocialCreditCode,
	}

	_, err = org.certificate.DuplicateCheck(ctx, option)
	if err == nil {
		return errors.New("duplicate information")
	}
	if !commonrepository.IsErrorResourceNotExists(err) {
		return err
	}

	certificateData := cmd.OrgCertificate
	certificateData.SetProcessingStatus()

	if err = org.certificate.Save(certificateData); err != nil {
		return err
	}

	return org.email.Send(cmd.OrgCertificate, cmd.ImageOfCertificate)
}

// GetCertification is a method of the orgCertificateService
// that retrieves the certificate information for an organization.
func (org *orgCertificateService) GetCertification(
	ctx context.Context, orgName, actor primitive.Account) (OrgCertificateDTO, error) {
	cert, err := org.certificate.Find(ctx, repository.FindOption{OrgName: orgName})
	if err != nil {
		if commonrepository.IsErrorResourceNotExists(err) {
			err = nil
		}

		return OrgCertificateDTO{}, err
	}

	isAdmin := false
	if actor != nil {
		err = org.perm.Check(ctx, actor, orgName, primitive.ObjTypeOrg, primitive.ActionWrite)
		if err == nil {
			isAdmin = true
		}
	}

	dto := toCertificationDTO(cert)
	if !isAdmin {
		dto.Masked()
	}

	return dto, nil
}

// DuplicateCheck is a method of the orgCertificateService
// that performs a duplicate check for organization certificates.
func (org *orgCertificateService) DuplicateCheck(
	ctx context.Context, cmd OrgCertificateDuplicateCheckCmd) (bool, error) {
	_, err := org.certificate.DuplicateCheck(ctx, cmd)
	if err != nil {
		if commonrepository.IsErrorResourceNotExists(err) {
			return true, nil
		}
	}

	return false, err
}

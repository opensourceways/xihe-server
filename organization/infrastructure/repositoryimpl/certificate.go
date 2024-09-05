/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/common/domain/crypto"
	"github.com/opensourceways/xihe-server/common/domain/primitive"
	postgresql "github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/organization/domain"
	orgprimitive "github.com/opensourceways/xihe-server/organization/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/domain/repository"
)

// NewCertificateImpl creates a new certificate repository implementation.
func NewCertificateImpl(db postgresql.Impl, enc crypto.Encrypter) (*certificateRepoImpl, error) {
	certificateTableName = db.TableName()
	err := db.DB().AutoMigrate(&CertificateDO{})

	return &certificateRepoImpl{Impl: db, e: enc}, err
}

type certificateRepoImpl struct {
	postgresql.Impl
	e crypto.Encrypter
}

// Save saves an organization certificate.
func (impl *certificateRepoImpl) Save(cert domain.OrgCertificate) error {
	do, err := toCertificateDo(cert, impl.e)
	if err != nil {
		return err
	}

	return impl.DB().Save(&do).Error
}

// Find finds an organization certificate.
func (impl *certificateRepoImpl) Find(
	ctx context.Context, option repository.FindOption) (domain.OrgCertificate, error) {
	do := CertificateDO{}
	if option.OrgName != nil {
		do.OrgName = option.OrgName.Account()
	}

	if option.CertificateOrgName != nil {
		do.CertificateOrgName = option.CertificateOrgName.AccountFullname()
	}

	if option.UnifiedSocialCreditCode != nil {
		do.USCC = option.UnifiedSocialCreditCode.USCC()
	}

	if err := impl.GetRecord(ctx, &do, &do); err != nil {
		return domain.OrgCertificate{}, err
	}

	return do.toCertificate(impl.e)
}

func (impl *certificateRepoImpl) GetAllCount() (total int64, err error) {
	query := impl.DB()
	query = query.Where(impl.EqualQuery(fieldStatus), orgprimitive.NewPassedStatus().CertificateStatus())
	if err = query.Count(&total).Error; err != nil {
		return
	}
	return
}

// get data by paging
func (impl *certificateRepoImpl) FindList(ctx context.Context, PageNum, PageSize int,
	orgType orgprimitive.CertificateOrgType) ([]domain.OrgCertificate, int, error) {
	query := impl.DB()
	if orgType != nil {
		query = query.Where(impl.EqualQuery(fieldCertOrgType), orgType.CertificateOrgType())
	}
	query = query.Where(impl.EqualQuery(fieldStatus), orgprimitive.NewPassedStatus().CertificateStatus())
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return []domain.OrgCertificate{}, 0, err
	}
	offset := 0
	if PageNum > 0 && PageSize > 0 {
		offset = (PageNum - 1) * PageSize
	}
	if offset > 0 {
		query = query.Limit(PageSize).Offset(offset)
	} else {
		query = query.Limit(PageSize)
	}
	var dos []CertificateDO
	err := query.Find(&dos).Error
	if err != nil || len(dos) == 0 {
		return []domain.OrgCertificate{}, 0, err
	}
	var do []domain.OrgCertificate
	for _, value := range dos {
		d, err := value.toCertificate(impl.e)
		if err != nil {
			return []domain.OrgCertificate{}, 0, err
		} else {
			do = append(do, d)
		}
	}
	return do, int(total), nil
}

func (impl *certificateRepoImpl) FinAllName(ctx context.Context) ([]string, error) {
	query := impl.DB()
	query = query.Where(impl.EqualQuery(fieldStatus),
		orgprimitive.NewPassedStatus().CertificateStatus())
	var dos []CertificateDO
	if err := query.Find(&dos).Error; err != nil {
		return nil, err
	}
	var result []string
	for _, value := range dos {
		result = append(result, value.OrgName)
	}
	return result, nil
}

// DuplicateCheck checks if the certificate already exists.
func (impl *certificateRepoImpl) DuplicateCheck(
	ctx context.Context, option repository.FindOption) (domain.OrgCertificate, error) {
	var do CertificateDO

	queryOr := impl.DB().Order(fieldID)
	if option.OrgName != nil {
		queryOr.Or(impl.EqualQuery(fieldOrg), option.OrgName.Account())
	}

	if option.CertificateOrgName != nil {
		queryOr.Or(impl.EqualQuery(fieldCertOrgName), option.CertificateOrgName.AccountFullname())
	}

	if option.UnifiedSocialCreditCode != nil {
		queryOr.Or(impl.EqualQuery(fieldUSCC), option.UnifiedSocialCreditCode.USCC())
	}

	if option.Phone != nil {
		queryOr.Or(impl.EqualQuery(fieldPhone), option.Phone.PhoneNumber())
	}

	query := impl.WithContext(ctx).
		Where(impl.EqualQuery(fieldStatus), orgprimitive.NewPassedStatus().CertificateStatus()).
		Where(queryOr)

	if err := impl.GetRecord(ctx, query, &do); err != nil {
		return domain.OrgCertificate{}, err
	}

	return do.toCertificate(impl.e)
}

// DeleteByOrgName deletes the organization certificate.
func (impl *certificateRepoImpl) DeleteByOrgName(orgName primitive.Account) error {
	return impl.DB().Where(impl.EqualQuery(fieldOrg), orgName.Account()).Delete(&CertificateDO{}).Error
}

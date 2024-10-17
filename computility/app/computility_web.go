/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package app provides application service
package app

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/computility/domain"
	"github.com/opensourceways/xihe-server/computility/domain/repository"
)

// ComputilityAppService is an interface for computility internal application service
type ComputilityAppService interface {
	GetAccountDetail(domain.ComputilityAccountIndex) (AccountQuotaDetailDTO, error)
}

// NewComputilityAppService creates a new instance of ComputilityAppService
func NewComputilityAppService(
	orgAdapter repository.ComputilityOrgRepositoryAdapter,
	detailAdapter repository.ComputilityDetailRepositoryAdapter,
	accountAdapter repository.ComputilityAccountRepositoryAdapter,
) ComputilityAppService {
	return &computilityAppService{
		orgAdapter:     orgAdapter,
		detailAdapter:  detailAdapter,
		accountAdapter: accountAdapter,
	}
}

type computilityAppService struct {
	orgAdapter     repository.ComputilityOrgRepositoryAdapter
	accountAdapter repository.ComputilityAccountRepositoryAdapter
	detailAdapter  repository.ComputilityDetailRepositoryAdapter
}

func (s *computilityAppService) GetAccountDetail(index domain.ComputilityAccountIndex) (
	AccountQuotaDetailDTO, error,
) {
	account, err := s.accountAdapter.FindByAccountIndex(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			empty := AccountQuotaDetailDTO{
				UserName:    index.UserName.Account(),
				ComputeType: index.ComputeType.ComputilityType(),
			}

			return empty, nil
		}

		e := xerrors.Errorf("find computility account error | user:%s, compute type:%s | err: %w",
			index.UserName.Account(), index.ComputeType.ComputilityType(), err,
		)

		logrus.Error(e)

		return AccountQuotaDetailDTO{}, allerror.New(
			allerror.ErrorCodeComputilityAccountFindError,
			"find computility account failed", e)
	}

	r := toAccountQuotaDetailDTO(&account)

	return r, nil
}

/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/computility/domain"
	"github.com/opensourceways/xihe-server/computility/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

// ComputilityInternalAppService is an interface for computility internal application service
type ComputilityInternalAppService interface {
	UserQuotaConsume(CmdToUserQuotaUpdate) error
}

// NewComputilityInternalAppService creates a new instance of ComputilityInternalAppService
func NewComputilityInternalAppService(
	detailAdapter repository.ComputilityDetailRepositoryAdapter,
	accountAdapter repository.ComputilityAccountRepositoryAdapter,
	accountRecordAtapter repository.ComputilityAccountRecordRepositoryAdapter,
) ComputilityInternalAppService {
	return &computilityInternalAppService{
		detailAdapter:        detailAdapter,
		accountAdapter:       accountAdapter,
		accountRecordAtapter: accountRecordAtapter,
	}
}

type computilityInternalAppService struct {
	accountAdapter       repository.ComputilityAccountRepositoryAdapter
	detailAdapter        repository.ComputilityDetailRepositoryAdapter
	accountRecordAtapter repository.ComputilityAccountRecordRepositoryAdapter
}

func (s *computilityInternalAppService) UserQuotaConsume(cmd CmdToUserQuotaUpdate) error {
	if cmd.Index.ComputeType.IsCpu() {
		return nil
	}

	user := cmd.Index.UserName
	_, err := s.accountRecordAtapter.FindByRecordIndex(cmd.Index)
	if err == nil {
		logrus.Errorf("user:%s already bind space:%s", user, cmd.Index.SpaceId.Identity())

		return nil
	}

	b, err := s.accountAdapter.CheckAccountExist(user)
	if err != nil {
		return err
	}
	if !b {
		e := xerrors.Errorf("user %s no quota balance for %s",
			user.Account(), cmd.Index.ComputeType.ComputilityType())

		logrus.Errorf("consume quota error| %s", e)

		return allerror.New(
			allerror.ErrorCodeNoNpuPermission,
			"no quota balance", e)
	}

	index := domain.ComputilityAccountIndex{
		UserName:    user,
		ComputeType: cmd.Index.ComputeType,
	}

	account, err := s.accountAdapter.FindByAccountIndex(index)
	if err != nil {
		logrus.Errorf("find user:%s account failed, %s", cmd.Index.UserName.Account(), err)

		return err
	}

	balance := account.QuotaCount - account.UsedQuota
	if balance < 1 {
		e := xerrors.Errorf("user %s insufficient computing quota balance", user.Account())

		logrus.Errorf("consume quota error| %s", e)

		return allerror.New(
			allerror.ErrorCodeInsufficientQuota,
			"insufficient computing quota balance", e)
	}

	err = s.accountAdapter.ConsumeQuota(account, cmd.QuotaCount)
	if err != nil {
		return err
	}

	err = s.accountRecordAtapter.Add(&domain.ComputilityAccountRecord{
		ComputilityAccountRecordIndex: cmd.Index,
		CreatedAt:                     utils.Now(),
		QuotaCount:                    cmd.QuotaCount,
		Version:                       0,
	})
	if err != nil {
		return err
	}

	return nil
}

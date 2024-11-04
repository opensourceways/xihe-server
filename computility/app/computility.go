/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/computility/domain"
	"github.com/opensourceways/xihe-server/computility/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

// ComputilityInternalAppService is an interface for computility internal application service
type ComputilityInternalAppService interface {
	UserQuotaConsume(CmdToUserQuotaUpdate) error
	UserQuotaRelease(CmdToUserQuotaUpdate) error

	SpaceCreateSupply(CmdToSupplyRecord) error
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

func (s *computilityInternalAppService) UserQuotaRelease(cmd CmdToUserQuotaUpdate) error {
	if cmd.Index.ComputeType.IsCpu() {
		return nil
	}

	record, err := s.accountRecordAtapter.FindByRecordIndex(cmd.Index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Errorf("user:%s has not cosume record to release", cmd.Index.UserName.Account())

			return nil
		}

		return err
	}

	accountIndex := domain.ComputilityAccountIndex{
		UserName:    cmd.Index.UserName,
		ComputeType: cmd.Index.ComputeType,
	}

	account, err := s.accountAdapter.FindByAccountIndex(accountIndex)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Errorf("user:%s is not a computility account, can not release quota", cmd.Index.UserName.Account())
			return nil
		}
		return err
	}

	if account.UsedQuota == 0 {
		logrus.Errorf("user:%s has no quota to release", cmd.Index.UserName.Account())
		return nil
	}

	err = s.accountAdapter.ReleaseQuota(account, cmd.QuotaCount)
	if err != nil {
		return err
	}

	err = s.accountRecordAtapter.Delete(record.Id)
	if err != nil {
		logrus.Errorf("delete user:%s account record failed, %s", cmd.Index.UserName.Account(), err)

		return err
	}

	err = s.accountAdapter.CancelAccount(accountIndex)
	if err != nil {
		logrus.Errorf("cancel user:%s account failed, %s", cmd.Index.UserName.Account(), err)

		return err
	}

	return nil
}

func (s *computilityInternalAppService) SpaceCreateSupply(cmd CmdToSupplyRecord) error {
	if cmd.Index.ComputeType.IsCpu() {
		return nil
	}

	b, err := s.accountAdapter.CheckAccountExist(cmd.Index.UserName)
	if err != nil {
		return err
	}
	if !b {
		e := xerrors.Errorf("user %s no permission for npu space", cmd.Index.UserName)

		return allerror.New(
			allerror.ErrorCodeNoNpuPermission,
			"no permission for npu space", e)
	}

	record, err := s.accountRecordAtapter.FindByRecordIndex(cmd.Index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			logrus.Errorf("user %s has not cosume record to release", cmd.Index.UserName.Account())
			return nil
		}
		return err
	}

	record.SpaceId = cmd.NewSpaceId

	err = s.accountRecordAtapter.Save(&record)
	if err != nil {
		logrus.Errorf("user %s no permission for %s space", cmd.Index.UserName, cmd.Index.ComputeType.ComputilityType())
	}

	return nil
}

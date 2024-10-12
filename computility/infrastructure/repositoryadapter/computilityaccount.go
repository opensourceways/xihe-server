/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"

	"github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/computility/domain"
	primitive "github.com/opensourceways/xihe-server/domain"
)

// computilityAccountAdapter is an implementation of the ComputilityAccountRepository interface.
type computilityAccountAdapter struct {
	daoImpl
}

// Add adds a new computility account to the database and returns an error if any occurs.
func (adapter *computilityAccountAdapter) Add(
	d *domain.ComputilityAccount,
) error {
	d.Id = primitive.CreateIdentity(primitive.GetId())

	do := toComputilityAccountDO(d)
	return adapter.db().Clauses(clause.Returning{}).Create(&do).Error

}

// Delete deletes a computility account in the database and returns an error if any occurs.
func (adapter *computilityAccountAdapter) Delete(id primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&computilityAccountDO{Id: id.Integer()},
	)
}

// FindByUserName finds a computility account in the repository based on the username.
func (adapter *computilityAccountAdapter) FindByAccountIndex(index domain.ComputilityAccountIndex) (
	domain.ComputilityAccount, error,
) {
	do := computilityAccountDO{
		UserName:    index.UserName.Account(),
		ComputeType: index.ComputeType.ComputilityType(),
	}

	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
	result := computilityAccountDO{}
	if err := adapter.daoImpl.GetRecord(&do, &result); err != nil {
		return domain.ComputilityAccount{}, err
	}

	return result.toComputilityAccount(), nil
}

// CheckAccountExist returns if account exists in the database and returns an error if any occurs.
func (adapter *computilityAccountAdapter) CheckAccountExist(userName primitive.Account) (
	bool, error,
) {
	do := computilityAccountDO{UserName: userName.Account()}

	result := computilityAccountDO{}
	if err := adapter.daoImpl.GetRecord(&do, &result); err != nil {
		if repository.IsErrorResourceNotExists(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// ConsumeQuota updates used_quota field
func (adapter *computilityAccountAdapter) ConsumeQuota(account domain.ComputilityAccount, quota int) error {
	do := toComputilityAccountDO(&account)

	do.Version += 1
	do.UsedQuota = do.UsedQuota + quota

	err := adapter.updateUsedQuota(do, account.Version)
	if err != nil {
		logrus.Errorf("computility db | increase %v used quota of user: %s failed", quota, account.UserName.Account())

		return err
	}

	logrus.Infof("computility db | increase %v used quota of user: %s success", quota, account.UserName.Account())

	return nil
}

// ReleaseQuota updates used_quota field
func (adapter *computilityAccountAdapter) ReleaseQuota(account domain.ComputilityAccount, quota int) error {
	do := toComputilityAccountDO(&account)

	do.Version += 1
	do.UsedQuota = do.UsedQuota - quota

	err := adapter.updateUsedQuota(do, account.Version)
	if err != nil {
		logrus.Errorf("computility db | decrease %v used quota of user: %s failed", quota, account.UserName.Account())

		return err
	}

	logrus.Infof("computility db | decrease %v used quota of user: %s success", quota, account.UserName.Account())

	return nil
}

// DecreaseAccountAssignedQuota updates quota_count field
func (adapter *computilityAccountAdapter) DecreaseAccountAssignedQuota(
	account domain.ComputilityAccount, quota int,
) error {
	do := toComputilityAccountDO(&account)

	do.Version += 1
	do.QuotaCount = do.QuotaCount - quota

	err := adapter.updateQuotaCount(do, account.Version)
	if err != nil {
		logrus.Errorf("computility db | decrease %v quota of user: %s failed", quota, account.UserName.Account())

		return err
	}

	logrus.Infof("computility db | decrease %v quota of user: %s success", quota, account.UserName.Account())

	return nil
}

// IncreaseAccountAssignedQuota updates quota_count field
func (adapter *computilityAccountAdapter) IncreaseAccountAssignedQuota(
	account domain.ComputilityAccount, quota int,
) error {
	do := toComputilityAccountDO(&account)

	do.Version += 1
	do.QuotaCount = do.QuotaCount + quota

	err := adapter.updateQuotaCount(do, account.Version)
	if err != nil {
		logrus.Errorf("computility db | increase %v quota of user: %s failed", quota, account.UserName.Account())

		return err
	}

	logrus.Infof("computility db | increase %v quota of user: %s success", quota, account.UserName.Account())

	return nil
}

// CancelAccount deletes a computility account in the repository based on the ComputilityAccountIndex.
func (adapter *computilityAccountAdapter) CancelAccount(index domain.ComputilityAccountIndex) error {
	account, err := adapter.FindByAccountIndex(index)
	if err != nil {
		return err
	}

	if account.QuotaCount == 0 && account.UsedQuota == 0 {
		err = adapter.Delete(account.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateUsedQuota updates the used quota for the Computility account with the specified version.
func (adapter *computilityAccountAdapter) updateUsedQuota(do computilityAccountDO, version int) error {
	result := adapter.db().Model(
		&computilityAccountDO{Id: do.Id},
	).Where(
		equalQuery(filedVersion), version,
	).Select(`*`).Omit(fieldCreatedAt, fieldQuotaCount).Updates(&do)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(errors.New("resource not found"))
	}

	return nil
}

// updateQuotaCount updates the quota count for the Computility account with the specified version.
func (adapter *computilityAccountAdapter) updateQuotaCount(do computilityAccountDO, version int) error {
	result := adapter.db().Model(
		&computilityAccountDO{Id: do.Id},
	).Where(
		equalQuery(filedVersion), version,
	).Select(`*`).Omit(fieldCreatedAt, fieldUsedQuota).Updates(&do)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(errors.New("resource not found"))
	}

	return nil
}

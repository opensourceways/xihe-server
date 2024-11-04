/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import (
	primitive "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/computility/domain"
)

var (
	computilityAccountTableName = ""
)

// TableName returns the table name for the domain.ComputilityAccount object.
func (do *computilityAccountDO) TableName() string {
	return computilityAccountTableName
}

// computilityAccountDO represents the database table for storing computility account information.
type computilityAccountDO struct {
	Id          int64  `gorm:"primaryKey"`
	UserName    string `gorm:"column:user_name;index:account_index,priority:1"`
	UsedQuota   int    `gorm:"column:used_quota"`
	CreatedAt   int64  `gorm:"column:created_at"`
	QuotaCount  int    `gorm:"column:quota_count"`
	ComputeType string `gorm:"column:compute_type;index:account_index,priority:2"`

	Version int `gorm:"column:version"`
}

// toComputilityAccountDO converts a domain.ComputilityAccount object to a computilityAccountDO object.
func toComputilityAccountDO(d *domain.ComputilityAccount) computilityAccountDO {
	return computilityAccountDO{
		Id:          d.Id.Integer(),
		UserName:    d.UserName.Account(),
		ComputeType: d.ComputeType.ComputilityType(),
		QuotaCount:  d.QuotaCount,
		UsedQuota:   d.UsedQuota,
		CreatedAt:   d.CreatedAt,
		Version:     d.Version,
	}
}

// toComputilityAccount converts a computilityAccountDO object to a domain.ComputilityAccount object.
func (do *computilityAccountDO) toComputilityAccount() domain.ComputilityAccount {
	return domain.ComputilityAccount{
		Id: primitive.CreateIdentity(do.Id),
		ComputilityAccountIndex: domain.ComputilityAccountIndex{
			UserName:    primitive.CreateAccount(do.UserName),
			ComputeType: primitive.CreateComputilityType(do.ComputeType),
		},
		UsedQuota:  do.UsedQuota,
		CreatedAt:  do.CreatedAt,
		QuotaCount: do.QuotaCount,
		Version:    do.Version,
	}
}

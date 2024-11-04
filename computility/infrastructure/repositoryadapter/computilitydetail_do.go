/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import (
	primitive "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/computility/domain"
)

var (
	computilityDetailTableName = ""
)

func (do *computilityDetailDO) TableName() string {
	return computilityDetailTableName
}

type computilityDetailDO struct {
	Id          int64  `gorm:"primarykey"`
	OrgName     string `gorm:"column:org_name;index:orgname_index"`
	UserName    string `gorm:"column:user_name;index:username_index"`
	CreatedAt   int64  `gorm:"column:created_at"`
	QuotaCount  int    `gorm:"column:quota_count"`
	ComputeType string `gorm:"column:compute_type"`

	Version int `gorm:"column:version"`
}

func toComputilityDetailDO(d *domain.ComputilityDetail) computilityDetailDO {
	return computilityDetailDO{
		Id:          d.Id.Integer(),
		UserName:    d.UserName.Account(),
		OrgName:     d.OrgName.Account(),
		CreatedAt:   d.CreatedAt,
		QuotaCount:  d.QuotaCount,
		ComputeType: d.ComputeType.ComputilityType(),
		Version:     d.Version,
	}
}

func (do *computilityDetailDO) toComputilityDetail() domain.ComputilityDetail {
	return domain.ComputilityDetail{
		Id: primitive.CreateIdentity(do.Id),
		ComputilityIndex: domain.ComputilityIndex{
			UserName: primitive.CreateAccount(do.UserName),
			OrgName:  primitive.CreateAccount(do.OrgName),
		},
		CreatedAt:   do.CreatedAt,
		ComputeType: primitive.CreateComputilityType(do.ComputeType),
		QuotaCount:  do.QuotaCount,
		Version:     do.Version,
	}
}

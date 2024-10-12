/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import (
	"github.com/opensourceways/xihe-server/computility/domain"
	primitive "github.com/opensourceways/xihe-server/domain"
)

var (
	computilityAccountRecordTableName = ""
)

func (do *computilityAccountRecordDO) TableName() string {
	return computilityAccountRecordTableName
}

type computilityAccountRecordDO struct {
	Id          int64  `gorm:"primaryKey"`
	UserName    string `gorm:"column:user_name;index:record_index,unique,priority:1"`
	SpaceId     int64  `gorm:"column:space_id;index:record_index,unique,priority:2"`
	CreatedAt   int64  `gorm:"column:created_at"`
	QuotaCount  int    `gorm:"column:quota_count"`
	ComputeType string `gorm:"column:compute_type"`

	Version int `gorm:"column:version"`
}

func toComputilityAccountRecordDO(d *domain.ComputilityAccountRecord) computilityAccountRecordDO {
	return computilityAccountRecordDO{
		Id:          d.Id.Integer(),
		UserName:    d.UserName.Account(),
		SpaceId:     d.SpaceId.Integer(),
		QuotaCount:  d.QuotaCount,
		CreatedAt:   d.CreatedAt,
		ComputeType: d.ComputeType.ComputilityType(),
		Version:     d.Version,
	}
}

func (do *computilityAccountRecordDO) toComputilityAccountRecord() domain.ComputilityAccountRecord {
	return domain.ComputilityAccountRecord{
		Id: primitive.CreateIdentity(do.Id),
		ComputilityAccountRecordIndex: domain.ComputilityAccountRecordIndex{
			UserName:    primitive.CreateAccount(do.UserName),
			SpaceId:     primitive.CreateIdentity(do.SpaceId),
			ComputeType: primitive.CreateComputilityType(do.ComputeType),
		},
		CreatedAt:  do.CreatedAt,
		QuotaCount: do.QuotaCount,
		Version:    do.Version,
	}
}

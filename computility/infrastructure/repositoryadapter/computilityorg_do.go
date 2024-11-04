/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import (
	primitive "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/computility/domain"
)

const (
	filedId          = "id"
	filedVersion     = "version"
	filedOrgName     = "org_name"
	filedUserName    = "user_name"
	fieldUsedQuota   = "used_quota"
	fieldCreatedAt   = "created_at"
	fieldQuotaCount  = "quota_count"
	filedComputeType = "compute_type"
)

var (
	computilityOrgTableName = ""
)

func (do *computilityOrgDO) TableName() string {
	return computilityOrgTableName
}

type computilityOrgDO struct {
	Id                 int64  `gorm:"primarykey;autoIncrement"`
	OrgId              int64  `gorm:"column:org_id;index:orgid_index"`
	OrgName            string `gorm:"column:org_name;index:orgname_index"`
	UsedQuota          int    `gorm:"column:used_quota"`
	QuotaCount         int    `gorm:"column:quota_count"`
	ComputeType        string `gorm:"column:compute_type"`
	DefaultAssignQuota int    `gorm:"column:default_assign_quota"`

	Version int `gorm:"column:version"`
}

func toComputilityOrgDO(d *domain.ComputilityOrg) computilityOrgDO {
	return computilityOrgDO{
		Id:                 d.Id.Integer(),
		OrgId:              d.OrgId.Integer(),
		OrgName:            d.OrgName.Account(),
		UsedQuota:          d.UsedQuota,
		QuotaCount:         d.QuotaCount,
		ComputeType:        d.ComputeType.ComputilityType(),
		DefaultAssignQuota: d.DefaultAssignQuota,
		Version:            d.Version,
	}
}

func (do *computilityOrgDO) toComputilityOrg() domain.ComputilityOrg {
	return domain.ComputilityOrg{
		Id:                 primitive.CreateIdentity(do.Id),
		OrgId:              primitive.CreateIdentity(do.OrgId),
		OrgName:            primitive.CreateAccount(do.OrgName),
		UsedQuota:          do.UsedQuota,
		QuotaCount:         do.QuotaCount,
		ComputeType:        primitive.CreateComputilityType(do.ComputeType),
		DefaultAssignQuota: do.DefaultAssignQuota,
		Version:            do.Version,
	}
}

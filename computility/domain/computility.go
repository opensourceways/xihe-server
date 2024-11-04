/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

import (
	primitive "github.com/opensourceways/xihe-server/domain"
)

// ComputilityDetail represents the detail of computility
type ComputilityDetail struct {
	ComputilityIndex

	Id          primitive.Identity
	CreatedAt   int64
	QuotaCount  int
	ComputeType primitive.ComputilityType

	Version int
}

// ComputilityAccount represents an account in the Computility system.
type ComputilityAccount struct {
	ComputilityAccountIndex

	Id         primitive.Identity
	UsedQuota  int
	QuotaCount int
	CreatedAt  int64

	Version int
}

// ComputilityOrg represents an organization in the Computility system.
type ComputilityOrg struct {
	Id                 primitive.Identity
	OrgId              primitive.Identity
	OrgName            primitive.Account
	UsedQuota          int
	QuotaCount         int
	ComputeType        primitive.ComputilityType
	DefaultAssignQuota int

	Version int
}

// ComputilityIndex represents an index for Computility entities.
type ComputilityIndex struct {
	OrgName  primitive.Account
	UserName primitive.Account
}

// ComputilityAccountIndex represents an index for Computility accounts.
type ComputilityAccountIndex struct {
	UserName    primitive.Account
	ComputeType primitive.ComputilityType
}

// RecallInfoList represents a list of recall information.
type RecallInfoList struct {
	InfoList []RecallInfo
}

// RecallInfo represents recall information for a user.
type RecallInfo struct {
	UserName    primitive.Account
	QuotaCount  int
	ComputeType primitive.ComputilityType
}

// ComputilityAccountRecordIndex represents an index for Computility account records.
type ComputilityAccountRecordIndex struct {
	UserName    primitive.Account
	SpaceId     primitive.Identity
	ComputeType primitive.ComputilityType
}

// ComputilityAccountRecord represents a record of a Computility account.
type ComputilityAccountRecord struct {
	ComputilityAccountRecordIndex

	Id         primitive.Identity
	CreatedAt  int64
	QuotaCount int

	Version int
}

/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package repositoryimpl provides implementations of repository interfaces for the organization domain.
package repositoryimpl

import (
	postgresql "github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
)

const (
	fieldID          = "id"
	fieldName        = "name"
	fieldOwner       = "owner"
	fieldCount       = "count"
	fieldAccount     = "account"
	fieldBio         = "bio"
	fieldAvatarId    = "avatar_id"
	fieldVersion     = "version"
	fieldType        = "type"
	fieldUser        = "user_name"
	fieldOrg         = "org_name"
	fieldRole        = "role"
	fieldInvitee     = "user_name"
	fieldInviter     = "inviter"
	fieldStatus      = "status"
	fieldCertOrgName = "certificate_org_name"
	fieldUSCC        = "uscc"
	fieldPhone       = "phone"
	fieldCertOrgType = "certificate_org_type"
	fieldCreatedAt   = "created_at"
	fieldUpdateAt    = "updated_at"
)

var certificateTableName string

// Member represents a member in the database.
type Member struct {
	postgresql.CommonModel
	Username string `gorm:"column:user_name;index:username_index"`
	FullName string `gorm:"column:full_name"`
	UserId   int64  `gorm:"column:user_id;index:userid_index"`
	Orgname  string `gorm:"column:org_name;index:orgname_index"`
	OrgId    int64  `gorm:"column:org_id;index:orgid_index"`
	Role     string `gorm:"column:role"`
	Type     string `gorm:"column:type"`
	Version  int    `gorm:"column:version"`
}

// Approve both request and approve use the same DO
type Approve struct {
	postgresql.CommonModel

	Username string `gorm:"column:user_name;index:username_index"`
	UserId   int64  `gorm:"column:user_id;index:userid_index"`
	Orgname  string `gorm:"column:org_name;index:orgname_index"`
	OrgId    int64  `gorm:"column:org_id;index:orgid_index"`
	Role     string `gorm:"column:role"`
	Expire   int64  `gorm:"column:expire"`  // approve only attr
	Inviter  string `gorm:"column:inviter"` // approve only attr
	Status   string `gorm:"column:status;index:status_index"`
	Type     string `gorm:"column:type"`
	By       string `gorm:"column:by"`
	Msg      string `gorm:"column:msg"`
	Version  int    `gorm:"column:version"`
}

// CertificateDO represents a certificate in the database.
type CertificateDO struct {
	postgresql.CommonModel

	OrgName            string `gorm:"column:org_name;uniqueIndex"`
	Phone              string `gorm:"column:phone;index"`
	USCC               string `gorm:"column:uscc;index"`
	CertificateOrgName string `gorm:"column:certificate_org_name;index"`
	CertificateOrgType string `gorm:"column:certificate_org_type"`
	Status             string `gorm:"column:status"`
	Reason             string `gorm:"column:reason"`
	Identity           string `gorm:"identity"`
}

// TableName returns the table name for the CertificateDO struct.
func (do CertificateDO) TableName() string {
	return certificateTableName
}

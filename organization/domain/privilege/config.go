/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package privilege provides privilege config
package privilege

// Config represents the configuration structure.
type Config struct {
	Npu     PrivilegeConfig `json:"npu"`
	Disable PrivilegeConfig `json:"disable"`
}

// PrivilegeConfig represents the privilege configuration structure.
type PrivilegeConfig struct {
	Orgs []OrgIndex `json:"orgs"`
}

// OrgIndex represents an organization index structure.
type OrgIndex struct {
	OrgId   string `json:"org_id"`
	OrgName string `json:"org_name"`
}

/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package permission provides a data structure for defining permissions and rules.
package permission

// Config represents the configuration for permissions.
type Config struct {
	Permissions []PermObject `json:"permissions"`
}

// PermObject represents a permission object.
type PermObject struct {
	ObjectType string `json:"object_type"`
	Rules      []Rule `json:"rules"`
}

// Rule represents a permission rule.
type Rule struct {
	Role      string   `json:"role"`
	Operation []string `json:"operation"`
}

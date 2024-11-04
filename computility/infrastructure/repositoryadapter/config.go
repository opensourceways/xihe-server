/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package repositoryadapter provides an adapter for working with computility repositories.
package repositoryadapter

// Tables is a struct that represents tables of computility.
type Tables struct {
	ComputilityOrg           string `json:"computility_org"            required:"true"`
	ComputilityDetail        string `json:"computility_detail"         required:"true"`
	ComputilityAccount       string `json:"computility_account"        required:"true"`
	ComputilityAccountRecord string `json:"computility_account_record" required:"true"`
}

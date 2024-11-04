/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import "fmt"

// OrgQuotaCount is an interface that represents an OrgQuotaCount.
type OrgQuotaCount interface {
	OrgQuotaCount() int
}

// NewOrgQuotaCount creates a new OrgQuotaCount with the given number.
func NewOrgQuotaCount(v int) (OrgQuotaCount, error) {
	if v < 0 || v > 1000 {
		return nil, fmt.Errorf("invalid quota count, should between %d and %d",
			0, 1000)
	}

	return orgQuotaCount(v), nil
}

// CreateOrgQuotaCountount is usually called internally, such as repository.
func CreateOrgQuotaCountount(v int) OrgQuotaCount {
	return orgQuotaCount(v)
}

type orgQuotaCount int

// OrgQuotaCount returns the org quota count.
func (r orgQuotaCount) OrgQuotaCount() int {
	return int(r)
}

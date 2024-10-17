package app

import "github.com/opensourceways/xihe-server/computility/domain"

// CmdToUserQuotaUpdate is a struct used for user quota update.
type CmdToUserQuotaUpdate struct {
	Index      domain.ComputilityAccountRecordIndex
	QuotaCount int
}

// AccountQuotaDetailDTO is a struct used for account quota detail.
type AccountQuotaDetailDTO struct {
	UserName     string `json:"user_name"`
	UsedQuota    int    `json:"used_quota"`
	TotalQuota   int    `json:"total_quota"`
	ComputeType  string `json:"compute_type"`
	QuotaBalance int    `json:"quota_balance"`
}

// toAccountQuotaDetailDTO converts a domain.ComputilityAccount object to an AccountQuotaDetailDTO.
func toAccountQuotaDetailDTO(a *domain.ComputilityAccount) AccountQuotaDetailDTO {
	return AccountQuotaDetailDTO{
		UserName:     a.UserName.Account(),
		UsedQuota:    a.UsedQuota,
		TotalQuota:   a.QuotaCount,
		QuotaBalance: a.QuotaCount - a.UsedQuota,
		ComputeType:  a.ComputeType.ComputilityType(),
	}
}

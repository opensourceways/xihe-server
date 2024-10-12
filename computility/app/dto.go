package app

import "github.com/opensourceways/xihe-server/computility/domain"

// CmdToUserQuotaUpdate is a struct used for user quota update.
type CmdToUserQuotaUpdate struct {
	Index      domain.ComputilityAccountRecordIndex
	QuotaCount int
}

package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

// WorkSubmittedEvent
type WorkSubmittedEvent struct {
	Id            string `json:"id"`
	Phase         string `json:"phase"`
	OBSPath       string `json:"obs_path"`
	PlayerId      string `json:"pid"`
	CompetitionId string `json:"cid"`
}

// CompetitorAppliedEvent
type CompetitorAppliedEvent struct {
	Account         types.Account
	CompetitionName types.CompetitionName
}

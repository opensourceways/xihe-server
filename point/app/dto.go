package app

import (
	"time"

	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/point/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type CmdToAddPointItem struct {
	Account common.Account
	Type    string
	Desc    string
	Time    int64
}

func (cmd *CmdToAddPointItem) dateAndTime() (string, string) {
	now := time.Now().Unix()

	if cmd.Time > now || cmd.Time < (now-minValueOfInvlidTime) {
		return "", ""
	}

	return utils.DateAndTime(cmd.Time)
}

type UserPointDetailsDTO struct {
	Total   int              `json:"total"`
	Details []PointDetailDTO `json:"details"`
}

type PointDetailDTO struct {
	Type string `json:"type"`

	domain.PointDetail
}

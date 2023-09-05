package app

import (
	"time"

	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type CmdToAddPointsItem struct {
	Account common.Account
	Type    string
	Desc    string
	Time    int64
}

func (cmd *CmdToAddPointsItem) dateAndTime() (string, string) {
	now := time.Now().Unix()

	if cmd.Time > now || cmd.Time < (now-minValueOfInvlidTime) {
		return "", ""
	}

	return utils.DateAndTime(cmd.Time)
}

type UserPointsDetailsDTO struct {
	Total   int               `json:"total"`
	Details []PointsDetailDTO `json:"details"`
}

type PointsDetailDTO struct {
	Type string `json:"type"`

	domain.PointsDetail
}

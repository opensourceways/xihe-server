package domain

import (
	"fmt"

	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type CompetitionMsg comsg.MsgNormal

func NewCompetitionApplyMsg(
	user types.Account, competitionName string,
) *CompetitionMsg {
	desc := fmt.Sprintf("apply competiton %s", competitionName)

	return &CompetitionMsg{
		User:      user.Account(),
		Desc:      desc,
		CreatedAt: utils.Now(),
	}
}

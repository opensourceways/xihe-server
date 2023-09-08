package message

import (
	"fmt"

	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type CompetitionMsg comsg.MsgNormal

type CompetitionMessageProducer interface {
	NotifyCalcScore(*domain.SubmissionMessage) error
	SendApplyCompetitionMsg(*CompetitionMsg) error
}

func (msg *CompetitionMsg) GenApplyMsg(
	user types.Account, competitionName string,
) {
	desc := fmt.Sprintf("apply competiton %s", competitionName)

	*msg = CompetitionMsg{
		User:      user.Account(),
		Desc:      desc,
		CreatedAt: utils.Now(),
	}
}

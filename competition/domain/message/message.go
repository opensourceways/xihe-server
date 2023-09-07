package message

import (
	"fmt"

	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	MsgTypeCompetitionApply = "msg_type_competition_apply"
)

type CompetitionMsg comsg.MsgNormal

type CompetitionMessageProducer interface {
	NotifyCalcScore(*domain.SubmissionMessage) error
	SendCompetitionMsg(*CompetitionMsg) error
}

func (msg *CompetitionMsg) GenApplyMsg(
	user types.Account, competitionName string,
) {
	desc := fmt.Sprintf("apply competiton %s", competitionName)

	*msg = CompetitionMsg{
		Type:      MsgTypeCompetitionApply,
		User:      user.Account(),
		Desc:      desc,
		Details:   nil,
		CreatedAt: utils.Now(),
	}
}

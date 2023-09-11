package message

import (
	"github.com/opensourceways/xihe-server/competition/domain"
)

type CompetitionMessageProducer interface {
	NotifyCalcScore(*domain.SubmissionMessage) error
	SendCompetitionApplyMsg(*domain.CompetitionMsg) error
}

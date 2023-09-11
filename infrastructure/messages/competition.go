package messages

import (
	"github.com/opensourceways/xihe-server/competition/domain"
)

func (s sender) NotifyCalcScore(v *domain.SubmissionMessage) error {
	return s.send(topics.Submission, v)
}

func (s sender) SendCompetitionApplyMsg(v *domain.CompetitionMsg) error {
	v.Type = topics.CompetitionApply.Name

	return s.send(topics.CompetitionApply.Topic, v)
}

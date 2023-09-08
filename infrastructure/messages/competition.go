package messages

import (
	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/message"
)

func (s sender) NotifyCalcScore(v *domain.SubmissionMessage) error {
	return s.send(topics.Submission, v)
}

func (s sender) SendApplyCompetitionMsg(v *message.CompetitionMsg) error {
	v.Type = topics.ApplyCompetition.Name

	return s.send(topics.ApplyCompetition.Topic, v)
}

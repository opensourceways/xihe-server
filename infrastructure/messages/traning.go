package messages

import "github.com/opensourceways/xihe-server/domain"

func (s sender) CreateTraining(msg *domain.Msg) error {
	msg.Type = topics.TrainingCreate.Name

	return s.send(topics.TrainingCreate.Topic, &msg)
}

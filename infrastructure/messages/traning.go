package messages

import "github.com/opensourceways/xihe-server/domain/message"

func (s sender) CreateTraining(msg *message.MsgTraining) error {
	return s.send(topics.Training, &msg)
}

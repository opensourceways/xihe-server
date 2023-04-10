package messages

import "github.com/opensourceways/xihe-server/bigmodel/domain/message"

func (s sender) UpdateWuKongTask(v *message.MsgWuKongLinks) error {
	return s.send(topics.Async, v)
}

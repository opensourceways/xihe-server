package messages

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain/message"
)

func (s sender) SendBigModelMsg(v *message.MsgTask) error {
	return s.send(topics.BigModel, v)
}

type BigModelMessageHandler interface {
	HandleEventBigModelWuKongInferenceStart(*message.MsgTask) error
	HandleEventBigModelWuKongInferenceError(*message.MsgTask) error
	HandleEventBigModelWuKongAsyncInferenceStart(*message.MsgTask) error
	HandleEventBigModelWuKongAsyncInferenceFinish(*message.MsgTask) error
}

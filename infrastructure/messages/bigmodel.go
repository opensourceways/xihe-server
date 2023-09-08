package messages

import (
	bigmodelmsg "github.com/opensourceways/xihe-server/bigmodel/domain/message"
)

// producer
func (s sender) SendBigModelMsg(v *bigmodelmsg.MsgTask) error {
	return s.send(topics.BigModel, v)
}

func (s sender) SendBigmodelPublicMsg(v *bigmodelmsg.MsgTask) error {
	v.Type = topics.PublicPicture.Name

	return s.send(topics.PublicPicture.Topic, v)
}

// comsumer
type BigModelMessageHandler interface {
	// wukong
	HandleEventBigModelWuKongInferenceStart(*bigmodelmsg.MsgTask) error
	HandleEventBigModelWuKongInferenceError(*bigmodelmsg.MsgTask) error
	HandleEventBigModelWuKongAsyncTaskStart(*bigmodelmsg.MsgTask) error
	HandleEventBigModelWuKongAsyncTaskFinish(*bigmodelmsg.MsgTask) error
}

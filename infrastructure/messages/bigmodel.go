package messages

import (
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
)

// producer
func (s sender) SendBigModelMsg(v *bigmodeldomain.MsgTask) error {
	return s.send(topics.BigModel, v)
}

func (s sender) SendBigmodelPublicMsg(v *bigmodeldomain.MsgTask) error {
	v.Type = topics.PublicPicture.Name

	return s.send(topics.PublicPicture.Topic, v)
}

// comsumer
type BigModelMessageHandler interface {
	// wukong
	HandleEventBigModelWuKongInferenceStart(*bigmodeldomain.MsgTask) error
	HandleEventBigModelWuKongInferenceError(*bigmodeldomain.MsgTask) error
	HandleEventBigModelWuKongAsyncTaskStart(*bigmodeldomain.MsgTask) error
	HandleEventBigModelWuKongAsyncTaskFinish(*bigmodeldomain.MsgTask) error
}

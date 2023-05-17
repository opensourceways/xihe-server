package messages

import (
	"github.com/opensourceways/xihe-server/cloud/domain/message"
)

func (s sender) SubscribeCloud(v *message.MsgCloudConf) error {
	return s.send(topics.Cloud, v)
}

func (s sender) ReleasePod(v *message.MsgPod) error {
	return s.send(topics.Cloud, v)
}

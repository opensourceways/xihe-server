package messageadapter

import (
	cloudmsg "github.com/opensourceways/xihe-server/cloud/domain/message"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
)

type cloudRecordEventPublisher struct {
	sender message.Sender
}

func (p cloudRecordEventPublisher) Publish(event *cloudmsg.CloudRecordEvent) error {
	return p.sender.AddOperateLogForCloudCreated(event.Owner, event.ClouId)
}

func NewCloudRecordEventPublisher(topics *messages.Topics, p common.Publisher) *cloudRecordEventPublisher {
	return &cloudRecordEventPublisher{
		sender: messages.NewMessageSender(topics, p),
	}
}

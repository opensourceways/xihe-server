package message

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

type AsyncMessageProducer interface {
	message.Sender

	SendBigModelMsg(*domain.MsgTask) error
	SendBigmodelPublicMsg(*domain.MsgTask) error
}

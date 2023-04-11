package messages

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain/message"
	asyncdomain "github.com/opensourceways/xihe-server/async-server/domain"
)

func (s sender) CreateWuKongTask(v *message.MsgTask) error {
	return s.send(topics.Async, v)
}

type AsyncCreateWuKongTaskMessageHandler interface {
	HandleEventAsyncCreateWuKongTask(info *asyncdomain.WuKongRequest) error
}

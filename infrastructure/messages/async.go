package messages

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain/message"
	asyncrepo "github.com/opensourceways/xihe-server/async-server/domain/repository"
)

func (s sender) UpdateWuKongTask(v *message.MsgTask) error {
	return s.send(topics.Async, v)
}

type AsyncUpdateWuKongTaskMessageHandler interface {
	HandleEventAsyncTaskWuKongUpdate(info *asyncrepo.WuKongResp) error
}

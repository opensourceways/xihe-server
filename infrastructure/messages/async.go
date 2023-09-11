package messages

import (
	asyncrepo "github.com/opensourceways/xihe-server/async-server/domain/repository"
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
)

func (s sender) UpdateWuKongTask(v *bigmodeldomain.MsgTask) error {
	return s.send(topics.Async, v)
}

type AsyncUpdateWuKongTaskMessageHandler interface {
	HandleEventAsyncTaskWuKongUpdate(info *asyncrepo.WuKongResp) error
}

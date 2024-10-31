package messageaimpl

import (
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

func NewMessageAdapter(topics *Topics, p common.Publisher) *messageAdapter {
	return &messageAdapter{topics: *topics, publisher: p}
}

type messageAdapter struct {
	topics    Topics
	publisher common.Publisher
}

func (impl *messageAdapter) SendSpaceAppCreateMsg(v *domain.SpaceAppCreateEvent) error {
	return impl.publisher.Publish(impl.topics.SpaceAppCreated, v, nil)
}

type Topics struct {
	// aicc finetune create
	SpaceAppCreated string `json:"spaceapp_created"`
}

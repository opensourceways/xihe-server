package infrastructure

import (
	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/space/domain"
)

func NewSpaceProducer(topics *Topics, p message.Publisher) *spaceProducer {
	return &spaceProducer{topics: *topics, publisher: p}
}

type spaceProducer struct {
	topics    Topics
	publisher message.Publisher
}

func (impl *spaceProducer) SendDeletedEvent(e *domain.DeleteSpaceEvent) error {
	return impl.publisher.Publish(impl.topics.SpaceDeleted, e, nil)
}

type Topics struct {
	SpaceDeleted string `json:"space_deleted"`
}

package messageadapter

import (
	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain"
)

func NewPublisher(cfg *Config) *publisher {
	return &publisher{*cfg}
}

type publisher struct {
	cfg Config
}

func (impl *publisher) SendWorkSubmittedEvent(v *domain.WorkSubmittedEvent) error {
	return message.Publish(impl.cfg.WorkSubmitted.Topic, v, nil)
}

// Config
type Config struct {
	WorkSubmitted message.TopicConfig `json:"work_submitted"`
}

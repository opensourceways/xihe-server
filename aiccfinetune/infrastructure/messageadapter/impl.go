package messageadapter

import (
	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	common "github.com/opensourceways/xihe-server/common/domain/message"
)

func NewMessageAdapter(cfg *Config, p common.Publisher) *messageAdapter {
	return &messageAdapter{cfg: *cfg, publisher: p}
}

type messageAdapter struct {
	cfg       Config
	publisher common.Publisher
}

func (impl *messageAdapter) SendAICCFinetuneCreateMsg(v *domain.AICCFinetuneCreateEvent) error {
	cfg := &impl.cfg.AICCFinetuneCreated

	msg := common.MsgNormal{
		User: v.User.Account(),
		Details: map[string]string{
			"id":    v.Id,
			"model": v.Model,
			"task":  v.Task,
		},
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

type Config struct {
	// wukong
	AICCFinetuneCreated common.TopicConfig `json:"aiccfinetune"`
}

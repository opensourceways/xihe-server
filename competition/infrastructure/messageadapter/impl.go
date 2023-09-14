package messageadapter

import (
	"fmt"

	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func MessageAdapter(cfg *Config, p common.Publisher) *messageAdapter {
	return &messageAdapter{cfg: *cfg, publisher: p}
}

type messageAdapter struct {
	cfg       Config
	publisher common.Publisher
}

func (impl *messageAdapter) SendWorkSubmittedEvent(v *domain.WorkSubmittedEvent) error {
	return impl.publisher.Publish(impl.cfg.WorkSubmitted.Topic, v, nil)
}

func (impl *messageAdapter) SendCompetitorAppliedEvent(v *domain.CompetitorAppliedEvent) error {
	cfg := &impl.cfg.CompetitorApplied

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      fmt.Sprintf("applied competition of %s", v.CompetitionName),
		CreatedAt: utils.Now(),
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Config
type Config struct {
	WorkSubmitted     common.TopicConfig `json:"work_submitted" required:"true"`
	CompetitorApplied common.TopicConfig `json:"competitor_applied" required:"true"`
}

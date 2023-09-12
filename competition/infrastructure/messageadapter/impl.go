package messageadapter

import (
	"fmt"

	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func NewPublisher(cfg *Config) *publisher {
	return &publisher{*cfg}
}

type publisher struct {
	cfg Config
}

func (impl *publisher) SendWorkSubmittedEvent(v *domain.WorkSubmittedEvent) error {
	return common.Publish(impl.cfg.WorkSubmitted.Topic, v, nil)
}

func (impl *publisher) SendCompetitorAppliedEvent(v *domain.CompetitorAppliedEvent) error {
	cfg := &impl.cfg.CompetitorApplied

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      fmt.Sprintf("applied competition of %s", v.CompetitionName),
		CreatedAt: utils.Now(),
	}

	return common.Publish(cfg.Topic, &msg, nil)
}

// Config
type Config struct {
	WorkSubmitted     common.TopicConfig `json:"work_submitted"`
	CompetitorApplied common.TopicConfig `json:"competitor_applied"`
}

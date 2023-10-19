package config

import (
	"github.com/opensourceways/xihe-server/aiccfinetune/infrastructure/aiccfinetuneimpl"
	"github.com/opensourceways/xihe-server/aiccfinetune/infrastructure/messageadapter"
)

type Config struct {
	aiccfinetuneimpl.Config

	Message messageadapter.Config `json:"message"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Config,
		&cfg.Message,
	}
}

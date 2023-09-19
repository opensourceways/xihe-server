package config

import (
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/infrastructure/trainingimpl"
)

type trainingConfig struct {
	trainingimpl.Config

	Message messages.TrainingConfig `json:"message"`
}

func (cfg *trainingConfig) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Config,
		&cfg.Message,
	}
}

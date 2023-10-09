package config

import (
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/messageadapter"
)

type Config struct {
	bigmodels.Config

	Message messageadapter.Config `json:"message"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Config,
		&cfg.Message,
	}
}

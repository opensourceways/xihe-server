package config

import (
	"github.com/opensourceways/xihe-server/infrastructure/authingimpl"
	"github.com/opensourceways/xihe-server/user/infrastructure/messageadapter"
)

type UserConfig struct {
	authingimpl.Config

	Message messageadapter.Config `json:"message"   required:"true"`
}

func (cfg *UserConfig) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Config,
		&cfg.Message,
	}
}

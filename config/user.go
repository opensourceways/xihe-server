package config

import (
	"github.com/opensourceways/xihe-server/infrastructure/authingimpl"
	"github.com/opensourceways/xihe-server/user/infrastructure/messageadapter"
)

type userConfig struct {
	authingimpl.Config

	Message messageadapter.Config `json:"message"   required:"true"`
}

func (cfg *userConfig) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Config,
		&cfg.Message,
	}
}

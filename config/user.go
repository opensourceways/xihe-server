package config

import (
	"github.com/opensourceways/xihe-server/user/infrastructure/messageadapter"
)

type userConfig struct {
	Message messageadapter.Config `json:"message"   required:"true"`
}

func (cfg *userConfig) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Message,
	}
}

package user

import (
	courseimpl "github.com/opensourceways/xihe-server/infrastructure/authingimpl"
	usermsg "github.com/opensourceways/xihe-server/user/infrastructure/messageadapter"
)

type Config struct {
	courseimpl.Config

	Message usermsg.Config `json:"message"   required:"true"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Config,
		&cfg.Message,
	}
}

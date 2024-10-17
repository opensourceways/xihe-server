package spaceapp

import (
	messageaimpl "github.com/opensourceways/xihe-server/spaceapp/infrastructure/messageimpl"
)

type Config struct {
	Message messageaimpl.Topics `json:"topics"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Message,
	}
}

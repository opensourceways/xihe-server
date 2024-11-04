package spaceapp

import (
	messageaimpl "github.com/opensourceways/xihe-server/spaceapp/infrastructure/messageimpl"
	"github.com/opensourceways/xihe-server/spaceapp/infrastructure/sseadapter"
)

type Config struct {
	Message    messageaimpl.Topics `json:"topics"`
	Controller sseadapter.Config   `json:"controller"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Message,
	}
}

package course

import (
	coursemsg "github.com/opensourceways/xihe-server/course/infrastructure/messageadapter"
)

type Config struct {
	Message coursemsg.Config `json:"message"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Message,
	}
}

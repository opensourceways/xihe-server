package competition

import (
	"github.com/opensourceways/xihe-server/common/config"
	competitionmsg "github.com/opensourceways/xihe-server/competition/infrastructure/messageadapter"
	"github.com/opensourceways/xihe-server/infrastructure/competitionimpl"
)

type Config struct {
	competitionimpl.Config

	Message competitionmsg.Config `json:"message"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Config,
		&cfg.Message,
	}
}

func (cfg *Config) SetDefault() {
	config.SetDefault(cfg.configItems())
}

func (cfg *Config) Validate() error {
	return config.Validate(cfg.configItems())
}

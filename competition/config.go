package competition

import (
	competitionmsg "github.com/opensourceways/xihe-server/competition/infrastructure/messageadapter"
	"github.com/opensourceways/xihe-server/infrastructure/competitionimpl"
)

type Config struct {
	competitionimpl.Config

	Message competitionmsg.Config `json:"message" required:"true"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Config,
		&cfg.Message,
	}
}

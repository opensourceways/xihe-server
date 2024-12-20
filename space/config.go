package space

import (
	"github.com/opensourceways/xihe-server/space/infrastructure"
	"github.com/opensourceways/xihe-server/space/infrastructure/repositoryimpl"
)

type Config struct {
	Topics infrastructure.Topics `json:"topics"`
	Tables repositoryimpl.Tables `json:"tables"`
}

// ConfigItems returns a slice of interfaces containing references to the Tables and Topics fields of the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
	}
}

// Init initializes the Config struct with default values.
func (cfg *Config) Init() {
}

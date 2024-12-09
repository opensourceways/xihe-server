package infrastructure

import (
	"github.com/opensourceways/xihe-server/filescan/infrastructure/repositoryadapter"
)

// Config is a struct that holds the configuration for tables and topics.
type FileScanConfig struct {
	Tables              repositoryadapter.Tables  `json:"tables"`
	ModerationPublisher ModerationPublisherConfig `json:"moderation_publisher"`
}

// ConfigItems returns a slice of interfaces containing references to the Tables and Topics fields of the Config struct.
func (cfg *FileScanConfig) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Tables,
	}
}

// Init initializes the Config struct with default values.
func (cfg *FileScanConfig) Init() {
}

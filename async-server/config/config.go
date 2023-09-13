package config

import (
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/async-server/infrastructure/poolimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/watchimpl"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
	"github.com/opensourceways/xihe-server/common/infrastructure/kafka"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
)

func LoadConfig(path string, cfg interface{}) error {
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	if f, ok := cfg.(ConfigSetDefault); ok {
		f.SetDefault()
	}

	if f, ok := cfg.(ConfigValidate); ok {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type ConfigValidate interface {
	Validate() error
}

type ConfigSetDefault interface {
	SetDefault()
}

type Config struct {
	MaxRetry int `json:"max_retry"`

	BigModel   bigmodels.Config `json:"bigmodel"     required:"true"`
	Postgresql PostgresqlConfig `json:"postgresql"   required:"true"`
	MQ         kafka.Config     `json:"mq"           required:"true"`
	MQTopics   messages.Topics  `json:"mq_topics"    required:"true"`
	Pool       poolimpl.Config  `json:"pool"         required:"true"`
	Watcher    watchimpl.Config `json:"watcher"      required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.BigModel,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Config,
		&cfg.MQ,
		&cfg.Pool,
	}
}

func (cfg *Config) SetDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(ConfigSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(ConfigValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type PostgresqlConfig struct {
	DB pgsql.Config `json:"db" required:"true"`

	repositoryimpl.Config
}

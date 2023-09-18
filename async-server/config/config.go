package config

import (
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/async-server/infrastructure/poolimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/watchimpl"
	bigmodel "github.com/opensourceways/xihe-server/bigmodel/config"
	common "github.com/opensourceways/xihe-server/common/config"
	"github.com/opensourceways/xihe-server/common/infrastructure/kafka"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
)

func LoadConfig(path string, cfg *Config) error {
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	cfg.setDefault()

	return cfg.validate()
}

type Config struct {
	MaxRetry int `json:"max_retry"`

	BigModel   bigmodel.Config  `json:"bigmodel"     required:"true"`
	Postgresql PostgresqlConfig `json:"postgresql"   required:"true"`
	MQ         kafka.Config     `json:"mq"           required:"true"`
	Pool       poolimpl.Config  `json:"pool"         required:"true"`
	Watcher    watchimpl.Config `json:"watcher"      required:"true"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.BigModel,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Config,
		&cfg.MQ,
		&cfg.Pool,
	}
}

func (cfg *Config) setDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	common.SetDefault(cfg)
}

func (cfg *Config) validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	return common.Validate(cfg)
}

type PostgresqlConfig struct {
	DB pgsql.Config `json:"db" required:"true"`

	repositoryimpl.Config
}

package config

import (
	"regexp"

	"github.com/opensourceways/community-robot-lib/utils"
	kfklib "github.com/opensourceways/kafka-lib/agent"
	redislib "github.com/opensourceways/redis-lib"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/watchimpl"
	serverconf "github.com/opensourceways/xihe-server/config"

	"github.com/opensourceways/xihe-server/async-server/infrastructure/poolimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
)

var reIpPort = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}:[1-9][0-9]*$`)

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
	Redis      serverconf.Redis `json:"redis"        required:"true"`
	MQ         serverconf.MQ    `json:"mq"           required:"true"`
	Pool       poolimpl.Config  `json:"pool"         required:"true"`
	Watcher    watchimpl.Config `json:"watcher"      required:"true"`
}

func (cfg *Config) GetKfkConfig() kfklib.Config {
	return kfklib.Config{
		Address: cfg.MQ.Address,
		Version: cfg.MQ.Version,
	}
}

func (cfg *Config) GetRedisConfig() redislib.Config {
	return redislib.Config{
		Address:  cfg.Redis.DB.Address,
		Password: cfg.Redis.DB.Password,
		DB:       cfg.Redis.DB.DB,
		Timeout:  cfg.Redis.DB.Timeout,
	}
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.BigModel,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Config,
		&cfg.Redis.DB,
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

package config

import (
	"errors"
	"regexp"
	"strings"

	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/async-server/infrastructure/poolimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
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
	MQ         MQ               `json:"mq"           required:"true"`
	Pool       poolimpl.Config  `json:"pool"         required:"true"`
}

func (cfg *Config) GetMQConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.MQ.ParseAddress(),
	}
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

type MQ struct {
	Address string          `json:"address" required:"true"`
	Topics  messages.Topics `json:"topics"  required:"true"`
}

func (cfg *MQ) Validate() error {
	if r := cfg.ParseAddress(); len(r) == 0 {
		return errors.New("invalid mq address")
	}

	return nil
}

func (cfg *MQ) ParseAddress() []string {
	v := strings.Split(cfg.Address, ",")
	r := make([]string, 0, len(v))
	for i := range v {
		if reIpPort.MatchString(v[i]) {
			r = append(r, v[i])
		}
	}

	return r
}

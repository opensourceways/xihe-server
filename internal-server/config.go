package main

import (
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/domain"
)

type configuration struct {
	Mongodb    config.Mongodb          `json:"mongodb"      required:"true"`
	Postgresql config.PostgresqlConfig `json:"postgresql"   required:"true"`
	Domain     domain.Config           `json:"domain"       required:"true"`
}

func (cfg *configuration) configItems() []interface{} {
	return []interface{}{
		&cfg.Mongodb,
		&cfg.Domain,
		&cfg.Postgresql.DB,
	}
}

func (cfg *configuration) SetDefault() {
	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(config.ConfigSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *configuration) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(config.ConfigValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *configuration) initDomainConfig() {
	domain.Init(&cfg.Domain)
}

package main

import (
	common "github.com/opensourceways/xihe-server/common/config"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func loadConfig(path string, cfg *configuration) error {
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	cfg.setDefault()

	return cfg.validate()
}

type configuration struct {
	Mongodb    config.Mongodb          `json:"mongodb"      required:"true"`
	Postgresql config.PostgresqlConfig `json:"postgresql"   required:"true"`
	Domain     domain.Config           `json:"domain"       required:"true"`
}

func (cfg *configuration) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Mongodb,
		&cfg.Domain,
		&cfg.Postgresql.DB,
	}
}

func (cfg *configuration) setDefault() {
	common.SetDefault(cfg)
}

func (cfg *configuration) validate() error {
	if err := utils.CheckConfig(cfg, ""); err != nil {
		return err
	}

	return common.Validate(cfg)
}

func (cfg *configuration) initDomainConfig() {
	domain.Init(&cfg.Domain)
}

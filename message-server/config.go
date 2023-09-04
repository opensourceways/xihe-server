package main

import (
	"github.com/opensourceways/community-robot-lib/utils"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	redislib "github.com/opensourceways/redis-lib"
	asyncrepoimpl "github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/cloud/infrastructure/cloudimpl"
	cloudrepoimpl "github.com/opensourceways/xihe-server/cloud/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/common/infrastructure/redis"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/evaluateimpl"
	"github.com/opensourceways/xihe-server/infrastructure/finetuneimpl"
	"github.com/opensourceways/xihe-server/infrastructure/inferenceimpl"
)

type configuration struct {
	MaxRetry         int    `json:"max_retry"`
	TrainingEndpoint string `json:"training_endpoint"  required:"true"`
	FinetuneEndpoint string `json:"finetune_endpoint"  required:"true"`

	Inference  inferenceimpl.Config `json:"inference"    required:"true"`
	Evaluate   evaluateConfig       `json:"evaluate"     required:"true"`
	Cloud      cloudConfig          `json:"cloud"        required:"true"`
	Mongodb    config.Mongodb       `json:"mongodb"      required:"true"`
	Postgresql PostgresqlConfig     `json:"postgresql"   required:"true"`
	Domain     domain.Config        `json:"domain"       required:"true"`
	MQ         config.MQ            `json:"mq"           required:"true"`
	Redis      redis.Config         `json:"redis"        required:"true"`
}

type PostgresqlConfig struct {
	DB pgsql.Config `json:"db" required:"true"`

	cloudconf cloudrepoimpl.Config
	asyncconf asyncrepoimpl.Config
}

func (cfg *configuration) getKfkConfig() kfklib.Config {
	return kfklib.Config{
		Address: cfg.MQ.Address,
	}
}

func (cfg *configuration) getRedisConfig() redislib.Config {
	return redislib.Config{
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		Timeout:  cfg.Redis.Timeout,
	}
}

func (cfg *configuration) configItems() []interface{} {
	return []interface{}{
		&cfg.Inference,
		&cfg.Evaluate,
		&cfg.Mongodb,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.cloudconf,
		&cfg.Postgresql.asyncconf,
		&cfg.Domain,
		&cfg.MQ,
		&cfg.Redis,
	}
}

func (cfg *configuration) SetDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

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

func (cfg *configuration) getFinetuneConfig() finetuneimpl.Config {
	return finetuneimpl.Config{
		Endpoint: cfg.FinetuneEndpoint,
	}
}

// evaluate
type evaluateConfig struct {
	SurvivalTime int `json:"survival_time"`

	evaluateimpl.Config
}

func (cfg *evaluateConfig) SetDefault() {
	if cfg.SurvivalTime <= 0 {
		cfg.SurvivalTime = 5 * 3600
	}

	var i interface{}
	i = &cfg.Config

	if f, ok := i.(config.ConfigSetDefault); ok {
		f.SetDefault()
	}
}

func (cfg *evaluateConfig) Validate() error {
	var i interface{}
	i = &cfg.Config

	if f, ok := i.(config.ConfigValidate); ok {
		return f.Validate()
	}

	return nil
}

// cloud
type cloudConfig struct {
	SurvivalTime int `json:"survival_time"`

	cloudimpl.Config
}

func (cfg *cloudConfig) SetDefault() {
	if cfg.SurvivalTime <= 0 {
		cfg.SurvivalTime = 5 * 3600
	}

	var i interface{}
	i = &cfg.Config

	if f, ok := i.(config.ConfigSetDefault); ok {
		f.SetDefault()
	}
}

func (cfg *cloudConfig) Validate() error {
	var i interface{}
	i = &cfg.Config

	if f, ok := i.(config.ConfigValidate); ok {
		return f.Validate()
	}

	return nil
}

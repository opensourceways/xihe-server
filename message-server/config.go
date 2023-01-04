package main

import (
	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/utils"

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

	Inference inferenceConfig `json:"inference"    required:"true"`
	Evaluate  evaluateConfig  `json:"evaluate"     required:"true"`
	Mongodb   config.Mongodb  `json:"mongodb"      required:"true"`
	Domain    domain.Config   `json:"domain"       required:"true"`
	MQ        config.MQ       `json:"mq"           required:"true"`
}

func (cfg *configuration) getMQConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.MQ.ParseAddress(),
	}
}

func (cfg *configuration) configItems() []interface{} {
	return []interface{}{
		&cfg.Inference,
		&cfg.Evaluate,
		&cfg.Mongodb,
		&cfg.Domain,
		&cfg.MQ,
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

type inferenceConfig struct {
	SurvivalTime int `json:"survival_time"`

	inferenceimpl.Config
}

func (cfg *inferenceConfig) SetDefault() {
	if cfg.SurvivalTime <= 0 {
		cfg.SurvivalTime = 5 * 3600
	}

	var i interface{}
	i = &cfg.Config

	if f, ok := i.(config.ConfigSetDefault); ok {
		f.SetDefault()
	}
}

func (cfg *inferenceConfig) Validate() error {
	var i interface{}
	i = &cfg.Config

	if f, ok := i.(config.ConfigValidate); ok {
		return f.Validate()
	}

	return nil
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

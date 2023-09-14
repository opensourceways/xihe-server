package main

import (
	"github.com/opensourceways/community-robot-lib/utils"

	asyncrepoimpl "github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/cloud/infrastructure/cloudimpl"
	cloudrepoimpl "github.com/opensourceways/xihe-server/cloud/infrastructure/repositoryimpl"
	common "github.com/opensourceways/xihe-server/common/config"
	"github.com/opensourceways/xihe-server/common/infrastructure/kafka"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/evaluateimpl"
	"github.com/opensourceways/xihe-server/infrastructure/finetuneimpl"
	"github.com/opensourceways/xihe-server/infrastructure/inferenceimpl"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/points"
	pointsdomain "github.com/opensourceways/xihe-server/points/domain"
)

func loadConfig(path string, cfg *configuration) error {
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	cfg.setDefault()

	return cfg.validate()
}

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
	MQ         kafka.Config         `json:"mq"           required:"true"`
	MQTopics   mqTopics             `json:"mq_topics"    required:"true"`
	Points     points.Config        `json:"points"`
}

type PostgresqlConfig struct {
	DB pgsql.Config `json:"db" required:"true"`

	cloudconf cloudrepoimpl.Config
	asyncconf asyncrepoimpl.Config
}

func (cfg *configuration) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Inference,
		&cfg.Evaluate,
		&cfg.Mongodb,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.cloudconf,
		&cfg.Postgresql.asyncconf,
		&cfg.Domain,
		&cfg.MQ,
		&cfg.MQTopics,
		&cfg.Points.Domain,
		&cfg.Points.Repo,
	}
}

func (cfg *configuration) setDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	common.SetDefault(cfg)
}

func (cfg *configuration) validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	return common.Validate(cfg)
}

func (cfg *configuration) initDomainConfig() {
	domain.Init(&cfg.Domain)
	pointsdomain.Init(&cfg.Points.Domain)
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

	common.SetDefault(&cfg.Config)
}

func (cfg *evaluateConfig) Validate() error {
	return common.Validate(&cfg.Config)
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

	common.SetDefault(&cfg.Config)
}

func (cfg *cloudConfig) Validate() error {
	return common.Validate(&cfg.Config)
}

type mqTopics struct {
	messages.Topics

	CompetitorApplied string `json:"competitor_applied" required:"true"`
	JupyterCreated    string `json:"jupyter_created"    required:"true"`
}

package main

import (
	aiccmq "github.com/opensourceways/xihe-server/aiccfinetune/messagequeue"

	asyncrepoimpl "github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	bigmodelmq "github.com/opensourceways/xihe-server/bigmodel/messagequeue"
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
	"github.com/opensourceways/xihe-server/messagequeue"
	pointsdomain "github.com/opensourceways/xihe-server/points/domain"
	pointsrepo "github.com/opensourceways/xihe-server/points/infrastructure/repositoryadapter"
	usermq "github.com/opensourceways/xihe-server/user/messagequeue"
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
	MaxRetry         int    `json:"max_retry"`
	FinetuneEndpoint string `json:"finetune_endpoint"  required:"true"`

	Inference    inferenceimpl.Config        `json:"inference"    required:"true"`
	Evaluate     evaluateConfig              `json:"evaluate"     required:"true"`
	Cloud        cloudConfig                 `json:"cloud"        required:"true"`
	Mongodb      config.Mongodb              `json:"mongodb"      required:"true"`
	Postgresql   PostgresqlConfig            `json:"postgresql"   required:"true"`
	Domain       domain.Config               `json:"domain"       required:"true"`
	MQ           kafka.Config                `json:"mq"           required:"true"`
	MQTopics     mqTopics                    `json:"mq_topics"    required:"true"`
	Points       pointsConfig                `json:"points"`
	Training     messagequeue.TrainingConfig `json:"training"`
	AICCFinetune aiccmq.AICCFinetuneConfig   `json:"aiccfinetune"`
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
		cfg.MaxRetry = 3
	}

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

// points
type pointsConfig struct {
	Domain pointsdomain.Config `json:"domain"`
	Repo   pointsrepo.Config   `json:"repo"`
}

func (cfg *pointsConfig) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Domain,
		&cfg.Repo,
	}
}

// user
type userConfig struct {
	BioSet       string `json:"bio_set"         required:"true"`
	AvatarSet    string `json:"avatar_set"      required:"true"`
	UserSignedUp string `json:"user_signed_up"  required:"true"`

	usermq.TopicConfig
}

// mqTopics
type mqTopics struct {
	messages.Topics

	SignIn string `json:"signin" required:"true"`

	// competition
	CompetitorApplied string `json:"competitor_applied" required:"true"`

	// cloud
	JupyterCreated string `json:"jupyter_created" required:"true"`

	// bigmodel
	BigModelTopics    bigmodelmq.TopicConfig `json:"bigmodel_topics"`
	PictureLiked      string                 `json:"picture_liked"            required:"true"`
	PicturePublicized string                 `json:"picture_publicized"       required:"true"`

	//course
	CourseApplied string `json:"course_applied" required:"true"`

	// training
	TrainingCreated string `json:"training_created" required:"true"`

	// resource
	ModelCreated      string `json:"model_created"      required:"true"`
	ProjectCreated    string `json:"project_created"    required:"true"`
	DatasetCreated    string `json:"dataset_created"    required:"true"`
	ModelDownloaded   string `json:"model_downloaded"   required:"true"`
	ProjectDownloaded string `json:"project_downloaded" required:"true"`
	DatasetDownloaded string `json:"dataset_downloaded" required:"true"`

	//user
	User userConfig `json:"user"`

	// aicc finetune
	AICCFinetuneCreated string `json:"aicc_finetune_created" required:"true"`
}

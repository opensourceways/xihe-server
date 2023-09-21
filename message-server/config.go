package main

import (
	"github.com/opensourceways/community-robot-lib/utils"

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
<<<<<<< HEAD
=======
	"github.com/opensourceways/xihe-server/user/infrastructure/messageadapter"
>>>>>>> 0b45df0 (update message of useraction)
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

	Inference  inferenceimpl.Config        `json:"inference"    required:"true"`
	Evaluate   evaluateConfig              `json:"evaluate"     required:"true"`
	Cloud      cloudConfig                 `json:"cloud"        required:"true"`
	Mongodb    config.Mongodb              `json:"mongodb"      required:"true"`
	Postgresql PostgresqlConfig            `json:"postgresql"   required:"true"`
	Domain     domain.Config               `json:"domain"       required:"true"`
	MQ         kafka.Config                `json:"mq"           required:"true"`
	MQTopics   mqTopics                    `json:"mq_topics"    required:"true"`
	Points     pointsConfig                `json:"points"`
	Training   messagequeue.TrainingConfig `json:"training"`
<<<<<<< HEAD
=======
<<<<<<< HEAD
<<<<<<< HEAD
	User       user.Config          `json:"user"         required:"true"`
=======
	User       user.Config                 `json:"user"         required:"true"`
>>>>>>> f320fe5 ( fix: update message of user config-)
=======
	User       messageadapter.Config       `json:"user"`
>>>>>>> 3afc252 (update message userinfo 566/11)
>>>>>>> 0b45df0 (update message of useraction)
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

// mqTopics
type mqTopics struct {
	messages.Topics

	// competition
	CompetitorApplied string `json:"competitor_applied" required:"true"`

	// cloud
	JupyterCreated string `json:"jupyter_created"    required:"true"`

	// bigmodel
	BigModelTopics    bigmodelmq.TopicConfig `json:"bigmodel_topics"`
	PicturePublicized string                 `json:"picture_publicized"  required:"true"`
	PictureLiked      string                 `json:"picture_liked"       required:"true"`

	//course
	CourseApplied string `json:"course_applied"                          required:"true"`

	// training
	TrainingCreated string `json:"training_created"                      required:"true"`

	//user
	UserSignedUp string `json:"user-signed-up"        required:"true"`
	BioSet       string `json:"bio_set"               required:"true"`
	AvatarSet    string `json:"avatar_set"            required:"true"`
}

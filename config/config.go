package config

import (
	"errors"
	"regexp"
	"strings"

	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/bigmodels"
	"github.com/opensourceways/xihe-server/infrastructure/challengeimpl"
	"github.com/opensourceways/xihe-server/infrastructure/competitionimpl"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/infrastructure/trainingimpl"
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
	MaxRetry        int `json:"max_retry"`
	ActivityKeepNum int `json:"activity_keep_num"`

	Competition competitionimpl.Config `json:"competition"  required:"true"`
	Challenge   challengeimpl.Config   `json:"challenge"    required:"true"`
	Training    trainingimpl.Config    `json:"training"     required:"true"`
	BigModel    bigmodels.Config       `json:"bigmodel"     required:"true"`
	Authing     AuthingService         `json:"authing"      required:"true"`
	Mongodb     Mongodb                `json:"mongodb"      required:"true"`
	Gitlab      gitlab.Config          `json:"gitlab"       required:"true"`
	Domain      domain.Config          `json:"domain"       required:"true"`
	App         app.Config             `json:"app"          required:"true"`
	API         controller.APIConfig   `json:"api"          required:"true"`
	MQ          MQ                     `json:"mq"           required:"true"`
}

func (cfg *Config) GetMQConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.MQ.ParseAddress(),
	}
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Competition,
		&cfg.Challenge,
		&cfg.Training,
		&cfg.BigModel,
		&cfg.Authing,
		&cfg.Domain,
		&cfg.Mongodb,
		&cfg.Gitlab,
		&cfg.App,
		&cfg.API,
		&cfg.MQ,
	}
}

func (cfg *Config) SetDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	if cfg.ActivityKeepNum <= 0 {
		cfg.ActivityKeepNum = 25
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

type Mongodb struct {
	DBName      string             `json:"db_name"       required:"true"`
	DBConn      string             `json:"db_conn"       required:"true"`
	Collections MongodbCollections `json:"collections"   required:"true"`
}

type MongodbCollections struct {
	Tag           string `json:"tag"             required:"true"`
	User          string `json:"user"            required:"true"`
	Like          string `json:"like"            required:"true"`
	Model         string `json:"model"           required:"true"`
	Login         string `json:"login"           required:"true"`
	LuoJia        string `json:"luojia"          required:"true"`
	WuKong        string `json:"wukong"          required:"true"`
	Dataset       string `json:"dataset"         required:"true"`
	Project       string `json:"project"         required:"true"`
	Activity      string `json:"activity"        required:"true"`
	Training      string `json:"training"        required:"true"`
	Finetune      string `json:"finetune"        required:"true"`
	Evaluate      string `json:"evaluate"        required:"true"`
	Inference     string `json:"inference"       required:"true"`
	AIQuestion    string `json:"aiquestion"      required:"true"`
	Competition   string `json:"competition"     required:"true"`
	QuestionPool  string `json:"question_pool"   required:"true"`
	WuKongPicture string `json:"wukong_picture"  required:"true"`
}

type AuthingService struct {
	APPId  string `json:"app_id" required:"true"`
	Secret string `json:"secret" required:"true"`
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

func (cfg *Config) InitDomainConfig() {
	domain.Init(&cfg.Domain)
}

func (cfg *Config) InitAppConfig() {
	app.Init(&cfg.App)
}

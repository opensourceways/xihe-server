package config

import (
	"errors"
	"regexp"
	"strings"

	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/utils"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	redislib "github.com/opensourceways/redis-lib"
	"github.com/opensourceways/xihe-server/app"
	asyncrepoimpl "github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
	cloudrepoimpl "github.com/opensourceways/xihe-server/cloud/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/common/infrastructure/redis"
	competitionmsg "github.com/opensourceways/xihe-server/competition/infrastructure/messageadapter"
	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/authingimpl"
	"github.com/opensourceways/xihe-server/infrastructure/challengeimpl"
	"github.com/opensourceways/xihe-server/infrastructure/competitionimpl"
	"github.com/opensourceways/xihe-server/infrastructure/finetuneimpl"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/infrastructure/trainingimpl"
	points "github.com/opensourceways/xihe-server/points/domain"
	pointsrepo "github.com/opensourceways/xihe-server/points/infrastructure/repositoryadapter"
)

var reIpPort = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}:[1-9][0-9]*$`)

const (
	kafkaDeaultVersion = "2.1.0"
)

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

type pointsConfig struct {
	Domain points.Config     `json:"domain"`
	Repo   pointsrepo.Config `json:"repo"`
}

type competitionConfig struct {
	competitionimpl.Config

	Message competitionmsg.Config `json:"message"`
}

type Config struct {
	MaxRetry        int `json:"max_retry"`
	ActivityKeepNum int `json:"activity_keep_num"`

	Competition competitionConfig    `json:"competition"  required:"true"`
	Challenge   challengeimpl.Config `json:"challenge"    required:"true"`
	Training    trainingimpl.Config  `json:"training"     required:"true"`
	Finetune    finetuneimpl.Config  `json:"finetune"     required:"true"`
	BigModel    bigmodels.Config     `json:"bigmodel"     required:"true"`
	Authing     authingimpl.Config   `json:"authing"      required:"true"`
	Mongodb     Mongodb              `json:"mongodb"      required:"true"`
	Postgresql  PostgresqlConfig     `json:"postgresql"   required:"true"`
	Redis       Redis                `json:"redis"        required:"true"`
	Gitlab      gitlab.Config        `json:"gitlab"       required:"true"`
	Domain      domain.Config        `json:"domain"       required:"true"`
	App         app.Config           `json:"app"          required:"true"`
	API         controller.APIConfig `json:"api"          required:"true"`
	MQ          MQ                   `json:"mq"           required:"true"`
	MQTopics    messages.Topics      `json:"mq_topics"    required:"true"`
	Points      pointsConfig         `json:"points"`
}

func (cfg *Config) GetMQConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.MQ.ParseAddress(),
	}
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
		&cfg.Competition.Config,
		&cfg.Competition.Message,
		&cfg.Challenge,
		&cfg.Training,
		&cfg.Finetune,
		&cfg.BigModel,
		&cfg.Authing,
		&cfg.Domain,
		&cfg.Mongodb,
		&cfg.Postgresql.DB,
		&cfg.Postgresql.Cloud,
		&cfg.Redis.DB,
		&cfg.Gitlab,
		&cfg.App,
		&cfg.API,
		&cfg.MQ,
		&cfg.MQTopics,
		&cfg.Points.Domain,
		&cfg.Points.Repo,
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
	DBCert      string             `json:"db_cert"       required:"true"`
	Collections MongodbCollections `json:"collections"   required:"true"`
}

type PostgresqlConfig struct {
	DB pgsql.Config `json:"db" required:"true"`

	Cloud cloudrepoimpl.Config
	Async asyncrepoimpl.Config
}

type Redis struct {
	DB redis.Config `json:"db" required:"true"`
}

type MongodbCollections struct {
	Tag               string `json:"tag"                    required:"true"`
	User              string `json:"user"                   required:"true"`
	Registration      string `json:"registration"           required:"true"`
	Like              string `json:"like"                   required:"true"`
	Model             string `json:"model"                  required:"true"`
	Login             string `json:"login"                  required:"true"`
	LuoJia            string `json:"luojia"                 required:"true"`
	WuKong            string `json:"wukong"                 required:"true"`
	Dataset           string `json:"dataset"                required:"true"`
	Project           string `json:"project"                required:"true"`
	Activity          string `json:"activity"               required:"true"`
	Training          string `json:"training"               required:"true"`
	Finetune          string `json:"finetune"               required:"true"`
	Evaluate          string `json:"evaluate"               required:"true"`
	Inference         string `json:"inference"              required:"true"`
	AIQuestion        string `json:"aiquestion"             required:"true"`
	Competition       string `json:"competition"            required:"true"`
	QuestionPool      string `json:"question_pool"          required:"true"`
	WuKongPicture     string `json:"wukong_picture"         required:"true"`
	CompetitionWork   string `json:"competition_work"       required:"true"`
	CompetitionPlayer string `json:"competition_player"     required:"true"`
	Course            string `json:"course"                 required:"true"`
	CoursePlayer      string `json:"course_player"          required:"true"`
	CourseWork        string `json:"course_work"            required:"true"`
	CourseRecord      string `json:"course_record"          required:"true"`
	CloudConf         string `json:"cloud_conf"             required:"true"`
	ApiApply          string `json:"api_apply"              required:"true"`
	ApiInfo           string `json:"api_info"               required:"true"`
	PointsTask        string `json:"points_task"            required:"true"`
	UserPoints        string `json:"user_points"            required:"true"`
}

type MQ struct {
	Address string `json:"address" required:"true"`
	Version string `json:"version"`
}

func (cfg *MQ) Validate() error {
	if r := cfg.ParseAddress(); len(r) == 0 {
		return errors.New("invalid mq address")
	}

	return nil
}

func (cfg *MQ) SetDefault() {
	if cfg.Version == "" {
		cfg.Version = kafkaDeaultVersion
	}
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

	points.Init(&cfg.Points.Domain)
}

func (cfg *Config) InitAppConfig() {
	app.Init(&cfg.App)
}

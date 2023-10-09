package config

import (
	redislib "github.com/opensourceways/redis-lib"

	"github.com/opensourceways/xihe-server/app"
	asyncrepoimpl "github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	bigmodel "github.com/opensourceways/xihe-server/bigmodel/config"
	cloudmsg "github.com/opensourceways/xihe-server/cloud/infrastructure/messageadapter"
	cloudrepoimpl "github.com/opensourceways/xihe-server/cloud/infrastructure/repositoryimpl"
	common "github.com/opensourceways/xihe-server/common/config"
	"github.com/opensourceways/xihe-server/common/infrastructure/kafka"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/common/infrastructure/redis"
	"github.com/opensourceways/xihe-server/competition"
	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/course"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/authingimpl"
	"github.com/opensourceways/xihe-server/infrastructure/challengeimpl"
	"github.com/opensourceways/xihe-server/infrastructure/finetuneimpl"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	pointsdomain "github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func LoadConfig(path string, cfg *Config) error {
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	cfg.setDefault()

	return cfg.validate()
}

type Config struct {
	MaxRetry        int `json:"max_retry"`
	ActivityKeepNum int `json:"activity_keep_num"`

	Competition competition.Config              `json:"competition"  required:"true"`
	Challenge   challengeimpl.Config            `json:"challenge"    required:"true"`
	Training    trainingConfig                  `json:"training"     required:"true"`
	Finetune    finetuneimpl.Config             `json:"finetune"     required:"true"`
	BigModel    bigmodel.Config                 `json:"bigmodel"     required:"true"`
	Authing     authingimpl.Config              `json:"authing"      required:"true"`
	Mongodb     Mongodb                         `json:"mongodb"      required:"true"`
	Postgresql  PostgresqlConfig                `json:"postgresql"   required:"true"`
	Redis       Redis                           `json:"redis"        required:"true"`
	Gitlab      gitlab.Config                   `json:"gitlab"       required:"true"`
	Domain      domain.Config                   `json:"domain"       required:"true"`
	App         app.Config                      `json:"app"          required:"true"`
	API         controller.APIConfig            `json:"api"          required:"true"`
	MQ          kafka.Config                    `json:"mq"           required:"true"`
	MQTopics    messages.Topics                 `json:"mq_topics"    required:"true"`
	SignIn      messages.SignInConfig           `json:"sign_in"      required:"true"`
	Points      pointsConfig                    `json:"points"`
	Cloud       cloudmsg.Config                 `json:"cloud"        required:"true"`
	Course      course.Config                   `json:"course"       required:"true"`
	Resource    messages.ResourceConfig         `json:"resource"     required:"true"`
	Download    messages.DownloadProducerConfig `json:"download"     required:"true"`
	User        userConfig                      `json:"user"`
	Like        messages.LikeConfig             `json:"like"`
}

func (cfg *Config) GetRedisConfig() redislib.Config {
	return redislib.Config{
		Address:  cfg.Redis.DB.Address,
		Password: cfg.Redis.DB.Password,
		DB:       cfg.Redis.DB.DB,
		Timeout:  cfg.Redis.DB.Timeout,
	}
}

func (cfg *Config) ConfigItems() []interface{} {
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
		&cfg.SignIn,
		&cfg.Points,
		&cfg.Cloud,
		&cfg.Download,
		&cfg.Course,
		&cfg.Resource,
		&cfg.User,
		&cfg.Like,
	}
}

func (cfg *Config) setDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	if cfg.ActivityKeepNum <= 0 {
		cfg.ActivityKeepNum = 25
	}

	common.SetDefault(cfg)
}

func (cfg *Config) validate() error {
	if err := utils.CheckConfig(cfg, ""); err != nil {
		return err
	}

	return common.Validate(cfg)
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

func (cfg *Config) InitDomainConfig() {
	domain.Init(&cfg.Domain)

	pointsdomain.Init(&cfg.Points.Domain)
}

func (cfg *Config) InitAppConfig() {
	app.Init(&cfg.App)
}

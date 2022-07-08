package config

import (
	"os"

	"github.com/opensourceways/xihe-server/utils"
)

func LoadConfig(path string) (*Config, error) {
	v := new(Config)

	if err := utils.LoadFromYaml(path, v); err != nil {
		return nil, err
	}

	v.setDefault()

	if err := v.validate(); err != nil {
		return nil, err
	}

	return v, nil
}

type Config struct {
	Authing AuthingService `json:"authing_service" required:"true"`
	Mongodb MongodbConfig  `json:"mongodb" required:"true"`
	Gitlab  GitlabConfig   `json:"gitlab" required:"true"`
}

func (cfg *Config) setDefault() {

	if os.Getenv("AUTHING_APP_ID") != "" {
		cfg.Authing.AppID = os.Getenv("AUTHING_APP_ID")
	}
	if os.Getenv("AUTHING_APP_SECRET") != "" {
		cfg.Authing.AppSecret = os.Getenv("AUTHING_APP_SECRET")
	}
	if os.Getenv("AUTHING_SECRET") != "" {
		cfg.Authing.Secret = os.Getenv("AUTHING_SECRET")
	}
	if os.Getenv("AUTHING_USER_POOL_ID") != "" {
		cfg.Authing.UserPoolId = os.Getenv("AUTHING_USER_POOL_ID")
	}
}

func (cfg *Config) validate() error {
	return nil
}

type MongodbConfig struct {
	MongodbConn       string `json:"mongodb_conn" required:"true"`
	DBName            string `json:"mongodb_db" required:"true"`
	ProjectCollection string `json:"project_collection" required:"true"`
}

type AuthingService struct {
	UserPoolId  string `json:"user_pool_id" required:"true"`
	Secret      string `json:"secret" required:"true"`
	AppID       string `json:"app_id" required:"true"`
	AppSecret   string `json:"app_secret" required:"true"`
	AuthingURL  string `json:"authing_url" required:"true"`
	RedirectURL string `json:"redirect_url" required:"true"`
}

type GitlabConfig struct {
	AcceesToken string `json:"accees_token" required:"true"`
	Host        string `json:"host" required:"true"`
}

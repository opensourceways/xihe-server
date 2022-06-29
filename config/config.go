package config

import (
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
	Mongodb MongodbConfig `json:"mongodb" required:"true"`
}

func (cfg *Config) setDefault() {
}

func (cfg *Config) validate() error {
	return nil
}

type MongodbConfig struct {
	MongodbConn       string `json:"mongodb_conn" required:"true"`
	DBName            string `json:"mongodb_db" required:"true"`
	ProjectCollection string `json:"project_collection" required:"true"`
}

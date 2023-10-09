package domain

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	MaxPointsOfDay int `json:"max_points_of_day"`
}

func (cfg *Config) SetDefault() {
	if cfg.MaxPointsOfDay <= 0 {
		cfg.MaxPointsOfDay = 50
	}
}

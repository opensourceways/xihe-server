package repositoryadapter

type Config struct {
	Keep int `json:"keep"`
}

func (cfg *Config) SetDefault() {
	if cfg.Keep <= 0 {
		cfg.Keep = 30
	}
}

package poolimpl

type Config struct {
	GoroutinePoolSize int `json:"goroutine_pool_size"`
}

func (cfg *Config) SetDefault() {
	if cfg.GoroutinePoolSize <= 0 {
		cfg.GoroutinePoolSize = 100
	}
}

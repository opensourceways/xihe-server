package poolimpl

import "github.com/panjf2000/ants/v2"

var gpool *ants.Pool

func Init(cfg *Config) (err error) {
	gpool, err = ants.NewPool(cfg.GoroutinePoolSize)
	if err != nil {
		return
	}

	return
}

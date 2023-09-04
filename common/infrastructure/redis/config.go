package redis

type Config struct {
	IdleSize int    `json:"idle_size"`
	NetWork  string `json:"network"`
	Address  string `json:"address"   required:"true"`
	Password string `json:"password"  required:"true"`
	KeyPair  string `json:"key_pair"  required:"true"`
	DB       int    `json:"db"`
	Timeout  int64  `json:"timeout"`
}

func (p *Config) SetDefault() {
	if p.IdleSize <= 0 {
		p.IdleSize = 20
	}

	if p.NetWork == "" {
		p.NetWork = "tcp"
	}

	if p.DB == 0 {
		p.DB = 0
	}

	if p.Timeout == 0 {
		p.Timeout = 10
	}
}

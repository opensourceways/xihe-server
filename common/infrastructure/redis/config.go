package redis

type Config struct {
	IdleSize int    `json:"idle_size" required:"true"`
	NetWork  string `json:"network"   required:"true"`
	Address  string `json:"address"   required:"true"`
	Password string `json:"password"  required:"true"`
	KeyPair  string `json:"key_pair"  required:"true"`
}

func (p *Config) SetDefault() {
	if p.IdleSize <= 0 {
		p.IdleSize = 20
	}

	if p.NetWork == "" {
		p.NetWork = "tcp"
	}
}

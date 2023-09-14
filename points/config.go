package points

import (
	points "github.com/opensourceways/xihe-server/points/domain"
	pointsrepo "github.com/opensourceways/xihe-server/points/infrastructure/repositoryadapter"
)

type Config struct {
	Domain points.Config     `json:"domain"`
	Repo   pointsrepo.Config `json:"repo"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Domain,
		&cfg.Repo,
	}
}

package config

import (
	points "github.com/opensourceways/xihe-server/points/domain"
	pointsrepo "github.com/opensourceways/xihe-server/points/infrastructure/repositoryadapter"
	"github.com/opensourceways/xihe-server/points/infrastructure/taskdocimpl"
)

type pointsConfig struct {
	Repo    pointsrepo.Config  `json:"repo"`
	Domain  points.Config      `json:"domain"`
	TaskDoc taskdocimpl.Config `json:"task_doc"`
}

func (cfg *pointsConfig) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Domain,
		&cfg.Repo,
		&cfg.TaskDoc,
	}
}

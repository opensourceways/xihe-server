package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Project interface {
	Save(domain.Project) (domain.Project, error)
}

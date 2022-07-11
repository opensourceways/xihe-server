package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Project interface {
	Save(domain.Project) (domain.Project, error)
	Update(domain.Project) (domain.Project, error)
	LikeCountIncrease(project_id, user_id string) error
}

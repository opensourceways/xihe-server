package infrastructure

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func NewProjectRepository(repo repository.Project) repository.Project {

	return Project{repo}
}

type Project struct {
	repo repository.Project
}

func (p Project) Save(entiy domain.Project) (r domain.Project, err error) {
	return
}
func (p Project) Update(entiy domain.Project) (r domain.Project, err error) {
	return
}
func (p Project) LikeCountIncrease(project_id, user_id string) error {

	return nil
}

package infrastructure

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func NewProjectInfra(repo repository.Project) repository.Project {

	return ProjectDB{repo}
}

type ProjectDB struct {
	repo repository.Project
}

func (p ProjectDB) Save(entiy domain.Project) (r domain.Project, err error) {
	return
}
func (p ProjectDB) Update(entiy domain.Project) (r domain.Project, err error) {
	return
}
func (p ProjectDB) LikeCountIncrease(project_id, user_id interface{}) error {

	return nil
}
func (p ProjectDB) GetBaseInfo(project_id interface{}) (data domain.Project, err error) {

	return
}

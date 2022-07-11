package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func NewProjectRepository(mapper ProjectMapper) repository.Project {
	return project{mapper}
}

type project struct {
	mapper ProjectMapper
}

func (impl project) Save(p domain.Project) (r domain.Project, err error) {
	return
}

type ProjectMapper interface {
}

package mongodb

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

func NewProjectRepository(collectionName string) app.ProjectRepository {
	return projectRepository{collectionName}
}

type projectRepository struct {
	collectionName string
}

func (p projectRepository) Save(domain.Project) error {
	return nil
}

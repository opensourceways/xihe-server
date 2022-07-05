package application

import (
	"github.com/opensourceways/xihe-server/domain/entity"
	"github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectApp struct {
	repo repository.ProjectRepository
}

func NewProjectAPP(repo repository.ProjectRepository) *ProjectApp {
	app := new(ProjectApp)
	app.repo = repo
	return app
}

func (f *ProjectApp) Save(item interface{}) (*mongo.InsertOneResult, error) {
	return f.repo.Save(item)
}

func (f *ProjectApp) GetAllFood() ([]entity.Project, error) {

	return nil, nil
}

package application

import (
	"github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
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

func (f *ProjectApp) Update(filter bson.M, item interface{}) (*mongo.UpdateResult, error) {
	return f.repo.Update(filter, item)
}

func (f *ProjectApp) Query(filter bson.M, offset, limit int64, order string) (result []bson.M, err error) {
	return f.repo.Query(filter, offset, limit, order)
}
func (f *ProjectApp) Get(filter bson.M) (result bson.M, err error) {
	return f.repo.Get(filter)
}

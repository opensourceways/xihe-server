package infrastructure

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectRepo struct {
	CommonRepo
}

func NewProjectRepository(mongodDB *mongo.Database) *ProjectRepo {
	repo := new(ProjectRepo)
	repo.Collection = mongodDB.Collection("project")
	return repo
}
func (r *ProjectRepo) Save(item interface{}) (*mongo.InsertOneResult, error) {
	return r.CommonRepo.Save(item)
}

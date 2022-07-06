package infrastructure

import (
	"go.mongodb.org/mongo-driver/bson"
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

func (r *ProjectRepo) Update(filter bson.M, item interface{}) (*mongo.UpdateResult, error) {
	return r.CommonRepo.Update(filter, item)
}

func (r *ProjectRepo) Query(filter bson.M, offset, limit int64, order string) (result []bson.M, err error) {
	return r.CommonRepo.Query(filter, offset, limit, order)
}

func (r *ProjectRepo) Get(filter bson.M) (result bson.M, err error) {
	return r.CommonRepo.Get(filter)
}

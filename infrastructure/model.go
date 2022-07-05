package infrastructure

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type ModelRepo struct {
	CommonRepo
}

func NewModelRepository(mongodDB *mongo.Database) *ModelRepo {
	repo := new(ModelRepo)
	repo.Collection = mongodDB.Collection("model")
	return repo
}

func (r *ModelRepo) Save(item interface{}) (*mongo.InsertOneResult, error) {
	return r.CommonRepo.Save(item)
}

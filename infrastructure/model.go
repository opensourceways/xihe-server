package infrastructure

import (
	"go.mongodb.org/mongo-driver/bson"
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

func (r *ModelRepo) Update(filter bson.M, item interface{}) (*mongo.UpdateResult, error) {
	return r.CommonRepo.Update(filter, item)
}

func (r *ModelRepo) Query(filter bson.M, offset, limit int64, order string) (result []bson.M, err error) {
	return r.CommonRepo.Query(filter, offset, limit, order)
}
func (r *ModelRepo) Get(filter bson.M) (result bson.M, err error) {
	return r.CommonRepo.Get(filter)
}

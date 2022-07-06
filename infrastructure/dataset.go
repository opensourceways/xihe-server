package infrastructure

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataSetRepo struct {
	CommonRepo
}

func NewDataSetRepository(mongodDB *mongo.Database) *DataSetRepo {
	repo := new(DataSetRepo)
	repo.Collection = mongodDB.Collection("dataset")
	return repo
}

func (r *DataSetRepo) Save(dataset interface{}) (*mongo.InsertOneResult, error) {
	return r.CommonRepo.Save(dataset)
}

func (r *DataSetRepo) Update(filter bson.M, item interface{}) (*mongo.UpdateResult, error) {
	return r.CommonRepo.Update(filter, item)
}

func (r *DataSetRepo) Query(filter bson.M, offset, limit int64, order string) (result []bson.M, err error) {
	return r.CommonRepo.Query(filter, offset, limit, order)
}
func (r *DataSetRepo) Get(filter bson.M) (result bson.M, err error) {
	return r.CommonRepo.Get(filter)
}

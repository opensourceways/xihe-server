package infrastructure

import (
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

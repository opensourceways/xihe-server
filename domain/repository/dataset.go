package repository

import "go.mongodb.org/mongo-driver/mongo"

type DataSetRepository interface {
	Save(item interface{}) (*mongo.InsertOneResult, error)
}

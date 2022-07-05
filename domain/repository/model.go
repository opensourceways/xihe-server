package repository

import "go.mongodb.org/mongo-driver/mongo"

type ModelRepository interface {
	Save(item interface{}) (*mongo.InsertOneResult, error)
}

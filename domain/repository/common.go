package repository

import "go.mongodb.org/mongo-driver/mongo"

type CommonRepository interface {
	Save(item interface{}) (*mongo.InsertOneResult, error)
}

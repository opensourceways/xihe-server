package repository

import "go.mongodb.org/mongo-driver/mongo"

type ProjectRepository interface {
	Save(item interface{}) (*mongo.InsertOneResult, error)
}

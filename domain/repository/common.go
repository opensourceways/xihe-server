package repository

import (
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

type CommonRepository interface {
	Save(item interface{}) (*mongo.InsertOneResult, error)
	Update(filter bson.M, item interface{}) (*mongo.UpdateResult, error)
	Query(filter bson.M, offset, limit int64, order string) (result []bson.M, err error)
	Get(filter bson.M) (result bson.M, err error)
}

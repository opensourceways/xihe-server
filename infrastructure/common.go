package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommonRepo struct {
	Collection *mongo.Collection
}

func (r *CommonRepo) Save(item interface{}) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(context.Background(), item)
}

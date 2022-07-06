package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommonRepo struct {
	Collection *mongo.Collection
}

func (r *CommonRepo) Save(item interface{}) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(context.Background(), item)
}

func (r *CommonRepo) Update(filter bson.M, item interface{}) (*mongo.UpdateResult, error) {
	return r.Collection.UpdateOne(context.Background(), filter, item, nil)
}

func (r *CommonRepo) Query(filter bson.M, offset, limit int64, order string) (result []bson.M, err error) {
	var cursor *mongo.Cursor
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{"updatedat", 1}})
	if cursor, err = r.Collection.Find(context.TODO(), filter, opts); err == nil {
		result = make([]bson.M, cursor.RemainingBatchLength())
		err = cursor.All(context.Background(), &result)
	}
	return
}

func (r *CommonRepo) Get(filter bson.M) (result bson.M, err error) {
	findResult := r.Collection.FindOne(context.TODO(), filter)
	resultRaw, err := findResult.DecodeBytes()
	if err != nil {
		return nil, err
	}
	result = make(bson.M)
	err = bson.Unmarshal(resultRaw, result)
	return
}

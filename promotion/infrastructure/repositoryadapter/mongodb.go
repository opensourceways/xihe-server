package repositoryadapter

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbClient interface {
	IsDocNotExists(error) bool
	IsDocExists(error) bool

	GetDoc(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error

	GetDocs(ctx context.Context, filterOfDoc bson.M, opts *options.FindOptions, result interface{}) error

	NewDocIfNotExist(ctx context.Context, filterOfDoc, docInfo bson.M) (string, error)

	PushElemArrayWithVersion(
		ctx context.Context, array string,
		filterOfDoc, value bson.M, version int, otherUpdate bson.M,
	) error

	GetArrayElem(
		ctx context.Context, array string,
		filterOfDoc, filterOfArray bson.M,
		project bson.M, result interface{},
	) error

	PushArrayElemAndInc(
		ctx context.Context,
		array string, filterOfDoc, value, updateCmd bson.M,
		version int,
	) error

	Count(ctx context.Context, filterOfDoc bson.M, opts *options.CountOptions) (int64, error)

	Aggregate(ctx context.Context, pipeline mongo.Pipeline, opts *options.AggregateOptions, result interface{}) error
}

func withContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second, // TODO use config
	)
	defer cancel()

	return f(ctx)
}

func genDoc(doc interface{}) (m bson.M, err error) {
	v, err := json.Marshal(doc)
	if err != nil {
		return
	}

	if err = json.Unmarshal(v, &m); err != nil {
		return
	}

	return
}

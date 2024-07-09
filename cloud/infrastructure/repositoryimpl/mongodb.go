package repositoryimpl

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodbClient interface {
	IsDocNotExists(error) bool
	IsDocExists(error) bool

	GetDoc(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error

	GetDocs(ctx context.Context, filterOfDoc bson.M, opts *options.FindOptions, result interface{}) error
}

func withContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second, // TODO use config
	)
	defer cancel()

	return f(ctx)
}

package repositoryimpl

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoCmdSet       = "$set"
	mongoCmdPush      = "$push"
	mongoCmdElemMatch = "$elemMatch"
	mongoCmdMatch     = "$match"
	mongoCmdEqual     = "$eq"
	mongoCmdCount     = "$count"
	mongoCmdInc       = "$inc"
)

type mongodbClient interface {
	IsDocNotExists(error) bool
	IsDocExists(error) bool
	Collection() *mongo.Collection
	ObjectIdFilter(s string) (bson.M, error)
	NewDocIfNotExist(ctx context.Context, filterOfDoc, docInfo bson.M) (string, error)
	GetDoc(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error
	GetDocs(ctx context.Context, filterOfDoc bson.M, opts *options.FindOptions, result interface{}) error
	UpdateDoc(ctx context.Context, filterOfDoc, update bson.M, op string, version int) error
	UpdateIncDoc(ctx context.Context, filterOfDoc, update bson.M, version int) error
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

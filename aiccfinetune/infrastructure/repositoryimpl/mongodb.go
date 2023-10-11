package repositoryimpl

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	mongoCmdSet  = "$set"
	mongoCmdPush = "$push"
)

var (
	errDocNotExists = errors.New("doc doesn't exist")
	errDocNoUpdate  = errors.New("no update")
)

type mongodbClient interface {
	IsDocNotExists(error) bool
	IsDocExists(error) bool

	Collection() *mongo.Collection

	ObjectIdFilter(s string) (bson.M, error)

	AppendElemMatchToFilter(array string, exists bool, cond, filter bson.M)

	GetDoc(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error

	GetDocs(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error

	NewDocIfNotExist(ctx context.Context, filterOfDoc, docInfo bson.M) (string, error)

	GetArrayElem(ctx context.Context, array string, filterOfDoc, filterOfArray bson.M, project bson.M, result interface{}) error

	UpdateDoc(ctx context.Context, filterOfDoc, update bson.M, op string, version int) error

	ModifyArrayElem(ctx context.Context, array string, filterOfDoc, filterOfArray, updateCmd bson.M, op string) (bool, error)

	PullArrayElem(ctx context.Context, array string, filterOfDoc, filterOfArray bson.M) error
	// modifyArrayElemWithoutVersion
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

package repositoryimpl

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type mongodbClient interface {
	IsDocNotExists(error) bool
	IsDocExists(error) bool

	ObjectIdFilter(s string) (bson.M, error)
	NewDocIfNotExist(ctx context.Context, filterOfDoc, docInfo bson.M) (string, error)
	GetDoc(ctx context.Context, filterOfDoc, project bson.M, result interface{}) error

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

	err = json.Unmarshal(v, &m)

	return
}

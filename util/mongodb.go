package util

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type client struct {
	c                 *mongo.Client
	db                *mongo.Database
	projectCollection string
	modelCollection   string
	datasetCollection string
	userCollection    string
}

func InitMongoDB() (*client, error) {
	user := GetConfig().Database.DBUser
	pswd := GetConfig().Database.Password
	host := GetConfig().Database.DBHost
	port := GetConfig().Database.DBPort
	dbName := GetConfig().Database.DbName
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, pswd, host, port)

	c, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}

	if err = withContext(c.Connect); err != nil {
		return nil, err
	}

	// verify if database connection is created successfully
	err = withContext(func(ctx context.Context) error {
		return c.Ping(ctx, nil)
	})
	if err != nil {
		return nil, err
	}

	cli := &client{
		c:  c,
		db: c.Database(dbName),
	}

	return cli, nil
}

func (this *client) Close() error {
	return withContext(this.c.Disconnect)
}

func (c *client) collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

func (this *client) doTransaction(f func(mongo.SessionContext) error) error {

	callback := func(sc mongo.SessionContext) (interface{}, error) {
		return nil, f(sc)
	}

	s, err := this.c.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start mongodb session: %s", err.Error())
	}

	ctx := context.Background()
	defer s.EndSession(ctx)

	_, err = s.WithTransaction(ctx, callback)
	return err
}

func objectIDToUID(oid primitive.ObjectID) string {
	return oid.Hex()
}

func toUID(oid interface{}) (string, error) {
	v, ok := oid.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("retrieve id failed")
	}
	return v.Hex(), nil
}

func withContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return f(ctx)
}

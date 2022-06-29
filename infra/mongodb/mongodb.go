package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cli *client

func Initialize(conn, dbName string) error {
	c, err := mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		return err
	}

	if err = withContext(c.Connect); err != nil {
		return err
	}

	// verify if database connection is created successfully
	err = withContext(func(ctx context.Context) error {
		return c.Ping(ctx, nil)
	})
	if err != nil {
		return err
	}

	cli = &client{
		c:  c,
		db: c.Database(dbName),
	}

	return nil
}

func Close() error {
	if cli != nil {
		return cli.disconnect()
	}

	return nil
}

type client struct {
	c  *mongo.Client
	db *mongo.Database
}

func (cli *client) disconnect() error {
	return withContext(cli.c.Disconnect)
}

func (cli *client) collection(name string) *mongo.Collection {
	return cli.db.Collection(name)
}

func withContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return f(ctx)
}

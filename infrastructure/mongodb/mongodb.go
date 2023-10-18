package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cli *client

func Initialize(conn, dbName, dbCert string) error {
	ca, err := ioutil.ReadFile(dbCert)
	if err != nil {
		return err
	}

	if err := os.Remove(dbCert); err != nil {
		return err
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		return fmt.Errorf("failed to append certs PEM")
	}

	tlsConfig := &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: true,
	}

	uri := conn
	clientOpt := options.Client().ApplyURI(uri)
	clientOpt.SetTLSConfig(tlsConfig)

	c, err := mongo.Connect(context.TODO(), clientOpt)
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
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second, // TODO use config
	)
	defer cancel()

	return f(ctx)
}

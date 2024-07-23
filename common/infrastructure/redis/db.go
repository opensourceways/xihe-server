package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

func Init(cfg *Config) error {
	ca, err := os.ReadFile(cfg.DBCert)
	if err != nil {
		return err
	}

	if err := os.Remove(cfg.DBCert); err != nil {
		return err
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		return fmt.Errorf("faild to append certs from PEM")
	}

	tlsConfig := &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: true, // #nosec G402 -- this is a false positive
	}

	client = redis.NewClient(&redis.Options{
		Addr:      cfg.Address,
		Password:  cfg.Password,
		DB:        0,
		TLSConfig: tlsConfig,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil

}

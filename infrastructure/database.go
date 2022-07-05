package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	ProjectRepo repository.ProjectRepository
	DataSetRepo repository.DataSetRepository
	ModelRepo   repository.ModelRepository
	UserRepo    repository.UserRepository
}

func withContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return f(ctx)
}

func NewRepositories() (*Repositories, error) {
	user := util.GetConfig().Database.DBUser
	pswd := util.GetConfig().Database.Password
	host := util.GetConfig().Database.DBHost
	port := util.GetConfig().Database.DBPort
	dbName := util.GetConfig().Database.DbName
	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, pswd, host, port)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connStr).SetConnectTimeout(5*time.Second))
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	// verify if database connection is created successfully
	err = withContext(func(ctx context.Context) error {
		return client.Ping(ctx, nil)
	})
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	return &Repositories{
		ProjectRepo: NewProjectRepository(db),
		DataSetRepo: NewDataSetRepository(db),
		ModelRepo:   NewModelRepository(db),
		UserRepo:    NewUserRepository(db),
	}, nil
}

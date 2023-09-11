package repositoryadapter

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain"
)

func TaskAdapter(cli mongodbClient) *taskAdapter {
	return &taskAdapter{cli}
}

type taskAdapter struct {
	cli mongodbClient
}

func (impl *taskAdapter) FindAllTasks() ([]domain.Task, error) {
	var dos []taskDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(ctx, nil, bson.M{fieldOlds: 0}, &dos)
	}

	if err := withContext(f); err != nil || len(dos) == 0 {
		return nil, err
	}

	r := make([]domain.Task, len(dos))
	for i := range dos {
		r[i] = dos[i].toTask()
	}

	return r, nil
}

func (impl *taskAdapter) Find(name string) (domain.Task, error) {
	var do taskDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(ctx, bson.M{fieldName: name}, bson.M{fieldOlds: 0}, &do)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return domain.Task{}, err
	}

	return do.toTask(), nil
}

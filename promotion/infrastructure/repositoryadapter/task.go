package repositoryadapter

import (
	"context"

	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func TaskAdapter(cli mongodbClient) repository.Task {
	return &taskAdapter{cli}
}

type taskAdapter struct {
	cli mongodbClient
}

func (impl *taskAdapter) Find(id string) (domain.Task, error) {
	var do taskDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, docIdFilter(id),
			nil, &do,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return domain.Task{}, err
	}

	return do.toTask()
}

func (impl *taskAdapter) FindAll() (tasks []domain.Task, err error) {
	var dos []taskDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, nil, nil, &dos,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return nil, err
	}

	tasks = make([]domain.Task, len(dos))
	for i := range dos {
		if tasks[i], err = dos[i].toTask(); err != nil {
			return
		}
	}

	return
}

func docIdFilter(id string) bson.M {
	return bson.M{fieldId: id}
}

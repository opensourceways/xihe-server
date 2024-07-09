package repositoryadapter

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain"
)

func TaskAdapter(cli mongodbClient) *taskAdapter {
	return &taskAdapter{cli}
}

type taskAdapter struct {
	cli mongodbClient
}

func (impl *taskAdapter) docFilter(tid string) bson.M {
	return bson.M{fieldId: tid}
}

func (impl *taskAdapter) Add(t *domain.Task) error {
	do := totaskDO(t)

	doc, err := do.doc()
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := impl.cli.NewDocIfNotExist(ctx, impl.docFilter(t.Id), doc)

		return err
	}

	if err = withContext(f); err != nil && impl.cli.IsDocExists(err) {
		err = repoerr.NewErrorDuplicateCreating(err)
	}

	return err
}

func (impl *taskAdapter) FindAllTasks() ([]domain.Task, error) {
	var dos []taskDO

	f := func(ctx context.Context) error {
		opts := options.FindOptions{}
		return impl.cli.GetDocs(ctx, nil, opts.SetProjection(bson.M{fieldOlds: 0}), &dos)
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

func (impl *taskAdapter) Find(tid string) (domain.Task, error) {
	var do taskDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(ctx, impl.docFilter(tid), bson.M{fieldOlds: 0}, &do)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return domain.Task{}, err
	}

	return do.toTask(), nil
}

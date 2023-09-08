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

const (
	fieldName = "name"
	fieldOlds = "olds"
)

type taskDO struct {
	Name string   `bson:"name"  json:"name"`
	Kind string   `bson:"kind"  json:"kind"`
	Addr string   `bson:"addr"  json:"addr"`
	Rule ruleDO   `bson:"rule"  json:"rule"`
	Olds []ruleDO `bson:"olds"  json:"olds"`
}

func (do *taskDO) toTask() domain.Task {
	return domain.Task{
		Name: do.Name,
		Kind: do.Kind,
		Addr: do.Kind,
		Rule: do.Rule.toRule(),
	}
}

type ruleDO struct {
	OnceOnly       bool   `bson:"once_only"          json:"once_only"`
	Desc           string `bson:"desc"               json:"desc"`
	CreatedAt      string `bson:"created_at"         json:"created_at"`
	PointsPerOnce  int    `bson:"points_per_once"    json:"points_per_once"`
	MaxPointsOfDay int    `bson:"max_points_of_day"  json:"max_points_of_day"`
}

func (do *ruleDO) toRule() domain.Rule {
	return domain.Rule{
		OnceOnly:       do.OnceOnly,
		Desc:           do.Desc,
		CreatedAt:      do.CreatedAt,
		PointsPerOnce:  do.PointsPerOnce,
		MaxPointsOfDay: do.MaxPointsOfDay,
	}
}

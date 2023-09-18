package repositoryadapter

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/points/domain"
)

const (
	fieldId   = "id"
	fieldOlds = "olds"
)

func totaskDO(t *domain.Task) taskDO {
	return taskDO{
		Id:    t.Id,
		Names: t.Names,
		Kind:  t.Kind,
		Addr:  t.Addr,
		Rule:  toruleDO(&t.Rule),
		Olds:  []ruleDO{},
	}
}

func toruleDO(r *domain.Rule) ruleDO {
	return ruleDO{
		Descs:          r.Descs,
		CreatedAt:      r.CreatedAt,
		OnceOnly:       r.OnceOnly,
		PointsPerOnce:  r.PointsPerOnce,
		MaxPointsOfDay: r.MaxPointsOfDay,
	}
}

// taskDO
type taskDO struct {
	Id    string            `bson:"id"    json:"id"`
	Names map[string]string `bson:"name"  json:"name"`
	Kind  string            `bson:"kind"  json:"kind"`
	Addr  string            `bson:"addr"  json:"addr"`
	Rule  ruleDO            `bson:"rule"  json:"rule"`
	Olds  []ruleDO          `bson:"olds"  json:"olds"`
}

func (do *taskDO) doc() (bson.M, error) {
	return genDoc(do)
}

func (do *taskDO) toTask() domain.Task {
	return domain.Task{
		Id:    do.Id,
		Names: do.Names,
		Kind:  do.Kind,
		Addr:  do.Addr,
		Rule:  do.Rule.toRule(),
	}
}

// ruleDO
type ruleDO struct {
	Descs          map[string]string `bson:"desc"               json:"desc"`
	CreatedAt      string            `bson:"created_at"         json:"created_at"`
	OnceOnly       bool              `bson:"once_only"          json:"once_only"`
	PointsPerOnce  int               `bson:"points_per_once"    json:"points_per_once"`
	MaxPointsOfDay int               `bson:"max_points_of_day"  json:"max_points_of_day"`
}

func (do *ruleDO) toRule() domain.Rule {
	return domain.Rule{
		OnceOnly:       do.OnceOnly,
		Descs:          do.Descs,
		CreatedAt:      do.CreatedAt,
		PointsPerOnce:  do.PointsPerOnce,
		MaxPointsOfDay: do.MaxPointsOfDay,
	}
}

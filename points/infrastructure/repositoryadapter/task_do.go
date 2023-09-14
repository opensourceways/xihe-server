package repositoryadapter

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/points/domain"
)

const (
	fieldName = "name"
	fieldOlds = "olds"
)

func totaskDO(t *domain.Task) taskDO {
	return taskDO{
		Name: t.Name,
		Kind: t.Kind,
		Addr: t.Addr,
		Rule: toruleDO(&t.Rule),
		Olds: []ruleDO{},
	}
}

func toruleDO(r *domain.Rule) ruleDO {
	return ruleDO{
		OnceOnly:       r.OnceOnly,
		Desc:           r.Desc,
		CreatedAt:      r.CreatedAt,
		PointsPerOnce:  r.PointsPerOnce,
		MaxPointsOfDay: r.MaxPointsOfDay,
	}
}

// taskDO
type taskDO struct {
	Name string   `bson:"name"  json:"name"`
	Kind string   `bson:"kind"  json:"kind"`
	Addr string   `bson:"addr"  json:"addr"`
	Rule ruleDO   `bson:"rule"  json:"rule"`
	Olds []ruleDO `bson:"olds"  json:"olds"`
}

func (do *taskDO) doc() (bson.M, error) {
	return genDoc(do)
}

func (do *taskDO) toTask() domain.Task {
	return domain.Task{
		Name: do.Name,
		Kind: do.Kind,
		Addr: do.Addr,
		Rule: do.Rule.toRule(),
	}
}

// ruleDO
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

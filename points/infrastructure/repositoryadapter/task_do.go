package repositoryadapter

import "github.com/opensourceways/xihe-server/points/domain"

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

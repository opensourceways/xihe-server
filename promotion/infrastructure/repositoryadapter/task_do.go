package repositoryadapter

import (
	"github.com/opensourceways/xihe-server/promotion/domain"
)

const (
	fieldId = "id"
)

type taskDO struct {
	Id    string            `bson:"user"  json:"user"`
	Names map[string]string `bson:"names" json:"names"`
	Rule  ruleDO            `bson:"rule"  json:"rule"`
}

func (do *taskDO) toTask() (task domain.Task, err error) {
	names, err := domain.NewSentence(
		do.Names[domain.FieldEN], do.Names[domain.FieldZH])
	if err != nil {
		return
	}

	rule, err := do.Rule.toRule()
	if err != nil {
		return
	}

	return domain.Task{
		Id:    do.Id,
		Names: names,
		Rule:  rule,
	}, nil
}

type ruleDO struct {
	Descs     map[string]string `bson:"descs"      json:"descs"`
	CreatedAt string            `bson:"created_at" json:"created_at"`
	MaxPoints int               `bson:"max_points" json:"max_points"`
}

func (do *ruleDO) toRule() (domain.Rule, error) {
	descs, err := domain.NewSentence(
		do.Descs[domain.FieldEN], do.Descs[domain.FieldZH])
	if err != nil {
		return domain.Rule{}, err
	}

	return domain.Rule{
		Descs:     descs,
		CreatedAt: do.CreatedAt,
		MaxPoints: do.MaxPoints,
	}, nil

}

package repository

import "github.com/opensourceways/xihe-server/points/domain"

type PointsRule interface {
	FindPointsItemRules() ([]domain.PointsItemRule, error)
	PointsOfDay() (int, error)
}

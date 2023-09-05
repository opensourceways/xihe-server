package repository

import "github.com/opensourceways/xihe-server/point/domain"

type PointRule interface {
	FindAll() ([]domain.PointItemRule, error)
	PointsOfDay() (int, error)
}

package repository

import "github.com/opensourceways/xihe-server/point/domain"

type PointRule interface {
	Find(t string) (domain.PointRule, error)
}

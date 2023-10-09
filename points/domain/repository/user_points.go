package repository

import (
	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/domain"
)

type UserPoints interface {
	SavePointsItem(*domain.UserPoints, *domain.PointsItem) error
	Find(account common.Account, date string) (domain.UserPoints, error)
	FindAll(account common.Account) (domain.UserPoints, error)
}

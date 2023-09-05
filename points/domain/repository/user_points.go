package repository

import (
	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/domain"
)

type UserPointsDetails struct {
	Total int
	Items []domain.PointsItem
}

func (details *UserPointsDetails) DetailNum() int {
	n := 0
	for i := range details.Items {
		n += len(details.Items[i].Details)
	}

	return n
}

type UserPoints interface {
	SavePointsItem(*domain.UserPoints, *domain.PointsItem) error
	Find(account common.Account, date string) (domain.UserPoints, error)
	FindAll(account common.Account) (UserPointsDetails, error)
}

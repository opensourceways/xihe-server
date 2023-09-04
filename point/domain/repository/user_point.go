package repository

import (
	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/point/domain"
)

type UserPointDetails struct {
	Total int
	Items []domain.PointItem
}

func (details *UserPointDetails) DetailNum() int {
	n := 0
	for i := range details.Items {
		n += len(details.Items[i].Details)
	}

	return n
}

type UserPoint interface {
	SavePointItem(*domain.UserPoint, *domain.PointItem) error
	Find(account common.Account, date string) (domain.UserPoint, error)
	FindAll(account common.Account) (UserPointDetails, error)
}

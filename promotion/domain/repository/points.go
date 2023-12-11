package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
)

type Points interface {
	Update(user types.Account, item domain.Item, version int) error
	Find(types.Account) (domain.UserPoints, error)
	FindAll() ([]domain.UserPoints, error)
}

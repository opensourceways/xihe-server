package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
)

type Points interface {
	Update(user types.Account, item domain.Item, version int) error
	Find(user types.Account, promotionid string) (domain.UserPoints, error)
	FindAll(promotionid string) ([]domain.UserPoints, error)
}

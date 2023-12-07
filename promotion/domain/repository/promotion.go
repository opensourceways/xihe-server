package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
)

type Promotion interface {
	Find(promotionid string) (domain.Promotion, error)
	UserRegister(promotionid string, user types.Account) error
}

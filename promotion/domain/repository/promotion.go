package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
)

type PromotionRepo struct {
	domain.Promotion
	Version int
}

type Promotion interface {
	Find(promotionid string) (PromotionRepo, error)
	FindAll() ([]PromotionRepo, error)
	UserRegister(promotionid string, user types.Account, version int) error
}

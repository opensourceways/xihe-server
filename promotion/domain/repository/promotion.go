package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
)

const (
	SortFieldStartTime = "start_time"
	SortFieldPriority  = "priority"
	SortAsc            = "asc"
	SortDesc           = "desc"
)

type Promotion interface {
	FindById(string) (domain.Promotion, error)
	FindAll() ([]domain.Promotion, error)
	UserRegister(promotionid string, user types.Account, origin domain.Origin, version int) error
	FindByCustom(*PromotionsQuery) ([]domain.Promotion, error)
	Count(*PromotionsQuery) (int64, error)
}

type PromotionsQuery struct {
	domain.Promotion
	Status domain.PromotionStatus
	Offset int64
	Limit  int64
	Sort   [][2]string
}

package app

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
)

type PromotionService interface {
	GetPromotion(*PromotionCmd) (PromotionDTO, error)
	GetUserRegisterPromotion(*types.Account) ([]PromotionDTO, error)
}

func NewPromotionService(
	repo repository.Promotion,
) PromotionService {
	return &promotionService{
		repo: repo,
	}
}

type promotionService struct {
	repo repository.Promotion
}

func (s *promotionService) GetPromotion(cmd *PromotionCmd) (dto PromotionDTO, err error) {
	// find promotion
	p, err := s.repo.Find(cmd.promotionId)
	if err != nil {
		return
	}

	dto.toDTO(&p, &cmd.User)

	return
}

func (s *promotionService) GetUserRegisterPromotion(user *types.Account) (dtos []PromotionDTO, err error) {
	// find all promotions
	ps, err := s.repo.FindAll()

	// generate promotion dtos
	for i := range ps {
		if ps[i].HasRegister(*user) {
			dto := PromotionDTO{}
			dto.toDTO(&ps[i], user)
			dtos = append(dtos, dto)
		}
	}

	return
}

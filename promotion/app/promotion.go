package app

import "github.com/opensourceways/xihe-server/promotion/domain/repository"

type PromotionService interface {
	GetPromotion(*PromotionCmd) (PromotionDTO, error)
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

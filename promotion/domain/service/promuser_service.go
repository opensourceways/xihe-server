package service

import (
	"fmt"

	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain/user"
)

type PromotionUserService interface {
	Register(promotionid string, origin domain.Origin, ur *domain.UserRegistration) error
}

type promotionUserService struct {
	user user.User
	repo repository.Promotion
}

func NewPromotionUserService(
	user user.User,
	repo repository.Promotion,
) PromotionUserService {
	return &promotionUserService{
		user: user,
		repo: repo,
	}
}

func (s *promotionUserService) Register(pid string, origin domain.Origin, ur *domain.UserRegistration) error {
	// get promotion version
	p, err := s.repo.FindById(pid)
	if err != nil {
		return fmt.Errorf("find promotion error: %w", err)
	}

	// register promotion
	if err := s.repo.UserRegister(pid, ur.Account, origin, p.Version); err != nil {
		return fmt.Errorf("register promotion error: %w", err)
	}

	// update registration
	if err := s.user.UpdateRegister(ur); err != nil {
		return fmt.Errorf("update registration error: %w", err)
	}

	return nil
}

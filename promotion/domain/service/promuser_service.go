package service

import (
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain/user"
)

type PromotionUserService interface {
	Register(promotionid string, ur *domain.UserRegistration) error
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

func (s *promotionUserService) Register(pid string, ur *domain.UserRegistration) error {
	// register promotion
	if err := s.repo.UserRegister(pid, ur.Account); err != nil {
		return err
	}

	// update registration
	if err := s.user.UpdateRegister(ur); err != nil {
		return err
	}

	return nil
}

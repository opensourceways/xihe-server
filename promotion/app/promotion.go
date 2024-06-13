package app

import (
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain/service"
)

type PromotionService interface {
	GetPromotion(*PromotionCmd) (PromotionDTO, error)
	GetUserRegisterPromotion(*types.Account) ([]PromotionDTO, string, error)
	UserRegister(*UserRegistrationCmd) (code string, err error)
}

func NewPromotionService(
	service service.PromotionUserService,
	ptservice service.PointsTaskService,
	repo repository.Promotion,
) PromotionService {
	return &promotionService{
		service:   service,
		ptservice: ptservice,
		repo:      repo,
	}
}

type promotionService struct {
	service   service.PromotionUserService
	ptservice service.PointsTaskService
	repo      repository.Promotion
}

func (s *promotionService) GetPromotion(cmd *PromotionCmd) (dto PromotionDTO, err error) {
	// find promotion
	p, err := s.repo.Find(cmd.promotionId)
	if err != nil {
		return
	}

	// find user point
	up, err := s.ptservice.Find(cmd.User, cmd.promotionId)
	if err != nil {
		return
	}

	dto.toDTO(&p.Promotion, &cmd.User, up.Total)

	return
}

func (s *promotionService) GetUserRegisterPromotion(user *types.Account) (
	dtos []PromotionDTO, code string, err error,
) {
	// find all promotions
	ps, err := s.repo.FindAll()
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			code = errorUserRegistrationExists
		}

		return
	}

	for i := range ps {
		if ps[i].HasRegister(*user) {
			// get total of user points
			var total int
			up, err := s.ptservice.Find(*user, ps[i].Id)
			if err != nil {
				if !repoerr.IsErrorResourceNotExists(err) {
					return nil, "", err
				}
			} else {
				total = up.Total
			}

			// gen promotion dtos
			dto := PromotionDTO{}
			dto.toDTO(&ps[i].Promotion, user, total)
			dtos = append(dtos, dto)
		}
	}

	return
}

func (s *promotionService) UserRegister(cmd *UserRegistrationCmd) (code string, err error) {
	if err = s.service.Register(cmd.PromotionId, cmd.Origin, &cmd.UserRegistration); err != nil {
		if repoerr.IsErrorDuplicateCreating(err) {
			code = errorUserRegistrationExists
		}
	}

	return
}

package app

import (
	"errors"

	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain/service"
)

type PromotionService interface {
	GetPromotion(*PromotionCmd) (PromotionDTO, error)
	GetUserRegisterPromotion(types.Account) ([]PromotionDTO, string, error)
	UserRegister(*UserRegistrationCmd) (code string, err error)
	Get(*PromotionCmd) (PromotionDTO, error)
	List(*ListPromotionsCmd) (PromotionsDTO, error)
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
	p, err := s.repo.FindById(cmd.Id)
	if err != nil {
		return
	}

	// find user point
	up, err := s.ptservice.Find(cmd.User, cmd.Id)
	if err != nil {
		return
	}

	err = dto.toDTO(&p, cmd.User, up.Total)

	return
}

func (s *promotionService) GetUserRegisterPromotion(user types.Account) (
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
		if ps[i].HasRegister(user) {
			// get total of user points
			var total int
			up, err := s.ptservice.Find(user, ps[i].Id)
			if err != nil {
				if !repoerr.IsErrorResourceNotExists(err) {
					return nil, "", err
				}
			} else {
				total = up.Total
			}

			// gen promotion dtos
			dto := PromotionDTO{}
			err = dto.toDTO(&ps[i], user, total)
			if err != nil {
				return nil, "", err
			}
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

func (s *promotionService) Get(cmd *PromotionCmd) (PromotionDTO, error) {
	dto := PromotionDTO{}

	p, err := s.repo.FindById(cmd.Id)
	if err != nil {
		return dto, err
	}

	err = dto.toDTO(&p, cmd.User, 0)

	return dto, err
}

func (s *promotionService) List(cmd *ListPromotionsCmd) (PromotionsDTO, error) {
	promotionsDTO := PromotionsDTO{
		Items: make([]PromotionDTO, 0),
	}

	query := cmd.toPromotionsQuery()

	total, err := s.repo.Count(query)
	if err != nil {
		return promotionsDTO, err
	}

	if total == 0 {
		return promotionsDTO, nil
	}

	if query.Offset >= total {
		return promotionsDTO, repoerr.NewExcendMaximumPageNumError(errors.New("excend the maximum page number"))
	}

	promotions, err := s.repo.FindByCustom(query)
	if err != nil {
		return promotionsDTO, err
	}

	promotionsDTO.Items = make([]PromotionDTO, 0, len(promotions))
	for i := range promotions {
		promotionDTO := PromotionDTO{}
		if err := promotionDTO.toDTO(&promotions[i], cmd.User, 0); err != nil {
			return promotionsDTO, err
		}
		promotionsDTO.Items = append(promotionsDTO.Items, promotionDTO)
	}
	promotionsDTO.Total = total

	return promotionsDTO, nil
}

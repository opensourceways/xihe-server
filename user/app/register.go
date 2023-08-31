package app

import (
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

type RegService interface {
	// register
	UpsertUserRegInfo(*UserRegisterInfoCmd) error
	GetUserRegInfo(domain.Account) (UserRegisterInfoDTO, error)
}

var _ RegService = (*regService)(nil)

func NewRegService(
	regRepo repository.UserReg,
) *regService {
	return &regService{
		regRepo: regRepo,
	}
}

type regService struct {
	regRepo repository.UserReg
}

func (s *regService) UpsertUserRegInfo(cmd *UserRegisterInfoCmd) (err error) {
	r := new(domain.UserRegInfo)
	cmd.toUserRegInfo(r)

	if err = s.regRepo.AddUserRegInfo(r); err != nil {
		if !repoerr.IsErrorDuplicateCreating(err) {
			return
		}

		u, err := s.regRepo.GetUserRegInfo(cmd.Account)
		if err != nil {
			return err
		}

		if err = s.regRepo.UpdateUserRegInfo(r, u.Version); err != nil {
			return err
		}
	}

	return nil
}

func (s *regService) GetUserRegInfo(user domain.Account) (dto UserRegisterInfoDTO, err error) {
	u, err := s.regRepo.GetUserRegInfo(user)
	if err != nil {
		return
	}

	dto.toUserRegInfoDTO(&u)

	return
}

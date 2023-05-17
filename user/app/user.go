package app

import (
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

type UserService interface {
	UpsertUserRegInfo(*UserRegisterInfoCmd) error
	GetUserRegInfo(types.Account) (UserRegisterInfoDTO, error)
}

var _ UserService = (*userService)(nil)

func NewUserService(
	userregRepo repository.UserReg,
) *userService {
	return &userService{
		userregRepo: userregRepo,
	}
}

type userService struct {
	userregRepo repository.UserReg
}

func (s *userService) UpsertUserRegInfo(cmd *UserRegisterInfoCmd) (err error) {
	r := new(domain.UserRegInfo)
	cmd.toUserRegInfo(r)

	if err = s.userregRepo.AddUserRegInfo(r); err != nil {
		if !repoerr.IsErrorDuplicateCreating(err) {
			return
		}

		u, err := s.userregRepo.GetUserRegInfo(cmd.Account)
		if err != nil {
			return err
		}

		if err = s.userregRepo.UpdateUserRegInfo(r, u.Version); err != nil {
			return err
		}
	}

	return nil
}

func (s *userService) GetUserRegInfo(user types.Account) (dto UserRegisterInfoDTO, err error) {
	u, err := s.userregRepo.GetUserRegInfo(user)
	if err != nil {
		return
	}

	dto.toUserRegInfoDTO(&u)

	return
}

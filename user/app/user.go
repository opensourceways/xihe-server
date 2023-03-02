package app

import (
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

type UserService interface {
	AddUserRegInfo(*UserRegisterInfoCmd) error
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

func (s *userService) AddUserRegInfo(cmd *UserRegisterInfoCmd) (err error) {
	r := new(domain.UserRegInfo)
	cmd.toUserRegInfo(r)

	return s.userregRepo.AddUserRegInfo(r)
}

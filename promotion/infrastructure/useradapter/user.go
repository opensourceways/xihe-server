package useradapter

import (
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/user"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
)

func NewUserAdapter(s userapp.RegService) user.User {
	return &userAdapter{s}
}

type userAdapter struct {
	s userapp.RegService
}

func (impl *userAdapter) UpdateRegister(ur *domain.UserRegistration) error {
	cmd := userapp.UserRegisterInfoCmd(
		userdomain.UserRegInfo(*ur),
	)

	return impl.s.UpsertUserRegInfo(&cmd)
}

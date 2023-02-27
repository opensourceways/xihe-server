package usercli

import (
	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/user"
	userApp "github.com/opensourceways/xihe-server/user/app"
)

func NewUserCli(c userApp.UserService) user.User {
	return &userImpl{c}
}

type userImpl struct {
	srv userApp.UserService
}

func (impl *userImpl) AddUserRegInfo(s *domain.Student) (err error) {
	return impl.srv.AddUserRegInfo(toUserRegisterInfoCmd(s))
}

func toUserRegisterInfoCmd(s *domain.Student) *userApp.UserRegisterInfoCmd {
	return &userApp.UserRegisterInfoCmd{
		Account:  s.Account.Account(),
		Name:     s.Name.StudentName(),
		City:     s.City.City(),
		Email:    s.Email.Email(),
		Phone:    s.Phone.Phone(),
		Identity: s.Identity.StudentIdentity(),
		Province: s.Province.Province(),
		Detail:   s.Detail,
	}
}

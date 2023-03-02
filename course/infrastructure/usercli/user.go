package usercli

import (
	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/user"
	userApp "github.com/opensourceways/xihe-server/user/app"
	userDomain "github.com/opensourceways/xihe-server/user/domain"
)

func NewUserCli(c userApp.UserService) user.User {
	return &userImpl{c}
}

type userImpl struct {
	srv userApp.UserService
}

func (impl *userImpl) AddUserRegInfo(s *domain.Student) (err error) {
	cmd := new(userApp.UserRegisterInfoCmd)
	if err = toUserRegisterInfoCmd(s, cmd); err != nil {
		return
	}

	return impl.srv.AddUserRegInfo(cmd)
}

func toUserRegisterInfoCmd(s *domain.Student, cmd *userApp.UserRegisterInfoCmd) (err error) {
	cmd.Account = s.Account

	if cmd.Name, err = userDomain.NewName(s.Name.StudentName()); err != nil {
		return
	}

	if cmd.City, err = userDomain.NewCity(s.City.City()); err != nil {
		return
	}

	if cmd.Email, err = userDomain.NewEmail(s.Email.Email()); err != nil {
		return
	}

	if cmd.Phone, err = userDomain.NewPhone(s.Phone.Phone()); err != nil {
		return
	}

	if cmd.Identity, err = userDomain.NewIdentity(s.Identity.StudentIdentity()); err != nil {
		return
	}

	if cmd.Province, err = userDomain.NewProvince(s.Province.Province()); err != nil {
		return
	}

	cmd.Detail = s.Detail

	return
}

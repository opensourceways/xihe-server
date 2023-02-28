package usersrv

import (
	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/user"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
)

func NewUserService(c userClient) user.User {
	return &userImpl{c}
}

type userImpl struct {
	cli userClient
}

func (impl *userImpl) AddUserRegInfo(s *domain.Student) (err error) {
	u := new(userdomain.UserRegInfo)
	if err = impl.toUserRegInfo(s, u); err != nil {
		return
	}

	return impl.cli.AddUserRegInfo(u)
}

func (impl *userImpl) toUserRegInfo(
	s *domain.Student,
	r *userdomain.UserRegInfo,
) (err error) {
	r.Account = s.Account

	r.Name, err = userdomain.NewName(s.Name.StudentName())
	if err != nil {
		return
	}

	r.City, err = userdomain.NewCity(s.City.City())
	if err != nil {
		return
	}

	r.Email, err = userdomain.NewEmail(s.Email.Email())
	if err != nil {
		return
	}

	r.Phone, err = userdomain.NewPhone(s.Phone.Phone())
	if err != nil {
		return
	}

	r.Identity, err = userdomain.NewIdentity(s.Identity.StudentIdentity())
	if err != nil {
		return
	}

	r.Province, err = userdomain.NewProvince(s.Province.Province())
	if err != nil {
		return
	}

	r.Detail = s.Detail

	return
}

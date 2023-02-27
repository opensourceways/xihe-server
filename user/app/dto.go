package app

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/user/domain"
)

type UserRegisterInfoCmd struct {
	Account  string
	Name     string
	City     string
	Email    string
	Phone    string
	Identity string
	Province string
	Detail   map[string]string
}

func (cmd *UserRegisterInfoCmd) toUserRegInfo(r *domain.UserRegInfo) (err error) {
	if r.Account, err = types.NewAccount(cmd.Account); err != nil {
		return
	}

	if r.Name, err = domain.NewName(cmd.Name); err != nil {
		return
	}

	if r.City, err = domain.NewCity(cmd.City); err != nil {
		return
	}

	if r.Email, err = domain.NewEmail(cmd.Email); err != nil {
		return
	}

	if r.Phone, err = domain.NewPhone(cmd.Phone); err != nil {
		return
	}

	if r.Identity, err = domain.NewIdentity(cmd.Identity); err != nil {
		return
	}

	if r.Province, err = domain.NewProvince(cmd.Province); err != nil {
		return
	}

	r.Detail = cmd.Detail

	return
}

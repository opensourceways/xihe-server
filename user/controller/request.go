package controller

import (
	"errors"

	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type UserRegistrationCmd domain.UserRegInfo

type ApplyRequest struct {
	Name     string            `json:"name"`
	City     string            `json:"city"`
	Email    string            `json:"email"`
	Phone    string            `json:"phone"`
	Identity string            `json:"identity"`
	Province string            `json:"province"`
	Detail   map[string]string `json:"detail"`
}

func (req *ApplyRequest) ToCmd(user types.Account) (cmd UserRegistrationCmd, err error) {
	if cmd.Name, err = domain.NewName(req.Name); err != nil {
		return
	}

	if cmd.City, err = domain.NewCity(req.City); err != nil {
		return
	}

	if cmd.Email, err = types.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Phone, err = domain.NewPhone(req.Phone); err != nil {
		return
	}

	if cmd.Identity, err = domain.NewIdentity(req.Identity); err != nil {
		return
	}

	if cmd.Province, err = domain.NewProvince(req.Province); err != nil {
		return
	}

	cmd.Detail = req.Detail
	cmd.Account = user

	err = cmd.Validate()

	return
}

func (cmd *UserRegistrationCmd) Validate() error {
	b := cmd.Account != nil &&
		cmd.Name != nil &&
		cmd.Email != nil &&
		cmd.Identity != nil

	if !b {
		return errors.New("invalid cmd")
	}

	for i := range cmd.Detail {
		if utils.StrLen(cmd.Detail[i]) > 20 {
			return errors.New("invalid detail")
		}
	}

	return nil
}

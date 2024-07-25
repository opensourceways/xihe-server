package controller

import (
	agreement "github.com/opensourceways/xihe-server/agreement/app"
	"github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/user/domain"
)

type UserAgreement struct {
	Type agreement.AgreementType `json:"type"`
}

type followingCreateRequest struct {
	Account string `json:"account" required:"true"`
}

type userDetail struct {
	*app.UserDTO

	Points     int  `json:"points"`
	IsFollower bool `json:"is_follower"`
}

type EmailCode struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (req *EmailCode) toCmd(user domain.Account) (cmd app.BindEmailCmd, err error) {
	if cmd.Email, err = domain.NewEmail(req.Email); err != nil {
		return
	}

	cmd.PassCode = req.Code
	cmd.User = user

	if cmd.PassWord, err = domain.NewPassword(apiConfig.DefaultPassword); err != nil {
		return
	}

	return
}

type EmailSend struct {
	Email string `json:"email"`
	Capt  string `json:"capt"`
}

func (req *EmailSend) toCmd(user domain.Account) (cmd app.SendBindEmailCmd, err error) {
	if cmd.Email, err = domain.NewEmail(req.Email); err != nil {
		return
	}

	cmd.User = user
	cmd.Capt = req.Capt

	return
}

func toCheckWhiteListCmd(user domain.Account, t string) (cmd app.UserWhiteListCmd, err error) {
	cmd.Account = user
	v, err := domain.NewWhiteListType(t)
	if err != nil {
		return
	}
	cmd.Type = v
	return
}

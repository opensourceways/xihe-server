package controller

import (
	agreement "github.com/opensourceways/xihe-server/agreement/app"
	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/user/domain"
	"golang.org/x/xerrors"
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

type ModifyInfo struct {
	AccountType string `json:"account_type"`
	OldAccount  string `json:"oldaccount"`
	OldCode     string `json:"oldcode"`
	Account     string `json:"account"`
	Code        string `json:"code"`
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

func (m *ModifyInfo) toCmd() (cmd app.CmdToModifyInfo, err error) {
	if m.Account == "" || m.OldAccount == "" || m.OldCode == "" || m.AccountType == "" || m.Code == "" {
		err = allerror.New(allerror.ErrorCodeUserNotFound, "",
			xerrors.Errorf("user update info empty"))
		return
	}
	cmd.Account = m.Account
	cmd.AccountType = m.AccountType
	cmd.OldCode = m.OldCode
	cmd.OldAccount = m.OldAccount
	cmd.Code = m.Code
	if cmd.Password, err = domain.NewPassword(apiConfig.DefaultPassword); err != nil {
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

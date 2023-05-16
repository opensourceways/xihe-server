package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/user/domain/login"
)

type EmailService interface {
	SendBindEmail(*SendBindEmailCmd) (string, error)
	VerifyBindEmail(*BindEmailCmd) (string, error)
}

func NewEmailService(
	auth authing.User,
	login login.Login,
	us UserService,
) EmailService {
	return &emailService{
		auth:  auth,
		login: login,
		us:    us,
	}
}

type emailService struct {
	auth  authing.User
	login login.Login
	us    UserService
}

func (s emailService) SendBindEmail(cmd *SendBindEmailCmd) (code string, err error) {
	return "", s.auth.SendBindEmail(cmd.Email.Email(), cmd.Capt)
}

func (s emailService) VerifyBindEmail(cmd *BindEmailCmd) (code string, err error) {
	info, err := s.login.GetAccessAndIdToken(cmd.User)
	if err != nil {
		return
	}

	if info.AccessToken == "" {
		code = errorNoAccessToken
		err = errors.New("cannot read access token")

		return
	}

	if err = s.auth.VerifyBindEmail(cmd.Email.Email(), cmd.PassCode, info.UserId); err != nil {
		return
	}

	// create platform account
	pfcmd := &CreatePlatformAccountCmd{
		email:    cmd.Email,
		account:  cmd.User,
		password: cmd.PassWord,
	}

	if err = s.us.CreatePlatformAccount(pfcmd); err != nil {
		return
	}

	return
}

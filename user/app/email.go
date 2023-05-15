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
) EmailService {
	return &emailService{
		auth:  auth,
		login: login,
	}
}

type emailService struct {
	auth  authing.User
	login login.Login
}

func (s emailService) SendBindEmail(cmd *SendBindEmailCmd) (code string, err error) {
	info, err := s.login.GetAccessAndIdToken(cmd.User)
	if err != nil {
		return
	}

	if info.AccessToken == "" {
		code = errorNoAccessToken
		err = errors.New("cannot read access token")

		return
	}

	return "", s.auth.SendBindEmail(info.AccessToken, cmd.Email.Email())
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

	return "", s.auth.VerifyBindEmail(info.AccessToken, cmd.Email.Email(), cmd.PassCode)
}

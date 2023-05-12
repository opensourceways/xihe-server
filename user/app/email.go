package app

import (
	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/login"
)

type EmailService interface {
	SendBindEmail(domain.Account) error
	VerifyBindEmail(domain.Account, string) error
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

func (s emailService) SendBindEmail(u domain.Account) (err error) {
	info, err := s.login.GetAccessAndIdToken(u)
	if err != nil {
		return
	}

	return s.auth.SendBindEmail(info.AccessToken)
}

func (s emailService) VerifyBindEmail(u domain.Account, code string) (err error) {
	info, err := s.login.GetAccessAndIdToken(u)
	if err != nil {
		return
	}

	return s.auth.VerifyBindEmail(info.AccessToken, code)
}

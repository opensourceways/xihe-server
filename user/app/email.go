package app

import (
	"errors"
	"github.com/opensourceways/xihe-server/user/domain"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/user/domain/login"
	messages "github.com/opensourceways/xihe-server/user/domain/message"
)

type EmailService interface {
	SendBindEmail(*SendBindEmailCmd) (string, error)
	VerifyBindEmail(*BindEmailCmd) (string, error)
}

func NewEmailService(
	auth authing.User,
	login login.Login,
	us UserService,
	producer messages.MessageProducer,
) EmailService {
	return &emailService{
		auth:     auth,
		login:    login,
		us:       us,
		producer: producer,
	}
}

type emailService struct {
	auth     authing.User
	login    login.Login
	us       UserService
	producer messages.MessageProducer
}

func (s emailService) SendBindEmail(cmd *SendBindEmailCmd) (code string, err error) {
	return s.auth.SendBindEmail(cmd.Email.Email(), cmd.Capt)
}

func (s emailService) VerifyBindEmail(cmd *BindEmailCmd) (code string, err error) {
	info, err := s.login.GetAccessAndIdToken(cmd.User)
	if err != nil {
		return
	}

	if info.UserId == "" {
		code = errorNoUserId
		err = errors.New("cannot read user id")

		return
	}

	if code, err = s.auth.VerifyBindEmail(cmd.Email.Email(), cmd.PassCode, info.UserId); err != nil {
		if !isCodeUserDuplicateBind(code) {
			return
		}

		// get authing email which saved in xihe
		var login app.LoginDTO
		login, err = s.login.GetAccessAndIdToken(cmd.User)
		if err != nil {
			return
		}

		// if bind email is authing email, bind email
		if login.Email != cmd.Email.Email() {
			return
		}

	}

	// create platform account
	pfcmd := &CreatePlatformAccountCmd{
		Email:    cmd.Email,
		Account:  cmd.User,
		Password: cmd.PassWord,
	}

	if err = s.us.NewPlatformAccountWithUpdate(pfcmd); err != nil {
		return
	}

	s.producer.SendBindEmailEvent(&domain.UserBindEmailEvent{
		Account: cmd.User,
		Email:   cmd.Email,
	})

	return
}

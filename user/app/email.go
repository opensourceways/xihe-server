package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/domain/repository"
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
		return
	}

	// create platform account
	pfcmd := &CreatePlatformAccountCmd{
		email:    cmd.Email,
		account:  cmd.User,
		password: cmd.PassWord,
	}

	dto, err := s.us.CreatePlatformAccount(pfcmd)
	if err != nil {
		return
	}

	// update user information
	updatecmd := &UpdatePlateformInfoCmd{
		PlatformInfoDTO: dto,
		User:            cmd.User,
		Email:           cmd.Email,
	}

	for i := 0; i <= 5; i++ {
		if err = s.us.UpdatePlateformInfo(updatecmd); err != nil {
			if !repository.IsErrorConcurrentUpdating(err) {
				return
			}
		} else {
			break
		}
	}

	return
}

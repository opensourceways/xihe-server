package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type LoginCreateCmd struct {
	Account domain.Account
	Info    string
	Email   domain.Email
	UserId  string
}

func (cmd *LoginCreateCmd) Validate() error {
	if b := cmd.Account != nil && cmd.Info != "" && cmd.Email != nil && cmd.UserId != ""; !b {
		return errors.New("invalid cmd of creating login")
	}

	return nil
}

func (cmd *LoginCreateCmd) toLogin() domain.Login {
	return domain.Login{
		Account: cmd.Account,
		Info:    cmd.Info,
		Email:   cmd.Email,
		UserId:  cmd.UserId,
	}
}

type LoginDTO struct {
	Info   string `json:"info"`
	Email  string `json:"email"`
	UserId string `json:"user_id"`
}

type LoginService interface {
	Create(*LoginCreateCmd) error
	Get(domain.Account) (LoginDTO, error)
	SignIn(domain.Account) error
}

func NewLoginService(repo repository.Login, sender message.UserSignedInMessageProducer) LoginService {
	return loginService{
		repo:   repo,
		sender: sender,
	}
}

type loginService struct {
	repo   repository.Login
	sender message.UserSignedInMessageProducer
}

func (s loginService) Create(cmd *LoginCreateCmd) error {
	v := cmd.toLogin()

	// new login
	return s.repo.Save(&v)
}

func (s loginService) Get(account domain.Account) (dto LoginDTO, err error) {
	v, err := s.repo.Get(account)
	if err != nil {
		return
	}

	s.toLoginDTO(&v, &dto)

	return
}

func (s loginService) toLoginDTO(u *domain.Login, dto *LoginDTO) {
	dto.Info = u.Info
	dto.Email = u.Email.Email()
	dto.UserId = u.UserId
}

func (s loginService) SignIn(account domain.Account) error {
	return s.sender.SendUserSignedIn(&domain.UserSignedInEvent{account})
}

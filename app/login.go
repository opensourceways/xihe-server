package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type LoginCreateCmd struct {
	Account domain.Account
	Info    string
}

func (cmd *LoginCreateCmd) Validate() error {
	if b := cmd.Account != nil && cmd.Info != ""; !b {
		return errors.New("invalid cmd of creating login")
	}

	return nil
}

func (cmd *LoginCreateCmd) toLogin() domain.Login {
	return domain.Login{
		Account: cmd.Account,
		Info:    cmd.Info,
	}
}

type LoginDTO struct {
	Info string `json:"info"`
}

type LoginService interface {
	Create(*LoginCreateCmd) error
	Get(domain.Account) (LoginDTO, error)
}

func NewLoginService(repo repository.Login) LoginService {
	return loginService{
		repo: repo,
	}
}

type loginService struct {
	repo repository.Login
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
}

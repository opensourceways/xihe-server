package logincli

import (
	typesapp "github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/login"
)

func NewLoginCli(c typesapp.LoginService) login.Login {
	return &loginImpl{c}
}

type loginImpl struct {
	s typesapp.LoginService
}

func (impl *loginImpl) GetAccessAndIdToken(u domain.Account) (dto typesapp.LoginDTO, err error) {
	return impl.s.Get(u)
}

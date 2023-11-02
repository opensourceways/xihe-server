package login

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/user/domain"
)

type Login interface {
	GetAccessAndIdToken(domain.Account) (app.LoginDTO, error)
}

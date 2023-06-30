package login

import (
	typesapp "github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/user/domain"
)

type Login interface {
	GetAccessAndIdToken(domain.Account) (typesapp.LoginDTO, error)
}

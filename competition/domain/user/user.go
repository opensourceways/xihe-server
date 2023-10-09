package user

import (
	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type User interface {
	AddUserRegInfo(*domain.Competitor) error
	GetUserRegInfo(types.Account) (domain.Competitor, error)
}

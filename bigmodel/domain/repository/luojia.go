package repository

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type LuoJia interface {
	Save(*domain.UserLuoJiaRecord) (domain.LuoJiaRecord, error)
	List(types.Account) ([]domain.LuoJiaRecord, error)
}

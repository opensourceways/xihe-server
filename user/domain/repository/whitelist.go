package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/user/domain"
)

type WhiteList interface {
	GetWhiteListInfo(types.Account, string) (domain.WhiteListInfo, error)
	FindByAccountAndWhitelistType(types.Account, []string) ([]domain.WhiteListInfo, error)
}

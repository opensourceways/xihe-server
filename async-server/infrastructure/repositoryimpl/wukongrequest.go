package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
)

func NewWuKongRequestRepo(cfg *Config) repository.WuKongRequest {
	return &wukongRequestRepoImpl{
		cli: pgsql.NewDBTable(cfg.Table.WukongRequest),
	}
}

type wukongRequestRepoImpl struct {
	cli pgsqlClient
}

func (impl *wukongRequestRepoImpl) HasNewRequest(time int64) (
	b bool, err error,
) {
	return
}

func (impl *wukongRequestRepoImpl) GetMultipleWuKongRequest(num int) (
	d []domain.WuKongRequest, err error,
) {
	return
}

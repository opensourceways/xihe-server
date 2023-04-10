package repositoryimpl

import (
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

func (impl *wukongRequestRepoImpl) GetNewRequest(time int64) (
	d []repository.WuKongTask, err error,
) {
	var twukong []TWukongTask

	impl.cli.DB().
		Where("created_at > ? AND status = ?", time, "waiting").
		Find(&twukong)

	d = make([]repository.WuKongTask, len(twukong))
	for i := range twukong {
		twukong[i].toWuKongTask(&d[i])
	}

	return
}

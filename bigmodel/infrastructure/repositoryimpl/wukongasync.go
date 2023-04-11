package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	types "github.com/opensourceways/xihe-server/domain"
)

func NewWuKongAsyncRepo(cfg *Config) repository.WuKongAsyncTask {
	return &wukongAsyncRepoImpl{
		cli: pgsql.NewDBTable(cfg.Table.WukongRequest),
	}
}

type wukongAsyncRepoImpl struct {
	cli pgsqlClient
}

func (impl *wukongAsyncRepoImpl) GetWaitingTaskRank(user types.Account, t commondomain.Time) (r int, err error) {
	var twukong []TWukongTask

	// 1. get all task before t
	err = impl.cli.DB().
		Where("created_at > ? and status = ?", t.Time(), "waiting").
		Find(&twukong).Error
	if err != nil {
		if impl.cli.IsRowNotFound(err) {
			err = repository.NewErrorResourceNotExists(err)

			return
		}

		return
	}

	// 2. is user in task
	f1 := func(v []TWukongTask) bool {
		for i := range twukong {
			if twukong[i].User == user.Account() {
				return true
			}
		}

		return false
	}

	if !f1(twukong) {
		return 0, nil
	}

	// 2. caculate rank
	f2 := func(v []TWukongTask) int {
		i := 1

		for j := range v {
			if v[j].CreatedAt <= t.Time() {
				i++
			}
		}

		return i
	}

	return f2(twukong), nil
}

func (impl *wukongAsyncRepoImpl) GetLastFinishedTask(user types.Account) (resp repository.WuKongTaskResp, err error) {
	var twukong TWukongTask

	filter := map[string]interface{}{
		fieldUserName: user.Account(),
	}

	order := "created_at DESC"

	if err = impl.cli.GetOrderOneRecord(filter, order, &twukong); err != nil {
		if impl.cli.IsRowNotFound(err) {
			err = repository.NewErrorResourceNotExists(err)

			return
		}

		return
	}

	err = twukong.toWuKongTaskResp(&resp)

	return
}

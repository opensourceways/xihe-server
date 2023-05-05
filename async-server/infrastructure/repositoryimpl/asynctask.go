package repositoryimpl

import (
	"fmt"

	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	types "github.com/opensourceways/xihe-server/domain"
)

func NewAsyncTaskRepo(cfg *Config) repository.AsyncTask {
	return &asyncTaskRepoImpl{
		cli: pgsql.NewDBTable(cfg.Table.AsyncTask),
	}
}

type asyncTaskRepoImpl struct {
	cli pgsqlClient
}

func (impl *asyncTaskRepoImpl) GetNewTask(taskType string, time int64) (
	d []repository.WuKongTask, err error,
) {
	var twukong []TAsyncTask

	impl.cli.DB().
		Where("task_type = ? AND created_at > ? AND status = ?", taskType, time, "waiting").
		Find(&twukong)

	d = make([]repository.WuKongTask, len(twukong))
	for i := range twukong {
		twukong[i].toWuKongTask(&d[i])
	}

	return
}

func (impl *asyncTaskRepoImpl) UpdateTask(resp *repository.WuKongResp) (err error) {

	v := NewTAsyncTask()
	v.toTAsyncTask(resp)

	filter := map[string]interface{}{
		fieldId: resp.WuKongTask.Id,
	}

	return impl.cli.UpdateRecord(filter, v) // TODO: don't cover old metadata
}

func (impl *asyncTaskRepoImpl) InsertTask(req *domain.WuKongRequest) error {

	v := NewTAsyncTask()
	v.toTWuKongTaskFromWuKongRequest(req)

	fmt.Printf("req: %+v\n", req)

	return impl.cli.Create(v)
}

func (impl *asyncTaskRepoImpl) GetWaitingTaskRank(user types.Account, t commondomain.Time, taskType string) (r int, err error) {
	var twukong []TAsyncTask

	// 1. get all task before t
	err = impl.cli.DB().
		Where("created_at > ? and status IN ? and task_type = ?", t.Time(), []string{"waiting", "running"}, taskType).
		Find(&twukong).Error
	if err != nil {
		if impl.cli.IsRowNotFound(err) {
			err = commonrepo.NewErrorResourceNotExists(err)

			return
		}

		return
	}

	// 2. is user in task
	f1 := func(v []TAsyncTask) bool {
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
	f2 := func(v []TAsyncTask) int {
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

func (impl *asyncTaskRepoImpl) GetLastFinishedTask(user types.Account, taskType string) (resp repository.WuKongResp, err error) {
	var twukong TAsyncTask

	filter := map[string]interface{}{
		fieldUserName: user.Account(),
		fieldTaskType: taskType,
		fieldStatus:   "finished",
	}

	order := "created_at DESC"

	if err = impl.cli.GetOrderOneRecord(filter, order, &twukong); err != nil {
		if impl.cli.IsRowNotFound(err) {
			err = commonrepo.NewErrorResourceNotExists(err)

			return
		}

		return
	}

	err = twukong.toWuKongTaskResp(&resp)

	return
}

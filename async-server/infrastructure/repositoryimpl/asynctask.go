package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
)

func NewWuKongRequestRepo(cfg *Config) repository.AsyncTask {
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

	var v TAsyncTask
	v.toTAsyncTask(resp)

	filter := map[string]interface{}{
		fieldId: resp.WuKongTask.Id,
	}

	return impl.cli.UpdateRecord(filter, v)
}

func (impl *asyncTaskRepoImpl) InsertTask(req *domain.WuKongRequest) error {

	v := new(TAsyncTask)
	v.toTWuKongTaskFromWuKongRequest(req)

	return impl.cli.Create(v)
}

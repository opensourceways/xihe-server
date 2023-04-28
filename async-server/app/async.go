package app

import (
	"github.com/opensourceways/xihe-server/async-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/async-server/domain/pool"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
)

type AsyncService interface {
	AsyncWuKong(taskType string, time int64) error
	AsyncWuKong4Img(taskType string, time int64) error
}

func NewAsyncService(
	bigmodel bigmodel.BigModel,
	pool pool.Pool,
	repo repository.AsyncTask,
) AsyncService {
	return &asyncService{
		bigmodel: bigmodel,
		pool:     pool,
		repo:     repo,
	}
}

type asyncService struct {
	bigmodel bigmodel.BigModel
	pool     pool.Pool
	repo     repository.AsyncTask
}

func (s *asyncService) AsyncWuKong(taskType string, time int64) (err error) {
	// 1. get waiting tasks before oder time
	var reqs []repository.WuKongTask
	if reqs, err = s.repo.GetNewTask(taskType, time); err != nil || len(reqs) == 0 { // TODO config
		return
	}

	// 2. get endpoint idle & idle worker
	ep, w := 0, 0
	if ep, err = s.bigmodel.GetIdleEndpoint(taskType); err != nil {
		return
	}

	w = s.pool.GetIdleWorker()

	// 3. check above
	f := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	if len(reqs) > f(ep, w) {
		reqs = reqs[:f(ep, w)]
	}

	// 4. do task in the goroutine pool
	var tasks pool.TaskList
	tasks.InitTaskList(reqs, s.bigmodel.WuKong)

	return s.pool.DoTasks(tasks)
}

func (s *asyncService) AsyncWuKong4Img(taskType string, time int64) (err error) {
	// 1. get waiting tasks before oder time
	var reqs []repository.WuKongTask
	if reqs, err = s.repo.GetNewTask(taskType, time); err != nil || len(reqs) == 0 { // TODO config
		return
	}

	// 2. get endpoint idle & idle worker
	ep, w := 0, 0
	if ep, err = s.bigmodel.GetIdleEndpoint(taskType); err != nil {
		return
	}

	w = s.pool.GetIdleWorker()

	// 3. check above
	f := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	if len(reqs) > f(ep, w) {
		reqs = reqs[:f(ep, w)]
	}

	// 4. do task in the goroutine pool
	var tasks pool.TaskList
	tasks.InitTaskListForWuKong4Img(reqs, s.bigmodel.WuKong4Img)

	return s.pool.DoTasks(tasks)
}

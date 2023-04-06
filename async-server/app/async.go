package app

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/async-server/domain/pool"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
)

type AsyncService interface {
	AsyncWuKong() error
}

func NewAsyncService(
	bigmodel bigmodel.BigModel,
	pool pool.Pool,
	repo repository.WuKongRequest,
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
	repo     repository.WuKongRequest
}

func (s *asyncService) AsyncWuKong() (err error) {
	// 1. get waiting tasks
	var reqs []domain.WuKongRequest
	if reqs, err = s.repo.GetMultipleWuKongRequest(8); err != nil || len(reqs) == 0 { // TODO config
		return
	}

	// 2. get endpoint idle & idle worker
	ep, w := 0, 0
	if ep, err = s.bigmodel.GetIdleEndpoint("wukong"); err != nil {
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

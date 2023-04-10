package pool

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
)

type TaskList []func()

type Pool interface {
	GetIdleWorker() int
	DoTasks(TaskList) error
}

func (r *TaskList) InitTaskList(reqs []repository.WuKongTask, f func(*domain.WuKongRequest) error) {
	*r = make(TaskList, len(reqs))

	for i := range reqs {
		([]func())(*r)[i] = func() {
			f(&reqs[i].WuKongRequest)
		}
	}
}

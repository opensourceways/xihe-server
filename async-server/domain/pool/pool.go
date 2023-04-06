package pool

import "github.com/opensourceways/xihe-server/async-server/domain"

type TaskList []func()

type Pool interface {
	GetIdleWorker() int
	DoTasks(TaskList) error
}

func (r *TaskList) InitTaskList(reqs []domain.WuKongRequest, f func(*domain.WuKongRequest) error) {
	*r = make(TaskList, len(reqs))

	for i := range reqs {
		([]func())(*r)[i] = func() {
			f(&reqs[i])
		}
	}
}

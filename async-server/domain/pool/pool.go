package pool

import (
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
)

type TaskList []func()

type Pool interface {
	GetIdleWorker() int
	DoTasks(TaskList) error
}

func (r *TaskList) InitTaskList(reqs []repository.WuKongTask, f func(*repository.WuKongTask) error) {
	*r = make(TaskList, len(reqs))

	// build new function with new address
	funcBuild := func(i int) func() {
		return func() {
			f(&reqs[i])
		}
	}

	for i := range reqs {
		([]func())(*r)[i] = funcBuild(i)
	}
}

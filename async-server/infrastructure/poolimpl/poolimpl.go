package poolimpl

import (
	"github.com/panjf2000/ants"

	"github.com/opensourceways/xihe-server/async-server/domain/pool"
)

func NewPoolImpl() pool.Pool {
	return &poolImpl{
		p: gpool,
	}
}

type poolImpl struct {
	p *ants.Pool
}

func (impl *poolImpl) GetIdleWorker() int {
	return impl.p.Free()
}

func (impl *poolImpl) DoTasks(tasks pool.TaskList) error {
	for i := range tasks {
		if err := impl.p.Submit(tasks[i]); err != nil {
			return err
		}

	}

	return nil
}

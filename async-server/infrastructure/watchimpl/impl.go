package watchimpl

import (
	"fmt"
	"sync"
	"time"

	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	"github.com/sirupsen/logrus"
)

type Watcher struct {
	repo repository.WuKongRequest

	handles map[string]func(int64) error
	wg      sync.WaitGroup
}

func NewWather(
	repo repository.WuKongRequest,
	handles map[string]func(int64) error,
) *Watcher {
	return &Watcher{
		repo:    repo,
		handles: handles,
	}
}

func (w *Watcher) watchRequset() (err error) {
	logrus.Debug("start watching request")

	const swapTime = time.Second * 30

	t := time.NewTicker(swapTime)
	defer t.Stop()

	for now := range t.C {

		for bname := range w.handles {
			w.wg.Add(1)
			go w.work(bname, now.Unix() - int64(swapTime))
		}
	}

	return
}

func (w *Watcher) work(bname string, time int64) (err error) {
	defer w.wg.Done()

	if v, ok := w.handles[bname]; !ok {
		return fmt.Errorf("internal error, cannot found the bigmodel name:%s", bname)
	} else {
		v(time)
	}

	return nil
}

func (w *Watcher) Run() {

	w.watchRequset()

}

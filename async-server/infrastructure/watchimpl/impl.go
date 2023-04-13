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
	cfg     Config
	timer   *time.Ticker
	wg      sync.WaitGroup
}

func NewWather(
	cfg Config,
	repo repository.WuKongRequest,
	handles map[string]func(int64) error,
) *Watcher {

	return &Watcher{
		repo:    repo,
		timer:   time.NewTicker(time.Duration(cfg.Time.TriggerTime) * time.Second),
		handles: handles,
		cfg:     cfg,
	}
}

func (w *Watcher) watchRequset() (err error) {
	logrus.Debug("start watching request")

	for now := range w.timer.C {

		for bname := range w.handles {
			w.wg.Add(1)
			go w.work(bname, now.Add(-time.Duration(w.cfg.Time.ScanTime)*time.Second).Unix())
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

func (w *Watcher) Exit() {
	w.timer.Stop()

	w.wg.Wait()
}

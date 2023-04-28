package repository

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
)

type WuKongTask struct {
	domain.WuKongRequest

	Id     uint64
	Status domain.TaskStatus
}

type WuKongResp struct {
	WuKongTask

	Links domain.Links
}

type AsyncTask interface {
	GetNewTask(taskType string, time int64) ([]WuKongTask, error)
	UpdateTask(*WuKongResp) error
	InsertTask(*domain.WuKongRequest) error
}

func (r *WuKongTask) SetDefaultStatusWuKongTask(req *domain.WuKongRequest) {
	r.Status, _ = domain.NewTaskStatus("waiting")
	r.WuKongRequest = *req
}

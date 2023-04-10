package repository

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
)

type WuKongTask struct {
	domain.WuKongRequest

	Status domain.TaskStatus
}

type WuKongRequest interface {
	GetNewRequest(time int64) ([]WuKongTask, error)
}

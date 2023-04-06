package repository

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
)

type WuKongRequest interface {
	HasNewRequest(time int64) (bool, error)
	GetMultipleWuKongRequest(num int) ([]domain.WuKongRequest, error)
}

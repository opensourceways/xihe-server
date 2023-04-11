package bigmodel

import (
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
)

type BigModel interface {
	GetIdleEndpoint(bid string) (int, error)
	WuKong(*repository.WuKongTask) error
}

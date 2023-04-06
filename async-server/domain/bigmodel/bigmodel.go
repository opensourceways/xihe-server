package bigmodel

import "github.com/opensourceways/xihe-server/async-server/domain"

type BigModel interface {
	GetIdleEndpoint(bid string) (int, error)
	WuKong(*domain.WuKongRequest) error
}

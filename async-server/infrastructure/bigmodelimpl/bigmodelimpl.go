package bigmodelimpl

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/bigmodel"
	bigmodelapp "github.com/opensourceways/xihe-server/bigmodel/app"
)

func NewBigModelImpl(s bigmodelapp.AsyncBigModelService) bigmodel.BigModel {
	return &bigmodelImpl{
		srv: s,
	}
}

type bigmodelImpl struct {
	srv bigmodelapp.AsyncBigModelService
}

func (impl *bigmodelImpl) GetIdleEndpoint(bid string) (
	c int, err error,
) {
	return impl.srv.GetIdleEndpoint(bid)
}

func (impl *bigmodelImpl) WuKong(d *domain.WuKongRequest) (err error) {
	cmd := bigmodelapp.WuKongCmd{
		Style: d.Style,
		Desc:  d.Desc,
	}

	return impl.srv.WuKong(d.Id, d.User, &cmd)
}

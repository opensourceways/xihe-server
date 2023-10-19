package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

const (
	fieldUserName = "username"
)

func (a *dApiInfo) toApiInfo(d *domain.ApiInfo) (err error) {
	d.Id = a.Id
	d.Name = a.Name
	if d.Doc, err = types.NewURL(a.Doc); err != nil {
		return
	}
	d.Endpoint = a.Endpoint
	return
}

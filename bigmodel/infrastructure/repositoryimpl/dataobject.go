package repositoryimpl

import (
	asyncdomain "github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

const (
	fieldUserName = "username"
)

func (table *TWukongTask) toWuKongTaskResp(p *repository.WuKongTaskResp) (err error) {

	if p.Links, err = asyncdomain.NewLinks(table.Links); err != nil {
		return
	}

	if p.User, err = types.NewAccount(table.User); err != nil {
		return
	}

	if p.Desc, err = domain.NewWuKongPictureDesc(table.Desc); err != nil {
		return
	}

	if p.CreatedAt, err = commondomain.NewTime(table.CreatedAt); err != nil {
		return
	}

	if p.Status, err = asyncdomain.NewTaskStatus(table.Status); err != nil {
		return
	}

	p.Id = table.Id
	p.Style = table.Style

	return
}

func (a *dApiInfo) toApiInfo(d *domain.ApiInfo) (err error) {
	d.Id = a.Id
	d.Name = a.Name
	if d.Doc, err = types.NewURL(a.Doc); err != nil {
		return
	}
	d.Endpoint = a.Endpoint
	return
}

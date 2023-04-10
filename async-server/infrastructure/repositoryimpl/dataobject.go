package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

func (table *TWukongTask) toWuKongTask(p *repository.WuKongTask) (err error) {

	if p.User, err = types.NewAccount(table.User); err != nil {
		return
	}

	if p.Desc, err = bigmodeldomain.NewWuKongPictureDesc(table.Desc); err != nil {
		return
	}

	if p.CreatedAt, err = commondomain.NewTime(table.CreatedAt); err != nil {
		return
	}

	if p.Status, err = domain.NewTaskStatus(table.Status); err != nil {
		return
	}

	p.Id = table.Id
	p.Style = table.Style

	return
}

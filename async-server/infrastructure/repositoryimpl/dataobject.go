package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

const (
	fieldId = "id"
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

func (table *TWukongTask) toTWuKongTaskFromWuKongRequest(req *domain.WuKongRequest) {

	task := new(repository.WuKongTask)
	task.SetDefaultStatusWuKongTask(req)

	table.toTWuKongTaskFromWuKongTask(task)
}

func (table *TWukongTask) toTWuKongTaskFromWuKongTask(task *repository.WuKongTask) {

	*table = TWukongTask{
		Id: task.Id,
	}

	if task.User != nil {
		table.User = task.User.Account()
	}

	if task.Style != "" {
		table.Style = task.Style
	}

	if task.Desc != nil {
		table.Desc = task.Desc.WuKongPictureDesc()
	}

	if task.Status != nil {
		table.Status = task.Status.TaskStatus()
	}

	if task.CreatedAt != nil {
		table.CreatedAt = task.CreatedAt.Time()
	}

}

func (table *TWukongTask) toTWuKongTask(resp *repository.WuKongResp) {

	table.toTWuKongTaskFromWuKongTask(&resp.WuKongTask)

	if resp.Links != nil {
		table.Links = resp.Links.StringLinks()
	}

}

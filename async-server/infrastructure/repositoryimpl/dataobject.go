package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

const (
	fieldId       = "id"
	fieldUserName = "username"
	fieldTaskType = "task_type"
	fieldStatus   = "status"
)

func (table *TAsyncTask) toWuKongTask(p *repository.WuKongTask) (err error) {

	if p.User, err = types.NewAccount(table.User); err != nil {
		return
	}

	if table.MetaData["desc"] != nil {
		if p.Desc, err = bigmodeldomain.NewWuKongPictureDesc(table.MetaData["desc"].(string)); err != nil {
			return
		}
	}

	if p.CreatedAt, err = commondomain.NewTime(table.CreatedAt); err != nil {
		return
	}

	if p.Status, err = domain.NewTaskStatus(table.Status); err != nil {
		return
	}

	if table.TaskType != "" {
		if p.TaskType, err = domain.NewTaskType(table.TaskType); err != nil {
			return
		}
	}

	if table.MetaData["style"] != nil {
		p.Style = table.MetaData["style"].(string)
	}

	p.Id = table.Id

	return
}

func (table *TAsyncTask) toWuKongTaskResp(p *repository.WuKongResp) (err error) {
	table.toWuKongTask(&p.WuKongTask)

	if table.MetaData["links"] != nil {
		if p.Links, err = domain.NewLinks(table.MetaData["links"].(string)); err != nil {
			return
		}
	}

	return
}

func (table *TAsyncTask) toTWuKongTaskFromWuKongRequest(req *domain.WuKongRequest) {

	task := new(repository.WuKongTask)
	task.SetDefaultStatusWuKongTask(req)

	table.toTAsyncTaskFromWuKongTask(task)
}

func (table *TAsyncTask) toTAsyncTaskFromWuKongTask(task *repository.WuKongTask) {

	if task.User != nil {
		table.User = task.User.Account()
	}

	if task.TaskType != nil {
		table.TaskType = task.TaskType.TaskType()
	}

	if task.Style != "" {
		table.MetaData["style"] = task.Style
	}

	if task.Desc != nil {
		table.MetaData["desc"] = task.Desc.WuKongPictureDesc()
	}

	if task.Status != nil {
		table.Status = task.Status.TaskStatus()
	}

	if task.CreatedAt != nil {
		table.CreatedAt = task.CreatedAt.Time()
	}

	table.Id = task.Id

}

func (table *TAsyncTask) toTAsyncTask(resp *repository.WuKongResp) {

	table.toTAsyncTaskFromWuKongTask(&resp.WuKongTask)

	if resp.Links != nil {
		table.MetaData["links"] = resp.Links.StringLinks()
	}

}

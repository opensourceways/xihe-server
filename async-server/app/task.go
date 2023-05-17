package app

import (
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	commondomain "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type TaskService interface {
	GetWaitingTaskRank(types.Account, commondomain.Time, string) (int, error)
	GetLastFinishedTask(types.Account, string) (repository.WuKongResp, error)
}

func NewTaskService(
	repo repository.AsyncTask,
) TaskService {
	return &taskService{
		repo: repo,
	}
}

type taskService struct {
	repo repository.AsyncTask
}

func (s *taskService) GetWaitingTaskRank(user types.Account, time commondomain.Time, taskType string) (rank int, err error) {
	return s.repo.GetWaitingTaskRank(user, time, taskType)
}

func (s *taskService) GetLastFinishedTask(user types.Account, taskType string) (resp repository.WuKongResp, err error) {
	return s.repo.GetLastFinishedTask(user, taskType)
}

package app

import (
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
)

type TaskAppService interface {
	Add(t *domain.Task) error
}

func NewTaskAppService(
	repo repository.Task,
) *taskAppService {
	return &taskAppService{
		repo: repo,
	}
}

type taskAppService struct {
	repo repository.Task
}

func (s *taskAppService) Add(t *domain.Task) error {
	return s.repo.Add(t)
}

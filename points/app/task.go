package app

import (
	common "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain/service"
)

type TaskAppService interface {
	Add(t *domain.Task) error
	Doc(lang common.Language) (dto TaskDocDTO, err error)
}

func NewTaskAppService(
	s service.TaskService,
	repo repository.Task,
) *taskAppService {
	return &taskAppService{
		ts:   s,
		repo: repo,
	}
}

type taskAppService struct {
	ts   service.TaskService
	repo repository.Task
}

func (s *taskAppService) Add(t *domain.Task) error {
	return s.repo.Add(t)
}

func (s *taskAppService) Doc(lang common.Language) (dto TaskDocDTO, err error) {
	v, err := s.ts.Doc(lang)
	if err == nil {
		dto.Content = string(v)
	}

	return
}

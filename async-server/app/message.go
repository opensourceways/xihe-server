package app

import (
	"github.com/opensourceways/xihe-server/async-server/domain"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
)

type AsyncMessageService interface {
	UpdateWuKongTask(*repository.WuKongResp) error
	CreateWuKongTask(*domain.WuKongRequest) error
}

func NewAsyncMessageService(
	repo repository.AsyncTask,
) AsyncMessageService {
	return &asyncMessageService{
		repo: repo,
	}
}

type asyncMessageService struct {
	repo repository.AsyncTask
}

func (s *asyncMessageService) CreateWuKongTask(d *domain.WuKongRequest) error {
	return s.repo.InsertTask(d)
}

func (s *asyncMessageService) UpdateWuKongTask(resp *repository.WuKongResp) error {
	return s.repo.UpdateTask(resp)
}

package app

import (
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
)

type UserPointsAppMessageService interface {
	AddPointsItem(cmd *CmdToAddPointsItem) error
}

func NewUserPointsAppMessageService(
	tr repository.Task,
	repo repository.UserPoints,
) *userPointsAppMessageService {
	return &userPointsAppMessageService{
		tr:   tr,
		repo: repo,
	}
}

type userPointsAppMessageService struct {
	tr   repository.Task
	repo repository.UserPoints
}

func (s *userPointsAppMessageService) AddPointsItem(cmd *CmdToAddPointsItem) error {
	task, err := s.tr.Find(cmd.Task)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			return nil
		}

		return err
	}

	date, time := cmd.dateAndTime()
	if date == "" {
		return nil
	}

	up, err := s.repo.Find(cmd.Account, date)
	if err != nil {
		if !repoerr.IsErrorResourceNotExists(err) {
			return err
		}

		up = domain.UserPoints{
			User: cmd.Account,
			Date: date,
		}
	}

	item := up.AddPointsItem(&task, time, cmd.Desc)
	if item == nil {
		return nil
	}

	return s.repo.SavePointsItem(&up, item)
}

package app

import (
	"strconv"

	"github.com/sirupsen/logrus"

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
	task, err := s.tr.Find(cmd.TaskId)
	if err != nil {
		logrus.Errorf("No task found for task id: %s", cmd.TaskId)

		if repoerr.IsErrorResourceNotExists(err) {
			return nil
		}

		return err
	}

	date, time := cmd.dateAndTime()
	if date == "" {
		logrus.Errorf("Failed to get date and time for task: %s, time: %d.", cmd.TaskId, cmd.Time)

		return nil
	}

	up, err := s.repo.Find(cmd.Account, date)
	if err != nil {
		if !repoerr.IsErrorResourceNotExists(err) {
			return err
		}

		up = domain.UserPoints{User: cmd.Account}
	}

	item := up.AddPointsItem(&task, date, &domain.PointsDetail{
		Id:   strconv.FormatInt(cmd.Time, 10),
		Desc: cmd.Desc,
		Time: time,
	})
	if item == nil {
		return nil
	}

	return s.repo.SavePointsItem(&up, item)
}

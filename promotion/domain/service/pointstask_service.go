package service

import (
	"errors"

	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
)

type PointsTaskService interface {
	Save(user types.Account, taskid string) error
	Find(types.Account) (domain.UserPoints, error)
	GetAllUserPoints() ([]domain.UserPoints, error)
}

func NewPointsTaskService(
	pointsRepo repository.Points,
	taskRepo repository.Task,
) (PointsTaskService, error) {
	// get all task
	alltask, err := taskRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// make a task map
	taskmap := make(map[string]domain.Task, len(alltask))
	for i := range alltask {
		taskmap[alltask[i].Id] = alltask[i]
	}

	return &pointsTaskService{
		taskMap:    taskmap,
		pointsRepo: pointsRepo,
		taskRepo:   taskRepo,
	}, nil
}

type pointsTaskService struct {
	taskMap map[string]domain.Task

	pointsRepo repository.Points
	taskRepo   repository.Task
}

func (s *pointsTaskService) Find(u types.Account) (up domain.UserPoints, err error) {
	// get user's points
	p, err := s.pointsRepo.Find(u)
	if err != nil {
		return
	}

	flag, err := s.isUserPointInvalid(p)
	if err != nil {
		return
	}

	if flag {
		err = errors.New("user point invalid, over max allowed")

		return
	}

	return
}

func (s *pointsTaskService) GetAllUserPoints() (ups []domain.UserPoints, err error) {
	// get all points
	ups, err = s.pointsRepo.FindAll()
	if err != nil {
		return
	}

	// is point invalid
	flag, err := s.isUserPointInvalid(ups...)
	if err != nil {
		return
	}

	if flag {
		err = errors.New("user point invalid, over max allowed")

		return
	}

	return
}

func (s *pointsTaskService) isUserPointInvalid(up ...domain.UserPoints) (bool, error) {
	for i := range up {
		flag, err := s.isPointOverMaxAllowed(up[i].Items...)
		if err != nil {
			return true, err
		}

		if flag {
			return true, nil
		}
	}

	return false, nil
}

func (s *pointsTaskService) isPointOverMaxAllowed(item ...domain.Item) (bool, error) {
	for i := range item {
		task, ok := s.taskMap[item[i].TaskId]
		if !ok {
			return true, errors.New("invalid task id")
		}

		if !task.Rule.IsValidPoint(item[i].Points) {
			return true, nil
		}
	}

	return false, nil
}

func (s *pointsTaskService) Save(user types.Account, taskid string) error {
	// find userpoint version
	up, err := s.pointsRepo.Find(user)
	if err != nil {
		return err
	}

	// update
	return s.pointsRepo.Update(user, taskid, up.Version)
}
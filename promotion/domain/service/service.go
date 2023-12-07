package service

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
)

type PointsTaskService interface {
	Find(types.Account) (domain.UserPoints, error)
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

	return s.toUserPoints(p), nil
}

func (s *pointsTaskService) toUserPoints(p repository.Point) domain.UserPoints {
	// generate UserPoints
	var total int
	items := make([]domain.Item, len(p.Dones))
	for i := range p.Dones {
		task := s.taskMap[p.Dones[i].TaskId]
		items[i] = domain.Item{
			Id:       task.Id,
			TaskName: task.Names,
			Descs:    task.Rule.Descs,
			Date:     p.Dones[i].Date,
			Points:   task.Rule.Points,
		}

		total += task.Rule.Points
	}

	return domain.UserPoints{
		Total: total,
		Items: items,
	}
}

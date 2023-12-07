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
) PointsTaskService {
	return &pointsTaskService{
		pointsRepo: pointsRepo,
		taskRepo:   taskRepo,
	}
}

type pointsTaskService struct {
	pointsRepo repository.Points
	taskRepo   repository.Task
}

func (s pointsTaskService) Find(u types.Account) (up domain.UserPoints, err error) {
	// get user's points
	p, err := s.pointsRepo.Find(u)
	if err != nil {
		return
	}

	// get all task
	ts, err := s.taskRepo.FindAll()
	if err != nil {
		return
	}

	return toUserPoints(p, ts), nil
}

func toUserPoints(p repository.Point, allTask []domain.Task) domain.UserPoints {
	// make a task map
	taskmap := make(map[string]domain.Task, len(allTask))
	for i := range allTask {
		taskmap[allTask[i].Id] = allTask[i]
	}

	// generate UserPoints
	var total int
	items := make([]domain.Item, len(p.Dones))
	for i := range p.Dones {
		task := taskmap[p.Dones[i].TaskId]
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

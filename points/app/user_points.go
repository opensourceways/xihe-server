package app

import (
	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

const minValueOfInvlidTime = 24 * 3600 // second

type UserPointsAppService interface {
	Points(account common.Account) (int, error)
	GetPointsDetails(account common.Account) (dto UserPointsDetailsDTO, err error)
	GetTaskCompletionInfo(account common.Account) ([]TasksCompletionInfoDTO, error)
}

func NewUserPointsAppService(
	tr repository.Task,
	repo repository.UserPoints,
) *userPointsAppService {
	return &userPointsAppService{
		tr:   tr,
		repo: repo,
	}
}

type userPointsAppService struct {
	tr   repository.Task
	repo repository.UserPoints
}

func (s *userPointsAppService) Points(account common.Account) (int, error) {
	return 0, nil
	/* TODO retrieve back
	up, err := s.repo.Find(account, utils.Date())
	if err != nil {
		// if not exist
		return 0, nil
	}

	return up.Total, nil
	*/
}

func (s *userPointsAppService) GetPointsDetails(account common.Account) (dto UserPointsDetailsDTO, err error) {
	v, err := s.repo.FindAll(account)
	if err != nil {
		return
	}

	dto.Total = v.Total

	details := make([]PointsDetailDTO, 0, v.DetailNum())

	for i := range v.Items {
		t := v.Items[i].Task

		ds := v.Items[i].Details
		for j := range ds {
			details = append(details, PointsDetailDTO{
				Task:         t,
				PointsDetail: ds[j],
			})
		}
	}

	return
}

func (s *userPointsAppService) GetTaskCompletionInfo(account common.Account) ([]TasksCompletionInfoDTO, error) {
	tasks, err := s.tr.FindAllTasks()
	if err != nil {
		return nil, err
	}

	var isCompleted func(*domain.Task) bool

	up, err := s.repo.Find(account, utils.Date())
	if err != nil {
		// if not exist
		isCompleted = func(*domain.Task) bool {
			return false
		}

		return nil, err
	} else {
		isCompleted = up.IsCompleted
	}

	m := map[string]int{}
	r := []TasksCompletionInfoDTO{}

	for i := range tasks {
		t := &tasks[i]

		j, ok := m[t.Kind]
		if !ok {
			j = len(r)
			r = append(r, TasksCompletionInfoDTO{Kind: t.Kind})
		}

		r[j].add(t, isCompleted(t))
	}

	return r, nil
}

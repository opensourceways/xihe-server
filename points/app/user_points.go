package app

import (
	common "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

const minValueOfInvlidTime = 24 * 3600 // second

type UserPointsAppService interface {
	Points(account common.Account) (int, error)
	PointsDetails(account common.Account) (dto UserPointsDetailsDTO, err error)
	TasksOfDay(account common.Account) ([]TasksCompletionInfoDTO, error)
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
	up, err := s.repo.Find(account, utils.Date())
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			return 0, nil
		}

		return 0, err
	}

	return up.Total, nil
}

func (s *userPointsAppService) PointsDetails(account common.Account) (dto UserPointsDetailsDTO, err error) {
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

func (s *userPointsAppService) TasksOfDay(account common.Account) ([]TasksCompletionInfoDTO, error) {
	tasks, err := s.tr.FindAllTasks()
	if err != nil {
		return nil, err
	}

	var isCompleted func(*domain.Task) bool

	up, err := s.repo.Find(account, utils.Date())
	if err != nil {
		if !repoerr.IsErrorResourceNotExists(err) {
			return nil, err
		}

		isCompleted = func(*domain.Task) bool {
			return false
		}
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
			m[t.Kind] = j

			r = append(r, TasksCompletionInfoDTO{Kind: t.Kind})
		}

		r[j].add(t, isCompleted(t))
	}

	return r, nil
}

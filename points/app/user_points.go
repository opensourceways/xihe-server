package app

import (
	common "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

const minValueOfInvlidTime = 24 * 3600 // second

type UserPointsAppService interface {
	Points(account types.Account) (int, error)
	PointsDetails(account types.Account, lang common.Language) (dto UserPointsDetailsDTO, err error)
	TasksOfDay(account types.Account, lang common.Language) ([]TasksCompletionInfoDTO, error)
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

func (s *userPointsAppService) Points(account types.Account) (int, error) {
	up, err := s.repo.Find(account, utils.Date())
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			return 0, nil
		}

		return 0, err
	}

	return up.Total, nil
}

func (s *userPointsAppService) PointsDetails(account types.Account, lang common.Language) (
	dto UserPointsDetailsDTO, err error,
) {
	tasks, err := s.tr.FindAllTasks()
	if err != nil {
		return
	}

	m := map[string]string{}
	for i := range tasks {
		item := &tasks[i]

		m[item.Id] = item.Name(lang)
	}

	v, err := s.repo.FindAll(account)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	dto.Total = v.Total

	details := make([]PointsDetailDTO, 0, v.DetailsNum())

	for i := range v.Items {
		t := m[v.Items[i].TaskId]

		ds := v.Items[i].Details
		for j := range ds {
			details = append(details, PointsDetailDTO{
				Task:         t,
				PointsDetail: ds[j],
			})
		}
	}

	dto.Details = details

	return
}

func (s *userPointsAppService) TasksOfDay(account types.Account, lang common.Language) (
	[]TasksCompletionInfoDTO, error,
) {
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

		if t.IsPassiveTask() {
			continue
		}

		j, ok := m[t.Kind]
		if !ok {
			j = len(r)
			m[t.Kind] = j

			r = append(r, TasksCompletionInfoDTO{Kind: t.Kind})
		}

		r[j].add(t, isCompleted(t), lang)
	}

	return r, nil
}

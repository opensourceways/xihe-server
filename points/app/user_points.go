package app

import (
	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
	"github.com/opensourceways/xihe-server/points/domain/service"
	"github.com/opensourceways/xihe-server/utils"
)

const minValueOfInvlidTime = 24 * 3600 // second

type UserPointsAppService interface {
	AddPointsItem(cmd *CmdToAddPointsItem) error
	GetTotal(account common.Account) (int, error)
	GetPointsDetails(account common.Account) (dto UserPointsDetailsDTO, err error)
}

type userPointsAppService struct {
	repo repository.UserPoints
}

func (s userPointsAppService) AddPointsItem(cmd *CmdToAddPointsItem) error {
	calculator := service.PointsRuleService().Calculator(cmd.Type)
	if calculator == nil {
		return nil
	}

	date, time := cmd.dateAndTime()
	if date == "" {
		return nil
	}

	up, err := s.repo.Find(cmd.Account, date)
	if err != nil {
		// if not exist
		up = domain.UserPoints{
			User: cmd.Account,
			Date: date,
		}
	}

	detail := domain.PointsDetail{
		Time: time,
		Desc: cmd.Desc,
	}

	item := up.AddPointsItem(cmd.Type, &detail, calculator)
	if item == nil {
		return nil
	}

	return s.repo.SavePointsItem(&up, item)
}

func (s userPointsAppService) GetTotal(account common.Account) (int, error) {
	up, err := s.repo.Find(account, utils.Date())
	if err != nil {
		// if not exist
		return 0, nil
	}

	return up.Total, nil
}

func (s userPointsAppService) GetPointsDetails(account common.Account) (dto UserPointsDetailsDTO, err error) {
	v, err := s.repo.FindAll(account)
	if err != nil {
		return
	}

	dto.Total = v.Total

	details := make([]PointsDetailDTO, 0, v.DetailNum())

	for i := range v.Items {
		t := v.Items[i].Type

		ds := v.Items[i].Details
		for j := range ds {
			details = append(details, PointsDetailDTO{
				Type:         t,
				PointsDetail: ds[j],
			})
		}
	}

	return
}

package app

import (
	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/point/domain"
	"github.com/opensourceways/xihe-server/point/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

const minValueOfInvlidTime = 24 * 3600 // second

type userPointAppService struct {
	repo repository.UserPoint
	rule repository.PointRule
}

func (s userPointAppService) AddPointItem(cmd *CmdToAddPointItem) error {
	rule, err := s.rule.Find(cmd.Type)
	if err != nil {
		return err
	}

	date, time := cmd.dateAndTime()
	if date == "" {
		return nil
	}

	up, err := s.repo.Find(cmd.Account, date)
	if err != nil {
		// if not exist
		up = domain.UserPoint{
			User: cmd.Account,
			Date: date,
		}
	}

	detail := domain.PointDetail{
		Time: time,
		Desc: cmd.Desc,
	}

	item := up.AddPointItem(cmd.Type, &detail, &rule)
	if item == nil {
		return nil
	}

	return s.repo.SavePointItem(&up, item)
}

func (s userPointAppService) GetTotal(account common.Account) (int, error) {
	up, err := s.repo.Find(account, utils.Date())
	if err != nil {
		// if not exist
		return 0, nil
	}

	return up.Total, nil
}

func (s userPointAppService) GetPointDetails(account common.Account) (dto UserPointDetailsDTO, err error) {
	v, err := s.repo.FindAll(account)
	if err != nil {
		return
	}

	dto.Total = v.Total

	details := make([]PointDetailDTO, 0, v.DetailNum())

	for i := range v.Items {
		t := v.Items[i].Type

		ds := v.Items[i].Details
		for j := range ds {
			details = append(details, PointDetailDTO{
				Type:        t,
				PointDetail: ds[j],
			})
		}
	}

	return
}

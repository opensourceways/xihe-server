package app

import (
	"github.com/opensourceways/xihe-server/promotion/domain/service"
)

type Points interface {
	GetPoints(*PointsCmd) (PointsDTO, error)
}

func NewPointsService(
	service service.PointsTaskService,
) Points {
	return &pointService{
		service: service,
	}
}

type pointService struct {
	service service.PointsTaskService
}

func (s *pointService) GetPoints(cmd *PointsCmd) (dto PointsDTO, err error) {
	p, err := s.service.Find(cmd.User)
	if err != nil {
		return
	}

	return toPointsDTO(p, cmd.Lang), nil
}

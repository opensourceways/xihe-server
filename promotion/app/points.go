package app

import (
	"sort"

	"github.com/opensourceways/xihe-server/promotion/domain/service"
)

type Points interface {
	GetPoints(*PointsCmd) (PointsDTO, error)
	GetPointsRank() ([]PointsRankDTO, error)
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

func (s *pointService) GetPointsRank() (dtos []PointsRankDTO, err error) {
	// get all userpoints (not ordered)
	ups, err := s.service.GetAllUserPoints()

	// order uerpoints (desc)
	sort.Slice(ups, func(i, j int) bool {
		return ups[i].Total > ups[j].Total
	})

	// to dto
	dtos = make([]PointsRankDTO, len(ups))
	for i := range ups {
		dtos[i].toDTO(ups[i])
	}

	return
}
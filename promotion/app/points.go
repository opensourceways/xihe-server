package app

import (
	"sort"

	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain/service"
)

type PointsService interface {
	GetPoints(*PointsCmd) (PointsDTO, error)
	GetPointsRank(promotionid string) ([]PointsRankDTO, error)
}

func NewPointsService(
	service service.PointsTaskService,
) PointsService {
	return &pointService{
		service: service,
	}
}

type pointService struct {
	service service.PointsTaskService
}

func (s *pointService) GetPoints(cmd *PointsCmd) (dto PointsDTO, err error) {
	p, err := s.service.Find(cmd.User, cmd.Promotionid)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			return PointsDTO{}, nil
		}

		return
	}

	return toPointsDTO(p, cmd.Lang), nil
}

func (s *pointService) GetPointsRank(promotionid string) (dtos []PointsRankDTO, err error) {
	// get all userpoints (not ordered)
	ups, err := s.service.GetAllUserPoints(promotionid)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			return dtos, nil
		}

		return
	}

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

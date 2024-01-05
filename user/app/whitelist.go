package app

import (
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

type WhiteListService interface {
	// register
	CheckWhiteList(*UserWhiteListCmd) (bool, error)
}

var _ RegService = (*regService)(nil)

func NewWhiteListService(
	whiteRepo repository.WhiteList,
) *whiteListService {
	return &whiteListService{
		whiteRepo: whiteRepo,
	}
}

type whiteListService struct {
	whiteRepo repository.WhiteList
}

func (s *whiteListService) CheckWhiteList(cmd *UserWhiteListCmd) (v bool, err error) {
	u, err := s.whiteRepo.GetWhiteListInfo(cmd.Account, cmd.Type.WhiteListType())
	if err != nil {
		return
	}
	if !u.Enabled {
		v = false
		return
	}
	v = true
	return
}

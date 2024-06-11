package app

import (
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

type WhiteListService interface {
	// register
	CheckWhiteList(*UserWhiteListCmd) (*WhitelistDTO, error)
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

func (s *whiteListService) CheckWhiteList(cmd *UserWhiteListCmd) (*WhitelistDTO, error) {
	u, err := s.whiteRepo.GetWhiteListInfo(cmd.Account, cmd.Type.WhiteListType())
	if err != nil {
		return nil, err
	}

	w := &WhitelistDTO{}
	w.toWhitelistDTO(&u)

	return w, nil
}

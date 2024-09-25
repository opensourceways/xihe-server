package app

import (
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

type WhiteListService interface {
	// register
	CheckWhiteList(*UserWhiteListCmd) (*WhitelistDTO, error)
	CheckCloudWhitelist(domain.Account) (bool, bool, error)
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
	whitelistType := cmd.Type.WhiteListType()
	whitelistDTO := &WhitelistDTO{}

	if whitelistType == domain.WhitelistTypeCloud || whitelistType == domain.WhitelistTypeMultiCloud {
		useNPU, useMultiNPU, err := s.CheckCloudWhitelist(cmd.Account)
		if err != nil {
			return nil, err
		}

		whitelistDTO.Allowed = useNPU
		if whitelistType == domain.WhitelistTypeMultiCloud {
			whitelistDTO.Allowed = useMultiNPU

		}

		return whitelistDTO, nil
	}

	u, err := s.whiteRepo.GetWhiteListInfo(cmd.Account, cmd.Type.WhiteListType())
	if err != nil {
		return nil, err
	}

	whitelistDTO.toWhitelistDTO(&u)

	return whitelistDTO, nil
}

// CheckCloudWhitelist reports whether an account can use npu with single card or multiple cards
func (s *whiteListService) CheckCloudWhitelist(account domain.Account) (useNPU bool, useMultiNPU bool, err error) {
	var whitelistItems []domain.WhiteListInfo

	whitelistItems, err = s.whiteRepo.GetWhiteListInfoItems(
		account, []string{domain.WhitelistTypeCloud, domain.WhitelistTypeMultiCloud})
	if err != nil {
		return
	}

	for _, item := range whitelistItems {
		switch item.Type.WhiteListType() {
		case domain.WhitelistTypeCloud:
			useNPU = item.Enable()
		case domain.WhitelistTypeMultiCloud:
			useMultiNPU = item.Enable()
		}
	}

	// user who can use multiple cards can use single card
	if useMultiNPU {
		useNPU = useMultiNPU
	}

	return
}

package app

import (
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

type WhiteListService interface {
	// register
	CheckWhiteList(*UserWhiteListCmd) (*WhitelistDTO, error)
	CheckCloudWhitelist(domain.Account) (bool, bool, error)
	List(domain.Account) ([]string, error)
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
	whitelistDTO := &WhitelistDTO{}

	if cmd.Type.WhiteListType() == domain.WhitelistTypeCloud ||
		cmd.Type.WhiteListType() == domain.WhitelistTypeMultiCloud {
		useNPU, useMultiNPU, err := s.CheckCloudWhitelist(cmd.Account)
		if err != nil {
			return nil, err
		}

		whitelistDTO.Allowed = useNPU
		if cmd.Type.WhiteListType() == domain.WhitelistTypeMultiCloud {
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

	whitelistItems, err = s.whiteRepo.FindByAccountAndWhitelistType(
		account, []string{domain.WhitelistTypeCloud, domain.WhitelistTypeMultiCloud})
	if err != nil {
		return
	}

	pairs := s.getWhitelistTypePairs(whitelistItems)
	useNPU, useMultiNPU = pairs[domain.WhitelistTypeCloud], pairs[domain.WhitelistTypeMultiCloud]

	// user who uses multiple cards can use single card
	if useMultiNPU {
		useNPU = useMultiNPU
	}

	return
}

// List return a list of whiyelist by account
func (s *whiteListService) List(account domain.Account) ([]string, error) {
	items, err := s.whiteRepo.FindByAccountAndWhitelistType(account, nil)
	if err != nil {
		return nil, err
	}

	pairs := s.getWhitelistTypePairs(items)
	if pairs[domain.WhitelistTypeMultiCloud] {
		pairs[domain.WhitelistTypeCloud] = true
	}

	list := make([]string, 0, len(items))
	for k, v := range pairs {
		if v {
			list = append(list, k)
		}
	}

	return list, nil
}

func (s *whiteListService) getWhitelistTypePairs(items []domain.WhiteListInfo) map[string]bool {
	mp := make(map[string]bool)

	for _, v := range items {
		mp[v.Type.WhiteListType()] = v.Enable()
	}

	return mp
}

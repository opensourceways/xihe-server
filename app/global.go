package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
)

type GlobalResourceListCmd struct {
	repository.GlobalResourceListOption

	SortType domain.SortType
}

func (cmd *GlobalResourceListCmd) ToResourceListOption() repository.GlobalResourceListOption {
	// only allow to list public resources.
	type1, _ := domain.NewRepoType(domain.RepoTypePublic)
	type2, _ := domain.NewRepoType(domain.RepoTypeOnline)
	cmd.RepoType = append(cmd.RepoType, type1)
	cmd.RepoType = append(cmd.RepoType, type2)
	return cmd.GlobalResourceListOption
}

// Model
type GlobalModelsDTO struct {
	Total  int              `json:"total"`
	Models []GlobalModelDTO `json:"projects"`
}

type GlobalModelDTO struct {
	ModelSummaryDTO
	AvatarId string `json:"avatar_id"`
}

func (s modelService) ListGlobal(cmd *GlobalResourceListCmd) (
	dto GlobalModelsDTO, err error,
) {
	option := cmd.ToResourceListOption()

	var v repository.UserModelsInfo

	if cmd.SortType == nil {
		v, err = s.repo.ListGlobalAndSortByUpdateTime(&option)
	} else {
		switch cmd.SortType.SortType() {
		case domain.SortTypeUpdateTime:
			v, err = s.repo.ListGlobalAndSortByUpdateTime(&option)

		case domain.SortTypeFirstLetter:
			v, err = s.repo.ListGlobalAndSortByFirstLetter(&option)

		case domain.SortTypeDownloadCount:
			v, err = s.repo.ListGlobalAndSortByDownloadCount(&option)
		}
	}

	items := v.Models

	if err != nil || len(items) == 0 {
		return
	}

	// find avatars
	users := make([]userdomain.Account, len(items))
	for i := range items {
		users[i] = items[i].Owner
	}

	avatars, err := s.rs.FindUserAvater(users)
	if err != nil {
		return
	}

	// gen result
	dtos := make([]GlobalModelDTO, len(items))
	for i := range items {
		s.toModelSummaryDTO(&items[i], &dtos[i].ModelSummaryDTO)
		dtos[i].AvatarId = avatars[i]
	}

	dto.Total = v.Total
	dto.Models = dtos

	return
}

// Dataset
type GlobalDatasetsDTO struct {
	Total    int                `json:"total"`
	Datasets []GlobalDatasetDTO `json:"projects"`
}

type GlobalDatasetDTO struct {
	DatasetSummaryDTO
	AvatarId string `json:"avatar_id"`
}

func (s datasetService) ListGlobal(cmd *GlobalResourceListCmd) (
	dto GlobalDatasetsDTO, err error,
) {
	option := cmd.ToResourceListOption()

	var v repository.UserDatasetsInfo

	if cmd.SortType == nil {
		v, err = s.repo.ListGlobalAndSortByUpdateTime(&option)
	} else {
		switch cmd.SortType.SortType() {
		case domain.SortTypeUpdateTime:
			v, err = s.repo.ListGlobalAndSortByUpdateTime(&option)

		case domain.SortTypeFirstLetter:
			v, err = s.repo.ListGlobalAndSortByFirstLetter(&option)

		case domain.SortTypeDownloadCount:
			v, err = s.repo.ListGlobalAndSortByDownloadCount(&option)
		}
	}

	items := v.Datasets

	if err != nil || len(items) == 0 {
		return
	}

	// find avatars
	users := make([]userdomain.Account, len(items))
	for i := range items {
		users[i] = items[i].Owner
	}

	avatars, err := s.rs.FindUserAvater(users)
	if err != nil {
		return
	}

	// gen result
	dtos := make([]GlobalDatasetDTO, len(items))
	for i := range items {
		s.toDatasetSummaryDTO(&items[i], &dtos[i].DatasetSummaryDTO)
		dtos[i].AvatarId = avatars[i]
	}

	dto.Total = v.Total
	dto.Datasets = dtos

	return
}

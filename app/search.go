package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
)

type SearchDTO struct {
	User    UserSearchDTO     `json:"user"`
	Model   ResourceSearchDTO `json:"model"`
	Project ResourceSearchDTO `json:"project"`
	Dataset ResourceSearchDTO `json:"dataset"`
}

type UserSearchDTO struct {
	Top   []string `json:"top"`
	Total int      `json:"total"`
}

type ResourceSearchDTO struct {
	Top   []ResourceSummaryDTO `json:"top"`
	Total int                  `json:"total"`
}

type ResourceSummaryDTO struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func newResourceSearchOption(name string) repository.ResourceSearchOption {
	option := repository.ResourceSearchOption{
		Name:   name,
		TopNum: 3,
	}
	type1, _ := domain.NewRepoType(domain.RepoTypePublic)
	type2, _ := domain.NewRepoType(domain.RepoTypeOnline)
	option.RepoType = append(option.RepoType, type1)
	option.RepoType = append(option.RepoType, type2)

	return option
}

type SearchService interface {
	Search(name string) (dto SearchDTO)
}

func NewSearchService(
	user userrepo.User,
	model repository.Model,
	projectPg spacerepo.ProjectPg,
	dataset repository.Dataset,
) SearchService {
	return searchService{
		user:      user,
		model:     model,
		projectPg: projectPg,
		dataset:   dataset,
	}
}

type searchService struct {
	user      userrepo.User
	model     repository.Model
	projectPg spacerepo.ProjectPg
	dataset   repository.Dataset
}

func (s searchService) Search(name string) (dto SearchDTO) {
	option := newResourceSearchOption(name)

	if u, err := s.searchUser(name); err == nil {
		dto.User = u
	}

	v, err := s.search(&option, s.projectPg.Search)
	if err == nil {
		dto.Project = v
	}

	v, err = s.search(&option, s.model.Search)
	if err == nil {
		dto.Model = v
	}

	v, err = s.search(&option, s.dataset.Search)
	if err == nil {
		dto.Dataset = v
	}

	return
}

func (s searchService) search(
	option *repository.ResourceSearchOption,
	f func(*repository.ResourceSearchOption) (repository.ResourceSearchResult, error),
) (
	dto ResourceSearchDTO, err error,
) {
	v, err := f(option)
	if err != nil || v.Total == 0 {
		return
	}

	items := make([]ResourceSummaryDTO, len(v.Top))
	for i := range v.Top {
		item := &v.Top[i]

		items[i].Owner = item.Owner.Account()
		items[i].Name = item.Name.ResourceName()
	}

	dto.Top = items
	dto.Total = v.Total

	return
}

func (s searchService) searchUser(name string) (
	dto UserSearchDTO, err error,
) {
	v, err := s.user.Search(&userrepo.UserSearchOption{
		Name:   name,
		TopNum: 10,
	})
	if err != nil || v.Total == 0 {
		return
	}

	items := make([]string, len(v.Top))
	for i := range v.Top {
		items[i] = v.Top[i].Account()
	}

	dto.Top = items
	dto.Total = v.Total

	return
}

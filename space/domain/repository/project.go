package repository

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
)

type ProjectPropertyUpdateInfo struct {
	repository.ResourceToUpdate

	Property spacedomain.ProjectModifiableProperty
}

type UserProjectsInfo struct {
	Projects []spacedomain.ProjectSummary
	Total    int
}

type ProjectSummary struct {
	domain.ResourceSummary
	Tags []string
}

type Project interface {
	Save(*spacedomain.Project) (spacedomain.Project, error)
	Delete(*domain.ResourceIndex) error
	Get(domain.Account, string) (spacedomain.Project, error)
	GetByName(domain.Account, domain.ResourceName) (spacedomain.Project, error)
	GetByRepoId(domain.Identity) (spacedomain.Project, error)
	GetSummary(domain.Account, string) (ProjectSummary, error)
	GetSummaryByName(domain.Account, domain.ResourceName) (domain.ResourceSummary, error)

	FindUserProjects([]repository.UserResourceListOption) ([]spacedomain.ProjectSummary, error)

	ListAndSortByUpdateTime(domain.Account, *repository.ResourceListOption) (UserProjectsInfo, error)
	ListAndSortByFirstLetter(domain.Account, *repository.ResourceListOption) (UserProjectsInfo, error)
	ListAndSortByDownloadCount(domain.Account, *repository.ResourceListOption) (UserProjectsInfo, error)

	ListGlobalAndSortByUpdateTime(*repository.GlobalResourceListOption) (UserProjectsInfo, error)
	ListGlobalAndSortByFirstLetter(*repository.GlobalResourceListOption) (UserProjectsInfo, error)
	ListGlobalAndSortByDownloadCount(*repository.GlobalResourceListOption) (UserProjectsInfo, error)

	Search(*repository.ResourceSearchOption) (repository.ResourceSearchResult, error)

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	AddRelatedModel(*repository.RelatedResourceInfo) error
	RemoveRelatedModel(*repository.RelatedResourceInfo) error

	AddRelatedDataset(*repository.RelatedResourceInfo) error
	RemoveRelatedDataset(*repository.RelatedResourceInfo) error

	UpdateProperty(*ProjectPropertyUpdateInfo) error

	IncreaseFork(*domain.ResourceIndex) error
	IncreaseDownload(*domain.ResourceIndex) error
}
type ProjectPg interface {
	Save(*spacedomain.Project) (spacedomain.Project, error)
	GetByName(domain.Account, domain.ResourceName) (spacedomain.Project, error)
	Get(domain.Account, string) (spacedomain.Project, error)

	AddRelatedDataset(*repository.RelatedResourceInfo) error
	AddRelatedModel(*repository.RelatedResourceInfo) error

	ListAndSortByUpdateTime(domain.Account, *repository.ResourceListOption) (UserProjectsInfo, error)
	ListAndSortByFirstLetter(domain.Account, *repository.ResourceListOption) (UserProjectsInfo, error)
	ListAndSortByDownloadCount(domain.Account, *repository.ResourceListOption) (UserProjectsInfo, error)
}

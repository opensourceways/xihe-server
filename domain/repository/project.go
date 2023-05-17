package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserResourceListOption struct {
	Owner domain.Account
	Ids   []string
}

type ResourceSearchOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name     string
	TopNum   int
	RepoType domain.RepoType
}

type ResourceListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name     string
	RepoType domain.RepoType

	PageNum      int
	CountPerPage int
}

type GlobalResourceListOption struct {
	Level    domain.ResourceLevel
	Tags     []string
	TagKinds []string

	ResourceListOption
}

type RelatedResourceInfo struct {
	ResourceToUpdate

	RelatedResource domain.ResourceIndex
}

type ProjectPropertyUpdateInfo struct {
	ResourceToUpdate

	Property domain.ProjectModifiableProperty
}

type ResourceToUpdate struct {
	Owner     domain.Account
	Id        string
	Version   int
	UpdatedAt int64
}

type UserProjectsInfo struct {
	Projects []domain.ProjectSummary
	Total    int
}

type ResourceSearchResult struct {
	Top []domain.ResourceSummary

	Total int
}

type ProjectSummary struct {
	domain.ResourceSummary
	Tags []string
}

type Project interface {
	Save(*domain.Project) (domain.Project, error)
	Delete(*domain.ResourceIndex) error
	Get(domain.Account, string) (domain.Project, error)
	GetByName(domain.Account, domain.ResourceName) (domain.Project, error)
	GetSummary(domain.Account, string) (ProjectSummary, error)
	GetSummaryByName(domain.Account, domain.ResourceName) (domain.ResourceSummary, error)

	FindUserProjects([]UserResourceListOption) ([]domain.ProjectSummary, error)

	ListAndSortByUpdateTime(domain.Account, *ResourceListOption) (UserProjectsInfo, error)
	ListAndSortByFirstLetter(domain.Account, *ResourceListOption) (UserProjectsInfo, error)
	ListAndSortByDownloadCount(domain.Account, *ResourceListOption) (UserProjectsInfo, error)

	ListGlobalAndSortByUpdateTime(*GlobalResourceListOption) (UserProjectsInfo, error)
	ListGlobalAndSortByFirstLetter(*GlobalResourceListOption) (UserProjectsInfo, error)
	ListGlobalAndSortByDownloadCount(*GlobalResourceListOption) (UserProjectsInfo, error)

	Search(*ResourceSearchOption) (ResourceSearchResult, error)

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	AddRelatedModel(*RelatedResourceInfo) error
	RemoveRelatedModel(*RelatedResourceInfo) error

	AddRelatedDataset(*RelatedResourceInfo) error
	RemoveRelatedDataset(*RelatedResourceInfo) error

	UpdateProperty(*ProjectPropertyUpdateInfo) error

	IncreaseFork(*domain.ResourceIndex) error
	IncreaseDownload(*domain.ResourceIndex) error
}

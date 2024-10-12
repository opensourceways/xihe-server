package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserResourceListOption struct {
	Owner domain.Account
	Ids   []string
}

type ResourceListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name     string
	RepoType []domain.RepoType

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

type ResourceToUpdate struct {
	Owner     domain.Account
	Id        string
	Version   int
	UpdatedAt int64
}

type ResourceSearchOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name     string
	TopNum   int
	RepoType []domain.RepoType
}

type ResourceSearchResult struct {
	Top []domain.ResourceSummary

	Total int
}

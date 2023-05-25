package domain

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
)

type ResourceObjects struct {
	Type    ResourceType
	Objects []ResourceIndex
}

// ResourceObject
type ResourceObject struct {
	Type ResourceType

	ResourceIndex
}

func (r *ResourceObject) String() string {
	return fmt.Sprintf(
		"%s_%s_%s",
		r.Owner.Account(),
		r.Type.ResourceType(),
		r.Id,
	)
}

type ResourceIndex struct {
	Owner Account
	Id    string
}

type RelatedResources []ResourceIndex

func (r RelatedResources) Has(index *ResourceIndex) bool {
	v := sets.NewString()

	for i := range ([]ResourceIndex)(r) {
		v.Insert(
			r[i].Owner.Account() + r[i].Id,
		)
	}

	return v.Has(index.Owner.Account() + index.Id)
}

func (r RelatedResources) Count() int {
	return len(r)
}

type ReverselyRelatedResourceInfo struct {
	Promoter *ResourceIndex
	Resource *ResourceIndex
}

type ResourceSummary struct {
	Owner    Account
	Name     ResourceName
	Id       string
	RepoId   string
	RepoType RepoType
}

func (s *ResourceSummary) IsPrivate() bool {
	return s.RepoType.RepoType() == RepoTypePrivate
}

func (s *ResourceSummary) IsOnline() bool {
	return s.RepoType.RepoType() == RepoTypeOnline
}

func (s *ResourceSummary) ResourceIndex() ResourceIndex {
	return ResourceIndex{
		Owner: s.Owner,
		Id:    s.Id,
	}
}

package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ResourceListDO struct {
	Name         string
	RepoType     []string
	PageNum      int
	CountPerPage int
}

func toResourceListDO(r *repository.ResourceListOption) ResourceListDO {
	do := ResourceListDO{
		Name:         r.Name,
		PageNum:      r.PageNum,
		CountPerPage: r.CountPerPage,
	}

	if r.RepoType != nil {
		for i := range r.RepoType {
			do.RepoType = append(do.RepoType, r.RepoType[i].RepoType())
		}
	}

	return do
}

type ResourceObjectDO struct {
	Owner string
	Type  string
	Id    string
}

func (do *ResourceObjectDO) toResourceObject(r *domain.ResourceObject) (err error) {
	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Type, err = domain.NewResourceType(do.Type); err != nil {
		return
	}

	r.Id = do.Id

	return
}

func toResourceObjectDO(r *domain.ResourceObject) ResourceObjectDO {
	return ResourceObjectDO{
		Owner: r.Owner.Account(),
		Type:  r.Type.ResourceType(),
		Id:    r.Id,
	}
}

type ResourceIndexDO struct {
	Owner string
	Id    string
}

func (do *ResourceIndexDO) toResourceIndex(r *domain.ResourceIndex) (err error) {
	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	r.Id = do.Id

	return
}

func toResourceIndexDO(r *domain.ResourceIndex) ResourceIndexDO {
	return ResourceIndexDO{
		Owner: r.Owner.Account(),
		Id:    r.Id,
	}
}

func convertToResourceIndex(v []ResourceIndexDO) (r []domain.ResourceIndex, err error) {
	if len(v) == 0 {
		return
	}

	r = make([]domain.ResourceIndex, len(v))

	for i := range v {
		if err = v[i].toResourceIndex(&r[i]); err != nil {
			return
		}
	}

	return
}

func toReverselyRelatedResourceInfoDO(
	info *domain.ReverselyRelatedResourceInfo,
) ReverselyRelatedResourceInfoDO {
	return ReverselyRelatedResourceInfoDO{
		Promoter: toResourceIndexDO(info.Promoter),
		Resource: toResourceIndexDO(info.Resource),
	}
}

type ReverselyRelatedResourceInfoDO struct {
	Promoter ResourceIndexDO
	Resource ResourceIndexDO
}

type ResourceSummaryDO struct {
	Owner    string
	Name     string
	Id       string
	RepoId   string
	RepoType string
}

func (do *ResourceSummaryDO) toProject() (s domain.ResourceSummary, err error) {
	if s.Name, err = domain.NewResourceName(do.Name); err != nil {
		return
	}

	return s, do.convert(&s)
}

func (do *ResourceSummaryDO) toModel() (s domain.ResourceSummary, err error) {
	if s.Name, err = domain.NewResourceName(do.Name); err != nil {
		return
	}

	return s, do.convert(&s)
}

func (do *ResourceSummaryDO) toDataset() (s domain.ResourceSummary, err error) {
	if s.Name, err = domain.NewResourceName(do.Name); err != nil {
		return
	}

	return s, do.convert(&s)
}

func (do *ResourceSummaryDO) convert(s *domain.ResourceSummary) (err error) {
	if s.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if s.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
		return
	}

	s.Id = do.Id
	s.RepoId = do.RepoId

	return
}

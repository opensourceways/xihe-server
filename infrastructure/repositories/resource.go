package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func ToRelatedResourceDO(info *repository.RelatedResourceInfo) RelatedResourceDO {
	return RelatedResourceDO{
		ResourceToUpdateDO: ToResourceToUpdateDO(&info.ResourceToUpdate),
		ResourceOwner:      info.RelatedResource.Owner.Account(),
		ResourceId:         info.RelatedResource.Id,
	}
}

type RelatedResourceDO struct {
	ResourceToUpdateDO

	ResourceOwner string
	ResourceId    string
}

type ResourceToUpdateDO struct {
	Id        string
	Owner     string
	Version   int
	UpdatedAt int64
}

func ToResourceToUpdateDO(info *repository.ResourceToUpdate) ResourceToUpdateDO {
	return ResourceToUpdateDO{
		Id:        info.Id,
		Owner:     info.Owner.Account(),
		Version:   info.Version,
		UpdatedAt: info.UpdatedAt,
	}
}

type ResourceListDO struct {
	Name         string
	RepoType     []string
	PageNum      int
	CountPerPage int
}

func ToResourceListDO(r *repository.ResourceListOption) ResourceListDO {
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

func ToResourceIndexDO(r *domain.ResourceIndex) ResourceIndexDO {
	return ResourceIndexDO{
		Owner: r.Owner.Account(),
		Id:    r.Id,
	}
}

func ConvertToResourceIndex(v []ResourceIndexDO) (r []domain.ResourceIndex, err error) {
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
		Promoter: ToResourceIndexDO(info.Promoter),
		Resource: ToResourceIndexDO(info.Resource),
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

func (do *ResourceSummaryDO) ToProject() (s domain.ResourceSummary, err error) {
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

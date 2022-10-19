package repositories

import (
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type DatasetMapper interface {
	Insert(DatasetDO) (string, error)
	Get(string, string) (DatasetDO, error)
	GetByName(string, string) (DatasetDO, error)
	GetSummaryByName(string, string) (ResourceSummaryDO, error)

	ListUsersDatasets(map[string][]string) ([]DatasetSummaryDO, error)
	ListSummary(map[string][]string) ([]ResourceSummaryDO, error)

	List(string, *ResourceListDO) ([]DatasetSummaryDO, int, error)
	ListAndSortByUpdateTime(string, *ResourceListDO) ([]DatasetSummaryDO, int, error)
	ListAndSortByFirstLetter(string, *ResourceListDO) ([]DatasetSummaryDO, int, error)
	ListAndSortByDownloadCount(string, *ResourceListDO) ([]DatasetSummaryDO, int, error)

	AddLike(ResourceIndexDO) error
	RemoveLike(ResourceIndexDO) error

	AddRelatedProject(*ReverselyRelatedResourceInfoDO) error
	RemoveRelatedProject(*ReverselyRelatedResourceInfoDO) error

	AddRelatedModel(*ReverselyRelatedResourceInfoDO) error
	RemoveRelatedModel(*ReverselyRelatedResourceInfoDO) error

	UpdateProperty(*DatasetPropertyDO) error
}

func NewDatasetRepository(mapper DatasetMapper) repository.Dataset {
	return dataset{mapper}
}

type dataset struct {
	mapper DatasetMapper
}

func (impl dataset) Save(d *domain.Dataset) (r domain.Dataset, err error) {
	if d.Id != "" {
		err = errors.New("must be a new dataset")

		return
	}

	v, err := impl.mapper.Insert(impl.toDatasetDO(d))
	if err != nil {
		err = convertError(err)
	} else {
		r = *d
		r.Id = v
	}

	return
}

func (impl dataset) Get(owner domain.Account, identity string) (r domain.Dataset, err error) {
	v, err := impl.mapper.Get(owner.Account(), identity)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toDataset(&r)
	}

	return
}

func (impl dataset) GetByName(owner domain.Account, name domain.DatasetName) (
	r domain.Dataset, err error,
) {
	v, err := impl.mapper.GetByName(owner.Account(), name.DatasetName())
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toDataset(&r)
	}

	return
}

func (impl dataset) FindUserDatasets(opts []repository.UserResourceListOption) (
	[]domain.DatasetSummary, error,
) {
	do := make(map[string][]string)
	for i := range opts {
		do[opts[i].Owner.Account()] = opts[i].Ids
	}

	v, err := impl.mapper.ListUsersDatasets(do)
	if err != nil {
		return nil, convertError(err)
	}

	r := make([]domain.DatasetSummary, len(v))
	for i := range v {
		if err = v[i].toDatasetSummary(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl dataset) ListSummary(opts []repository.DatasetSummaryListOption) (
	[]domain.ResourceSummary, error,
) {
	m := map[string][]string{}
	all := sets.NewString()

	for i := range opts {
		owner := opts[i].Owner.Account()
		name := opts[i].Name.DatasetName()

		all.Insert(owner + name)

		if v, ok := m[owner]; ok {
			m[owner] = append(v, name)
		} else {
			m[owner] = []string{name}
		}
	}

	v, err := impl.mapper.ListSummary(m)
	if err != nil {
		return nil, convertError(err)
	}

	r := make([]domain.ResourceSummary, 0, len(opts))

	for i := range v {
		if all.Has(v[i].Owner + v[i].Name) {
			s, err := v[i].toDataset()
			if err != nil {
				return nil, err
			}

			r = append(r, s)
		}
	}

	return r, nil
}

func (impl dataset) GetSummaryByName(owner domain.Account, name domain.ResourceName) (
	domain.ResourceSummary, error,
) {
	v, err := impl.mapper.GetSummaryByName(owner.Account(), name.ResourceName())
	if err != nil {
		return domain.ResourceSummary{}, convertError(err)
	}

	return v.toDataset()
}

func (impl dataset) toDatasetDO(d *domain.Dataset) DatasetDO {
	return DatasetDO{
		Id:        d.Id,
		Owner:     d.Owner.Account(),
		Name:      d.Name.DatasetName(),
		FL:        d.Name.FirstLetterOfName(),
		Desc:      d.Desc.ResourceDesc(),
		RepoType:  d.RepoType.RepoType(),
		Protocol:  d.Protocol.ProtocolName(),
		Tags:      d.Tags,
		RepoId:    d.RepoId,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Version:   d.Version,
	}
}

type DatasetDO struct {
	Id            string
	Owner         string
	Name          string
	FL            byte
	Desc          string
	Protocol      string
	RepoType      string
	RepoId        string
	Tags          []string
	CreatedAt     int64
	UpdatedAt     int64
	Version       int
	LikeCount     int
	DownloadCount int

	RelatedModels   []ResourceIndexDO
	RelatedProjects []ResourceIndexDO
}

func (do *DatasetDO) toDataset(r *domain.Dataset) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewDatasetName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewResourceDesc(do.Desc); err != nil {
		return
	}

	if r.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
		return
	}

	if r.Protocol, err = domain.NewProtocolName(do.Protocol); err != nil {
		return
	}

	if r.RelatedModels, err = convertToResourceIndex(do.RelatedModels); err != nil {
		return
	}

	if r.RelatedProjects, err = convertToResourceIndex(do.RelatedProjects); err != nil {
		return
	}

	r.RepoId = do.RepoId
	r.Tags = do.Tags
	r.Version = do.Version
	r.CreatedAt = do.CreatedAt
	r.UpdatedAt = do.UpdatedAt
	r.LikeCount = do.LikeCount
	r.DownloadCount = do.DownloadCount

	return
}

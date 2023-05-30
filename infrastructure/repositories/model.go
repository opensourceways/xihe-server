package repositories

import (
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ModelMapper interface {
	Insert(ModelDO) (string, error)
	Delete(*ResourceIndexDO) error
	Get(string, string) (ModelDO, error)
	GetByName(string, string) (ModelDO, error)
	GetSummaryByName(string, string) (ResourceSummaryDO, error)

	ListUsersModels(map[string][]string) ([]ModelSummaryDO, error)
	ListSummary(map[string][]string) ([]ResourceSummaryDO, error)

	ListAndSortByUpdateTime(string, *ResourceListDO) ([]ModelSummaryDO, int, error)
	ListAndSortByFirstLetter(string, *ResourceListDO) ([]ModelSummaryDO, int, error)
	ListAndSortByDownloadCount(string, *ResourceListDO) ([]ModelSummaryDO, int, error)

	ListGlobalAndSortByUpdateTime(*GlobalResourceListDO) ([]ModelSummaryDO, int, error)
	ListGlobalAndSortByFirstLetter(*GlobalResourceListDO) ([]ModelSummaryDO, int, error)
	ListGlobalAndSortByDownloadCount(*GlobalResourceListDO) ([]ModelSummaryDO, int, error)

	Search(do *GlobalResourceListDO, topNum int) ([]ResourceSummaryDO, int, error)

	IncreaseDownload(ResourceIndexDO) error

	AddLike(ResourceIndexDO) error
	RemoveLike(ResourceIndexDO) error

	AddRelatedDataset(*RelatedResourceDO) error
	RemoveRelatedDataset(*RelatedResourceDO) error

	AddRelatedProject(*ReverselyRelatedResourceInfoDO) error
	RemoveRelatedProject(*ReverselyRelatedResourceInfoDO) error

	UpdateProperty(*ModelPropertyDO) error
}

func NewModelRepository(mapper ModelMapper) repository.Model {
	return model{mapper}
}

type model struct {
	mapper ModelMapper
}

func (impl model) Save(m *domain.Model) (r domain.Model, err error) {
	if m.Id != "" {
		err = errors.New("must be a new model")

		return
	}

	v, err := impl.mapper.Insert(impl.toModelDO(m))
	if err != nil {
		err = convertError(err)
	} else {
		r = *m
		r.Id = v
	}

	return
}

func (impl model) Delete(index *domain.ResourceIndex) (err error) {
	do := toResourceIndexDO(index)

	if err = impl.mapper.Delete(&do); err != nil {
		err = convertError(err)
	}

	return
}

func (impl model) Get(owner domain.Account, identity string) (r domain.Model, err error) {
	v, err := impl.mapper.Get(owner.Account(), identity)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toModel(&r)
	}

	return
}

func (impl model) GetByName(owner domain.Account, name domain.ResourceName) (
	r domain.Model, err error,
) {
	v, err := impl.mapper.GetByName(owner.Account(), name.ResourceName())
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toModel(&r)
	}

	return
}

func (impl model) FindUserModels(opts []repository.UserResourceListOption) (
	[]domain.ModelSummary, error,
) {
	do := make(map[string][]string)
	for i := range opts {
		do[opts[i].Owner.Account()] = opts[i].Ids
	}

	v, err := impl.mapper.ListUsersModels(do)
	if err != nil {
		return nil, convertError(err)
	}

	r := make([]domain.ModelSummary, len(v))
	for i := range v {
		if err = v[i].toModelSummary(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl model) ListSummary(opts []repository.ResourceSummaryListOption) (
	[]domain.ResourceSummary, error,
) {
	m := map[string][]string{}
	all := sets.NewString()

	for i := range opts {
		owner := opts[i].Owner.Account()
		name := opts[i].Name.ResourceName()

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
			s, err := v[i].toModel()
			if err != nil {
				return nil, err
			}

			r = append(r, s)
		}
	}

	return r, nil
}

func (impl model) GetSummaryByName(owner domain.Account, name domain.ResourceName) (
	domain.ResourceSummary, error,
) {
	v, err := impl.mapper.GetSummaryByName(owner.Account(), name.ResourceName())
	if err != nil {
		return domain.ResourceSummary{}, convertError(err)
	}

	return v.toModel()
}

func (impl model) toModelDO(m *domain.Model) ModelDO {
	do := ModelDO{
		Id:        m.Id,
		Owner:     m.Owner.Account(),
		Name:      m.Name.ResourceName(),
		FL:        m.Name.FirstLetterOfName(),
		RepoType:  m.RepoType.RepoType(),
		Protocol:  m.Protocol.ProtocolName(),
		Tags:      m.Tags,
		TagKinds:  m.TagKinds,
		RepoId:    m.RepoId,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Version:   m.Version,
	}

	if m.Desc != nil {
		do.Desc = m.Desc.ResourceDesc()
	}

	if m.Title != nil {
		do.Title = m.Title.ResourceTitle()
	}

	return do
}

type ModelDO struct {
	Id            string
	Owner         string
	Name          string
	FL            byte
	Desc          string
	Title         string
	Protocol      string
	RepoType      string
	RepoId        string
	Tags          []string
	TagKinds      []string
	CreatedAt     int64
	UpdatedAt     int64
	Version       int
	LikeCount     int
	DownloadCount int

	RelatedDatasets []ResourceIndexDO
	RelatedProjects []ResourceIndexDO
}

func (do *ModelDO) toModel(r *domain.Model) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewResourceName(do.Name); err != nil {
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

	if r.RelatedDatasets, err = convertToResourceIndex(do.RelatedDatasets); err != nil {
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

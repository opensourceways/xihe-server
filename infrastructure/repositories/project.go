package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ProjectMapper interface {
	Insert(ProjectDO) (string, error)
	Delete(*ResourceIndexDO) error
	Get(string, string) (ProjectDO, error)
	GetByName(string, string) (ProjectDO, error)
	GetSummary(string, string) (ProjectResourceSummaryDO, error)
	GetSummaryByName(string, string) (ResourceSummaryDO, error)

	ListUsersProjects(map[string][]string) ([]ProjectSummaryDO, error)

	ListAndSortByUpdateTime(string, *ResourceListDO) ([]ProjectSummaryDO, int, error)
	ListAndSortByFirstLetter(string, *ResourceListDO) ([]ProjectSummaryDO, int, error)
	ListAndSortByDownloadCount(string, *ResourceListDO) ([]ProjectSummaryDO, int, error)

	ListGlobalAndSortByUpdateTime(*GlobalResourceListDO) ([]ProjectSummaryDO, int, error)
	ListGlobalAndSortByFirstLetter(*GlobalResourceListDO) ([]ProjectSummaryDO, int, error)
	ListGlobalAndSortByDownloadCount(*GlobalResourceListDO) ([]ProjectSummaryDO, int, error)

	Search(do *GlobalResourceListDO, topNum int) ([]ResourceSummaryDO, int, error)

	IncreaseFork(ResourceIndexDO) error
	IncreaseDownload(ResourceIndexDO) error

	AddLike(ResourceIndexDO) error
	RemoveLike(ResourceIndexDO) error

	AddRelatedModel(*RelatedResourceDO) error
	RemoveRelatedModel(*RelatedResourceDO) error

	AddRelatedDataset(*RelatedResourceDO) error
	RemoveRelatedDataset(*RelatedResourceDO) error

	UpdateProperty(*ProjectPropertyDO) error
}

func NewProjectRepository(mapper ProjectMapper) repository.Project {
	return project{mapper}
}

type project struct {
	mapper ProjectMapper
}

func (impl project) Save(p *domain.Project) (r domain.Project, err error) {
	if p.Id != "" {
		err = errors.New("must be a new project")

		return
	}

	v, err := impl.mapper.Insert(impl.toProjectDO(p))
	if err != nil {
		err = convertError(err)
	} else {
		r = *p
		r.Id = v
	}

	return
}

func (impl project) Delete(index *domain.ResourceIndex) (err error) {
	do := toResourceIndexDO(index)

	if err = impl.mapper.Delete(&do); err != nil {
		err = convertError(err)
	}

	return
}

func (impl project) Get(owner domain.Account, identity string) (r domain.Project, err error) {
	v, err := impl.mapper.Get(owner.Account(), identity)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toProject(&r)
	}

	return
}

func (impl project) GetByName(owner domain.Account, name domain.ResourceName) (
	r domain.Project, err error,
) {
	v, err := impl.mapper.GetByName(owner.Account(), name.ResourceName())
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toProject(&r)
	}

	return
}

func (impl project) FindUserProjects(opts []repository.UserResourceListOption) (
	[]domain.ProjectSummary, error,
) {
	do := make(map[string][]string)
	for i := range opts {
		do[opts[i].Owner.Account()] = opts[i].Ids
	}

	v, err := impl.mapper.ListUsersProjects(do)
	if err != nil {
		return nil, convertError(err)
	}

	r := make([]domain.ProjectSummary, len(v))
	for i := range v {
		if err = v[i].toProjectSummary(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl project) GetSummary(owner domain.Account, projectId string) (
	r repository.ProjectSummary, err error,
) {
	v, err := impl.mapper.GetSummary(owner.Account(), projectId)
	if err != nil {
		err = convertError(err)

		return
	}

	if r.ResourceSummary, err = v.toProject(); err == nil {
		r.Tags = v.Tags
	}

	return
}

func (impl project) GetSummaryByName(owner domain.Account, name domain.ResourceName) (
	domain.ResourceSummary, error,
) {
	v, err := impl.mapper.GetSummaryByName(owner.Account(), name.ResourceName())
	if err != nil {
		return domain.ResourceSummary{}, convertError(err)
	}

	return v.toProject()
}

func (impl project) toProjectDO(p *domain.Project) ProjectDO {
	do := ProjectDO{
		Id:        p.Id,
		Owner:     p.Owner.Account(),
		Name:      p.Name.ResourceName(),
		FL:        p.Name.FirstLetterOfName(),
		Type:      p.Type.ProjType(),
		CoverId:   p.CoverId.CoverId(),
		RepoType:  p.RepoType.RepoType(),
		Protocol:  p.Protocol.ProtocolName(),
		Training:  p.Training.TrainingPlatform(),
		Tags:      p.Tags,
		TagKinds:  p.TagKinds,
		RepoId:    p.RepoId,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Version:   p.Version,
		Hardware:  p.Hardware.Hardware(),
		BaseImage: p.BaseImage.BaseImage(),
	}

	if p.Desc != nil {
		do.Desc = p.Desc.ResourceDesc()
	}

	if p.Title != nil {
		do.Title = p.Title.ResourceTitle()
	}
	return do
}

type ProjectDO struct {
	Id            string
	Owner         string
	Name          string
	FL            byte
	Desc          string
	Title         string
	Type          string
	Level         int
	CoverId       string
	Protocol      string
	Training      string
	RepoType      string
	RepoId        string
	Tags          []string
	TagKinds      []string
	CreatedAt     int64
	UpdatedAt     int64
	Version       int
	LikeCount     int
	ForkCount     int
	DownloadCount int

	Hardware  string
	BaseImage string

	RelatedModels   []ResourceIndexDO
	RelatedDatasets []ResourceIndexDO
}

func (do *ProjectDO) toProject(r *domain.Project) (err error) {
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

	if r.Title, err = domain.NewResourceTitle(do.Title); err != nil {
		return
	}

	if r.Type, err = domain.NewProjType(do.Type); err != nil {
		return
	}

	if r.CoverId, err = domain.NewCoverId(do.CoverId); err != nil {
		return
	}

	if r.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
		return
	}

	if r.Protocol, err = domain.NewProtocolName(do.Protocol); err != nil {
		return
	}

	if r.Training, err = domain.NewTrainingPlatform(do.Training); err != nil {
		return
	}

	if r.RelatedModels, err = convertToResourceIndex(do.RelatedModels); err != nil {
		return
	}

	if r.RelatedDatasets, err = convertToResourceIndex(do.RelatedDatasets); err != nil {
		return
	}

	r.Level = domain.NewResourceLevelByNum(do.Level)
	r.RepoId = do.RepoId
	r.Tags = do.Tags
	r.TagKinds = do.TagKinds
	r.Version = do.Version
	r.CreatedAt = do.CreatedAt
	r.UpdatedAt = do.UpdatedAt
	r.LikeCount = do.LikeCount
	r.ForkCount = do.ForkCount
	r.DownloadCount = do.DownloadCount

	return
}

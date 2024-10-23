package repositoryimpl

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
)

type ProjectMapper interface {
	Insert(ProjectDO) (string, error)
	Delete(*repositories.ResourceIndexDO) error
	Get(string, string) (ProjectDO, error)
	GetByName(string, string) (ProjectDO, error)
	GetById(string) (ProjectDO, error)
	GetSummary(string, string) (ProjectResourceSummaryDO, error)
	GetSummaryByName(string, string) (repositories.ResourceSummaryDO, error)

	ListUsersProjects(map[string][]string) ([]ProjectSummaryDO, error)

	ListAndSortByUpdateTime(string, *repositories.ResourceListDO) ([]ProjectSummaryDO, int, error)
	ListAndSortByFirstLetter(string, *repositories.ResourceListDO) ([]ProjectSummaryDO, int, error)
	ListAndSortByDownloadCount(string, *repositories.ResourceListDO) ([]ProjectSummaryDO, int, error)

	ListGlobalAndSortByUpdateTime(*repositories.GlobalResourceListDO) ([]ProjectSummaryDO, int, error)
	ListGlobalAndSortByFirstLetter(*repositories.GlobalResourceListDO) ([]ProjectSummaryDO, int, error)
	ListGlobalAndSortByDownloadCount(*repositories.GlobalResourceListDO) ([]ProjectSummaryDO, int, error)

	Search(do *repositories.GlobalResourceListDO, topNum int) ([]repositories.ResourceSummaryDO, int, error)

	IncreaseFork(repositories.ResourceIndexDO) error
	IncreaseDownload(repositories.ResourceIndexDO) error

	AddLike(repositories.ResourceIndexDO) error
	RemoveLike(repositories.ResourceIndexDO) error

	AddRelatedModel(*repositories.RelatedResourceDO) error
	RemoveRelatedModel(*repositories.RelatedResourceDO) error

	AddRelatedDataset(*repositories.RelatedResourceDO) error
	RemoveRelatedDataset(*repositories.RelatedResourceDO) error

	UpdateProperty(*ProjectPropertyDO) error
}

func NewProjectRepository(mapper ProjectMapper) spacerepo.Project {
	return project{mapper}
}

type project struct {
	mapper ProjectMapper
}

func (impl project) Save(p *spacedomain.Project) (r spacedomain.Project, err error) {
	if p.Id != "" {
		err = errors.New("must be a new project")

		return
	}

	v, err := impl.mapper.Insert(impl.toProjectDO(p))
	if err != nil {
		err = repositories.ConvertError(err)
	} else {
		r = *p
		r.Id = v
	}

	return
}

func (impl project) Delete(index *domain.ResourceIndex) (err error) {
	do := repositories.ToResourceIndexDO(index)

	if err = impl.mapper.Delete(&do); err != nil {
		err = repositories.ConvertError(err)
	}

	return
}

func (impl project) Get(owner domain.Account, identity string) (r spacedomain.Project, err error) {
	v, err := impl.mapper.Get(owner.Account(), identity)
	if err != nil {
		err = repositories.ConvertError(err)
	} else {
		err = v.toProject(&r)
	}

	return
}

func (impl project) GetByName(owner domain.Account, name domain.ResourceName) (
	r spacedomain.Project, err error,
) {
	v, err := impl.mapper.GetByName(owner.Account(), name.ResourceName())
	if err != nil {
		err = repositories.ConvertError(err)
	} else {
		err = v.toProject(&r)
	}

	return
}

func (impl project) GetById(id domain.Identity) (
	r spacedomain.Project, err error,
) {
	v, err := impl.mapper.GetById(id.Identity())
	if err != nil {
		err = repositories.ConvertError(err)
	} else {
		err = v.toProject(&r)
	}

	return
}

func (impl project) FindUserProjects(opts []repository.UserResourceListOption) (
	[]spacedomain.ProjectSummary, error,
) {
	do := make(map[string][]string)
	for i := range opts {
		do[opts[i].Owner.Account()] = opts[i].Ids
	}

	v, err := impl.mapper.ListUsersProjects(do)
	if err != nil {
		return nil, repositories.ConvertError(err)
	}

	r := make([]spacedomain.ProjectSummary, len(v))
	for i := range v {
		if err = v[i].toProjectSummary(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl project) GetSummary(owner domain.Account, projectId string) (
	r spacerepo.ProjectSummary, err error,
) {
	v, err := impl.mapper.GetSummary(owner.Account(), projectId)
	if err != nil {
		err = repositories.ConvertError(err)

		return
	}

	if r.ResourceSummary, err = v.ToProject(); err == nil {
		r.Tags = v.Tags
	}

	return
}

func (impl project) GetSummaryByName(owner domain.Account, name domain.ResourceName) (
	domain.ResourceSummary, error,
) {
	v, err := impl.mapper.GetSummaryByName(owner.Account(), name.ResourceName())
	if err != nil {
		return domain.ResourceSummary{}, repositories.ConvertError(err)
	}

	return v.ToProject()
}

func (impl project) toProjectDO(p *spacedomain.Project) ProjectDO {
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

	RelatedModels   []repositories.ResourceIndexDO
	RelatedDatasets []repositories.ResourceIndexDO
}

func (do *ProjectDO) toProject(r *spacedomain.Project) (err error) {
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

	if r.RelatedModels, err = repositories.ConvertToResourceIndex(do.RelatedModels); err != nil {
		return
	}

	if r.RelatedDatasets, err = repositories.ConvertToResourceIndex(do.RelatedDatasets); err != nil {
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

func (impl project) ListGlobalAndSortByUpdateTime(
	option *repository.GlobalResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByUpdateTime,
	)
}

func (impl project) ListGlobalAndSortByFirstLetter(
	option *repository.GlobalResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByFirstLetter,
	)
}

func (impl project) ListGlobalAndSortByDownloadCount(
	option *repository.GlobalResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByDownloadCount,
	)
}

func (impl project) listGlobal(
	option *repository.GlobalResourceListOption,
	f func(*repositories.GlobalResourceListDO) ([]ProjectSummaryDO, int, error),
) (
	info spacerepo.UserProjectsInfo, err error,
) {
	return impl.doList(func() ([]ProjectSummaryDO, int, error) {
		do := repositories.ToGlobalResourceListDO(option)

		return f(&do)
	})
}

func (impl project) Search(option *repository.ResourceSearchOption) (
	repository.ResourceSearchResult, error,
) {
	r := repository.ResourceSearchResult{}

	do := repositories.SearchOptionToListDO(option)
	v, total, err := impl.mapper.Search(&do, option.TopNum)
	if err != nil {
		return r, err
	}

	items := make([]domain.ResourceSummary, len(v))
	for i := range v {
		if items[i].Owner, err = domain.NewAccount(v[i].Owner); err != nil {
			return r, err
		}

		if items[i].Name, err = domain.NewResourceName(v[i].Name); err != nil {
			return r, err
		}
	}

	r.Top = items
	r.Total = total

	return r, nil
}

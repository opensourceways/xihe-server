package repositoryimpl

import (
	"fmt"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
)

func equalQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}
func (impl project) IncreaseFork(p *domain.ResourceIndex) error {
	err := impl.mapper.IncreaseFork(repositories.ToResourceIndexDO(p))
	if err != nil {
		err = repositories.ConvertError(err)
	}

	return err
}

func (impl project) IncreaseDownload(index *domain.ResourceIndex) error {
	err := impl.mapper.IncreaseDownload(repositories.ToResourceIndexDO(index))
	if err != nil {
		err = repositories.ConvertError(err)
	}

	return err
}

func (impl project) AddLike(p *domain.ResourceIndex) error {
	err := impl.mapper.AddLike(repositories.ToResourceIndexDO(p))
	if err != nil {
		err = repositories.ConvertError(err)
	}

	return err
}

func (impl project) RemoveLike(p *domain.ResourceIndex) error {
	err := impl.mapper.RemoveLike(repositories.ToResourceIndexDO(p))
	if err != nil {
		err = repositories.ConvertError(err)
	}

	return err
}

func (impl project) AddRelatedModel(info *repository.RelatedResourceInfo) error {
	do := repositories.ToRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedModel(&do); err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}

func (impl project) RemoveRelatedModel(info *repository.RelatedResourceInfo) error {
	do := repositories.ToRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedModel(&do); err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}

func (impl project) AddRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := repositories.ToRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedDataset(&do); err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}

func (impl project) RemoveRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := repositories.ToRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedDataset(&do); err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}

func (impl project) UpdateProperty(info *spacerepo.ProjectPropertyUpdateInfo) error {
	p := &info.Property

	do := ProjectPropertyDO{
		ResourceToUpdateDO: repositories.ToResourceToUpdateDO(&info.ResourceToUpdate),

		FL:                p.Name.FirstLetterOfName(),
		Name:              p.Name.ResourceName(),
		CoverId:           p.CoverId.CoverId(),
		RepoType:          p.RepoType.RepoType(),
		Tags:              p.Tags,
		TagKinds:          p.TagKinds,
		CommitId:          p.CommitId,
		NoApplicationFile: p.NoApplicationFile,
		Exception:         p.Exception.Exception(),
	}

	if p.Desc != nil {
		do.Desc = p.Desc.ResourceDesc()
	}

	if p.Title != nil {
		do.Title = p.Title.ResourceTitle()
	}

	if p.Level != nil {
		do.Level = p.Level.Int()
	}

	if err := impl.mapper.UpdateProperty(&do); err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}

type ProjectPropertyDO struct {
	repositories.ResourceToUpdateDO

	FL                byte
	Name              string
	Desc              string
	Title             string
	Level             int
	CoverId           string
	RepoType          string
	Tags              []string
	TagKinds          []string
	CommitId          string
	NoApplicationFile bool
	Exception         string
}

func (impl project) ListAndSortByUpdateTime(
	owner domain.Account, option *repository.ResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByUpdateTime,
	)
}

// ListAndSortByUpdateTime Pg
func (adapter *projectAdapter) ListAndSortByUpdateTime(
	owner domain.Account, option *repository.ResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return adapter.list(
		owner, option, adapter.daoImpl.ListAndSortByUpdateTime,
	)
}

func (impl project) ListAndSortByFirstLetter(
	owner domain.Account, option *repository.ResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByUpdateTime,
	)
}

// ListAndSortByFirstLetter Pg
func (adapter *projectAdapter) ListAndSortByFirstLetter(
	owner domain.Account, option *repository.ResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return adapter.list(
		owner, option, adapter.daoImpl.ListAndSortByFirstLetter,
	)
}
func (impl project) ListAndSortByDownloadCount(
	owner domain.Account, option *repository.ResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByDownloadCount,
	)
}

// // ListAndSortByDownloadCount Pg
func (adapter *projectAdapter) ListAndSortByDownloadCount(
	owner domain.Account, option *repository.ResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return adapter.list(
		owner, option, adapter.daoImpl.ListAndSortByDownloadCount,
	)
}

func (impl project) list(
	owner domain.Account,
	option *repository.ResourceListOption,
	f func(string, *repositories.ResourceListDO) ([]ProjectSummaryDO, int, error),
) (
	info spacerepo.UserProjectsInfo, err error,
) {
	return impl.doList(func() ([]ProjectSummaryDO, int, error) {
		do := repositories.ToResourceListDO(option)

		return f(owner.Account(), &do)
	})
}

// list Pg
func (adapter *projectAdapter) list(
	owner domain.Account,
	option *repository.ResourceListOption,
	f func(string, *repositories.ResourceListDO) ([]ProjectSummaryDO, int, error),
) (
	info spacerepo.UserProjectsInfo, err error,
) {
	return adapter.doList(func() ([]ProjectSummaryDO, int, error) {
		do := repositories.ToResourceListDO(option)

		return f(owner.Account(), &do)
	})
}

func (impl project) doList(
	f func() ([]ProjectSummaryDO, int, error),
) (
	info spacerepo.UserProjectsInfo, err error,
) {
	v, total, err := f()
	if err != nil {
		err = repositories.ConvertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r := make([]spacedomain.ProjectSummary, len(v))
	for i := range v {
		if err = v[i].toProjectSummary(&r[i]); err != nil {
			r = nil

			return
		}
	}

	info.Projects = r
	info.Total = total

	return
}

// doList Pg
func (adapter *projectAdapter) doList(
	f func() ([]ProjectSummaryDO, int, error),
) (
	info spacerepo.UserProjectsInfo, err error,
) {
	v, total, err := f()
	if err != nil {
		err = repositories.ConvertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r := make([]spacedomain.ProjectSummary, len(v))
	for i := range v {
		if err = v[i].toProjectSummary(&r[i]); err != nil {
			r = nil

			return
		}
	}

	info.Projects = r
	info.Total = total

	return
}

type ProjectResourceSummaryDO struct {
	repositories.ResourceSummaryDO

	Tags []string
}

type ProjectSummaryDO struct {
	Id            string
	Owner         string
	Name          string
	Desc          string
	Title         string
	Level         int
	CoverId       string
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	ForkCount     int
	DownloadCount int
	Hardware      string
	Type          string
}

func (do *ProjectSummaryDO) toProjectSummary(r *spacedomain.ProjectSummary) (err error) {
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

	if r.CoverId, err = domain.NewCoverId(do.CoverId); err != nil {
		return
	}

	if r.Hardware, err = domain.NewHardware(do.Hardware, do.Type); err != nil {
		return
	}

	r.Level = domain.NewResourceLevelByNum(do.Level)
	r.Tags = do.Tags
	r.UpdatedAt = do.UpdatedAt
	r.LikeCount = do.LikeCount
	r.ForkCount = do.ForkCount
	r.DownloadCount = do.DownloadCount

	return
}

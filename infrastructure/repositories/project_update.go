package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl project) IncreaseFork(p *domain.ResourceIndex) error {
	err := impl.mapper.IncreaseFork(toResourceIndexDO(p))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) AddLike(p *domain.ResourceIndex) error {
	err := impl.mapper.AddLike(toResourceIndexDO(p))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) RemoveLike(p *domain.ResourceIndex) error {
	err := impl.mapper.RemoveLike(toResourceIndexDO(p))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) AddRelatedModel(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedModel(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) RemoveRelatedModel(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedModel(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) AddRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedDataset(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) RemoveRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedDataset(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) UpdateProperty(info *repository.ProjectPropertyUpdateInfo) error {
	p := &info.Property

	do := ProjectPropertyDO{
		ResourceToUpdateDO: toResourceToUpdateDO(&info.ResourceToUpdate),

		FL:       p.Name.FirstLetterOfName(),
		Name:     p.Name.ProjName(),
		Desc:     p.Desc.ResourceDesc(),
		CoverId:  p.CoverId.CoverId(),
		RepoType: p.RepoType.RepoType(),
		Tags:     p.Tags,
	}

	if err := impl.mapper.UpdateProperty(&do); err != nil {
		return convertError(err)
	}

	return nil
}

type ProjectPropertyDO struct {
	ResourceToUpdateDO

	FL       byte
	Name     string
	Desc     string
	CoverId  string
	RepoType string
	Tags     []string
}

func toRelatedResourceDO(info *repository.RelatedResourceInfo) RelatedResourceDO {
	return RelatedResourceDO{
		ResourceToUpdateDO: toResourceToUpdateDO(&info.ResourceToUpdate),
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

func toResourceToUpdateDO(info *repository.ResourceToUpdate) ResourceToUpdateDO {
	return ResourceToUpdateDO{
		Id:        info.Id,
		Owner:     info.Owner.Account(),
		Version:   info.Version,
		UpdatedAt: info.UpdatedAt,
	}
}

func (impl project) ListAndSortByUpdateTime(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByUpdateTime,
	)
}

func (impl project) ListAndSortByFirstLetter(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByFirstLetter,
	)
}

func (impl project) ListAndSortByDownloadCount(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByDownloadCount,
	)
}

func (impl project) list(
	owner domain.Account,
	option *repository.ResourceListOption,
	f func(string, *ResourceListDO) ([]ProjectSummaryDO, int, error),
) (
	info repository.UserProjectsInfo, err error,
) {
	do := toResourceListDO(option)

	v, total, err := f(owner.Account(), &do)
	if err != nil {
		err = convertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r := make([]domain.ProjectSummary, len(v))
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

type ProjectSummaryDO struct {
	Id            string
	Owner         string
	Name          string
	Desc          string
	CoverId       string
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	ForkCount     int
	DownloadCount int
}

func (do *ProjectSummaryDO) toProjectSummary(r *domain.ProjectSummary) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewProjName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewResourceDesc(do.Desc); err != nil {
		return
	}

	if r.CoverId, err = domain.NewConverId(do.CoverId); err != nil {
		return
	}

	r.Tags = do.Tags
	r.UpdatedAt = do.UpdatedAt
	r.LikeCount = do.LikeCount
	r.ForkCount = do.ForkCount
	r.DownloadCount = do.DownloadCount

	return
}

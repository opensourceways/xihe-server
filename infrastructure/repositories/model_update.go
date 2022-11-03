package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl model) AddLike(m *domain.ResourceIndex) error {
	err := impl.mapper.AddLike(toResourceIndexDO(m))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl model) RemoveLike(m *domain.ResourceIndex) error {
	err := impl.mapper.RemoveLike(toResourceIndexDO(m))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl model) AddRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedDataset(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl model) RemoveRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedDataset(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl model) AddRelatedProject(info *domain.ReverselyRelatedResourceInfo) error {
	do := toReverselyRelatedResourceInfoDO(info)

	if err := impl.mapper.AddRelatedProject(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl model) RemoveRelatedProject(info *domain.ReverselyRelatedResourceInfo) error {
	do := toReverselyRelatedResourceInfoDO(info)

	if err := impl.mapper.RemoveRelatedProject(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl model) UpdateProperty(info *repository.ModelPropertyUpdateInfo) error {
	p := &info.Property

	do := ModelPropertyDO{
		ResourceToUpdateDO: toResourceToUpdateDO(&info.ResourceToUpdate),

		FL:       p.Name.FirstLetterOfName(),
		Name:     p.Name.ModelName(),
		Desc:     p.Desc.ResourceDesc(),
		RepoType: p.RepoType.RepoType(),
		Tags:     p.Tags,
	}

	if err := impl.mapper.UpdateProperty(&do); err != nil {
		return convertError(err)
	}

	return nil
}

type ModelPropertyDO struct {
	ResourceToUpdateDO

	FL       byte
	Name     string
	Desc     string
	RepoType string
	Tags     []string
}

func (impl model) ListAndSortByUpdateTime(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserModelsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByUpdateTime,
	)
}

func (impl model) ListAndSortByFirstLetter(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserModelsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByFirstLetter,
	)
}

func (impl model) ListAndSortByDownloadCount(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserModelsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByDownloadCount,
	)
}

func (impl model) list(
	owner domain.Account,
	option *repository.ResourceListOption,
	f func(string, *ResourceListDO) ([]ModelSummaryDO, int, error),
) (
	info repository.UserModelsInfo, err error,
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

	r := make([]domain.ModelSummary, len(v))
	for i := range v {
		if err = v[i].toModelSummary(&r[i]); err != nil {
			r = nil

			return
		}
	}

	info.Models = r
	info.Total = total

	return
}

type ModelSummaryDO struct {
	Id            string
	Owner         string
	Name          string
	Desc          string
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

func (do *ModelSummaryDO) toModelSummary(r *domain.ModelSummary) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewModelName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewResourceDesc(do.Desc); err != nil {
		return
	}

	r.Tags = do.Tags
	r.UpdatedAt = do.UpdatedAt
	r.LikeCount = do.LikeCount
	r.DownloadCount = do.DownloadCount

	return
}

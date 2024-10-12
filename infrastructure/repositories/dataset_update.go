package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl dataset) IncreaseDownload(index *domain.ResourceIndex) error {
	err := impl.mapper.IncreaseDownload(ToResourceIndexDO(index))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl dataset) AddLike(d *domain.ResourceIndex) error {
	err := impl.mapper.AddLike(ToResourceIndexDO(d))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl dataset) RemoveLike(d *domain.ResourceIndex) error {
	err := impl.mapper.RemoveLike(ToResourceIndexDO(d))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl dataset) UpdateProperty(info *repository.DatasetPropertyUpdateInfo) error {
	p := &info.Property

	do := DatasetPropertyDO{
		ResourceToUpdateDO: ToResourceToUpdateDO(&info.ResourceToUpdate),

		FL:       p.Name.FirstLetterOfName(),
		Name:     p.Name.ResourceName(),
		RepoType: p.RepoType.RepoType(),
		Tags:     p.Tags,
		TagKinds: p.TagKinds,
	}

	if p.Desc != nil {
		do.Desc = p.Desc.ResourceDesc()
	}

	if p.Title != nil {
		do.Title = p.Title.ResourceTitle()
	}

	if err := impl.mapper.UpdateProperty(&do); err != nil {
		return convertError(err)
	}

	return nil
}

type DatasetPropertyDO struct {
	ResourceToUpdateDO

	FL       byte
	Name     string
	Desc     string
	Title    string
	RepoType string
	Tags     []string
	TagKinds []string
}

func (impl dataset) ListAndSortByUpdateTime(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByUpdateTime,
	)
}

func (impl dataset) ListAndSortByFirstLetter(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByFirstLetter,
	)
}

func (impl dataset) ListAndSortByDownloadCount(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByDownloadCount,
	)
}

func (impl dataset) list(
	owner domain.Account,
	option *repository.ResourceListOption,
	f func(string, *ResourceListDO) ([]DatasetSummaryDO, int, error),
) (
	info repository.UserDatasetsInfo, err error,
) {
	return impl.doList(func() ([]DatasetSummaryDO, int, error) {
		do := ToResourceListDO(option)

		return f(owner.Account(), &do)
	})

}

func (impl dataset) doList(
	f func() ([]DatasetSummaryDO, int, error),
) (
	info repository.UserDatasetsInfo, err error,
) {
	v, total, err := f()
	if err != nil {
		err = convertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r := make([]domain.DatasetSummary, len(v))
	for i := range v {
		if err = v[i].toDatasetSummary(&r[i]); err != nil {
			r = nil

			return
		}
	}

	info.Datasets = r
	info.Total = total

	return
}

func (impl dataset) AddRelatedProject(info *domain.ReverselyRelatedResourceInfo) error {
	do := toReverselyRelatedResourceInfoDO(info)

	if err := impl.mapper.AddRelatedProject(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl dataset) RemoveRelatedProject(info *domain.ReverselyRelatedResourceInfo) error {
	do := toReverselyRelatedResourceInfoDO(info)

	if err := impl.mapper.RemoveRelatedProject(&do); err != nil {
		return convertError(err)
	}

	return nil
}
func (impl dataset) AddRelatedModel(info *domain.ReverselyRelatedResourceInfo) error {
	do := toReverselyRelatedResourceInfoDO(info)

	if err := impl.mapper.AddRelatedModel(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl dataset) RemoveRelatedModel(info *domain.ReverselyRelatedResourceInfo) error {
	do := toReverselyRelatedResourceInfoDO(info)

	if err := impl.mapper.RemoveRelatedModel(&do); err != nil {
		return convertError(err)
	}

	return nil
}

type DatasetSummaryDO struct {
	Id            string
	Owner         string
	Name          string
	Desc          string
	Title         string
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

func (do *DatasetSummaryDO) toDatasetSummary(r *domain.DatasetSummary) (err error) {
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

	r.Tags = do.Tags
	r.UpdatedAt = do.UpdatedAt
	r.LikeCount = do.LikeCount
	r.DownloadCount = do.DownloadCount

	return
}

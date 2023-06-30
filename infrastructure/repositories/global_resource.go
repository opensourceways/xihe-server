package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func searchOptionToListDO(
	opt *repository.ResourceSearchOption,
) (do GlobalResourceListDO) {
	do.Name = opt.Name

	if opt.RepoType != nil {
		for i := range opt.RepoType {
			do.RepoType = append(do.RepoType, opt.RepoType[i].RepoType())
		}
	}

	return
}

type GlobalResourceListDO struct {
	ResourceListDO
	Level    int
	Tags     []string
	TagKinds []string
}

func toGlobalResourceListDO(
	opt *repository.GlobalResourceListOption,
) (do GlobalResourceListDO) {
	do.ResourceListDO = toResourceListDO(&opt.ResourceListOption)

	if opt.Level != nil {
		do.Level = opt.Level.Int()
	}
	do.Tags = opt.Tags
	do.TagKinds = opt.TagKinds

	return
}

func (impl project) ListGlobalAndSortByUpdateTime(
	option *repository.GlobalResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByUpdateTime,
	)
}

func (impl project) ListGlobalAndSortByFirstLetter(
	option *repository.GlobalResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByFirstLetter,
	)
}

func (impl project) ListGlobalAndSortByDownloadCount(
	option *repository.GlobalResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByDownloadCount,
	)
}

func (impl project) listGlobal(
	option *repository.GlobalResourceListOption,
	f func(*GlobalResourceListDO) ([]ProjectSummaryDO, int, error),
) (
	info repository.UserProjectsInfo, err error,
) {
	return impl.doList(func() ([]ProjectSummaryDO, int, error) {
		do := toGlobalResourceListDO(option)

		return f(&do)
	})
}

func (impl project) Search(option *repository.ResourceSearchOption) (
	repository.ResourceSearchResult, error,
) {
	r := repository.ResourceSearchResult{}

	do := searchOptionToListDO(option)
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

// Model
func (impl model) ListGlobalAndSortByUpdateTime(
	option *repository.GlobalResourceListOption,
) (repository.UserModelsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByUpdateTime,
	)
}

func (impl model) ListGlobalAndSortByFirstLetter(
	option *repository.GlobalResourceListOption,
) (repository.UserModelsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByFirstLetter,
	)
}

func (impl model) ListGlobalAndSortByDownloadCount(
	option *repository.GlobalResourceListOption,
) (repository.UserModelsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByDownloadCount,
	)
}

func (impl model) listGlobal(
	option *repository.GlobalResourceListOption,
	f func(*GlobalResourceListDO) ([]ModelSummaryDO, int, error),
) (
	info repository.UserModelsInfo, err error,
) {
	return impl.doList(func() ([]ModelSummaryDO, int, error) {
		do := toGlobalResourceListDO(option)

		return f(&do)
	})
}

func (impl model) Search(option *repository.ResourceSearchOption) (
	repository.ResourceSearchResult, error,
) {
	r := repository.ResourceSearchResult{}

	do := searchOptionToListDO(option)
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

// Dataset
func (impl dataset) ListGlobalAndSortByUpdateTime(
	option *repository.GlobalResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByUpdateTime,
	)
}

func (impl dataset) ListGlobalAndSortByFirstLetter(
	option *repository.GlobalResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByFirstLetter,
	)
}

func (impl dataset) ListGlobalAndSortByDownloadCount(
	option *repository.GlobalResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.listGlobal(
		option, impl.mapper.ListGlobalAndSortByDownloadCount,
	)
}

func (impl dataset) listGlobal(
	option *repository.GlobalResourceListOption,
	f func(*GlobalResourceListDO) ([]DatasetSummaryDO, int, error),
) (
	info repository.UserDatasetsInfo, err error,
) {
	return impl.doList(func() ([]DatasetSummaryDO, int, error) {
		do := toGlobalResourceListDO(option)

		return f(&do)
	})
}

func (impl dataset) Search(option *repository.ResourceSearchOption) (
	repository.ResourceSearchResult, error,
) {
	r := repository.ResourceSearchResult{}

	do := searchOptionToListDO(option)
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

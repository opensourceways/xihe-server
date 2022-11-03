package mongodb

import (
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

type globalProject struct {
	owner string
	*projectItem
}

func (col project) ListGlobalAndSortByUpdateTime(do *repositories.GlobalResourceListDO) (
	[]repositories.ProjectSummaryDO, int, error,
) {

	f := func(items []globalProject) []globalProject {
		v := make([]updateAtSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = updateAtSortData{
				id:       item.Id,
				index:    i,
				updateAt: item.UpdatedAt,
			}
		}

		v = updateAtSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]globalProject, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col project) ListGlobalAndSortByFirstLetter(do *repositories.GlobalResourceListDO) (
	[]repositories.ProjectSummaryDO, int, error,
) {

	f := func(items []globalProject) []globalProject {
		v := make([]firstLetterSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = firstLetterSortData{
				index:    i,
				letter:   item.FL,
				updateAt: item.UpdatedAt,
			}
		}

		v = firstLetterSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]globalProject, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col project) ListGlobalAndSortByDownloadCount(do *repositories.GlobalResourceListDO) (
	[]repositories.ProjectSummaryDO, int, error,
) {

	f := func(items []globalProject) []globalProject {
		v := make([]downloadSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = downloadSortData{
				index:    i,
				download: item.DownloadCount,
				updateAt: item.UpdatedAt,
			}
		}

		v = downloadSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]globalProject, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col project) listGlobalResource(
	do *repositories.GlobalResourceListDO,
	sortAndPagination func(items []globalProject) []globalProject,
) (r []repositories.ProjectSummaryDO, total int, err error) {
	var v []dProject

	err = listGlobalResourceWithoutSort(
		col.collectionName, do, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	items := col.toGlobalProjects(v)

	total = len(items)

	if items = sortAndPagination(items); len(items) == 0 {
		return
	}

	r = make([]repositories.ProjectSummaryDO, len(items))
	for i := range items {
		col.toProjectSummaryDO(items[i].owner, items[i].projectItem, &r[i])
	}

	return
}

func (col project) toGlobalProjects(v []dProject) []globalProject {
	n := 0
	for i := range v {
		n += len(v[i].Items)
	}

	k := 0
	r := make([]globalProject, n)

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		for j := range items {
			r[k] = globalProject{
				owner:       owner,
				projectItem: &items[j],
			}
			k++
		}
	}

	return r
}

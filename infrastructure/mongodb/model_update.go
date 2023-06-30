package mongodb

import (
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func (col model) IncreaseDownload(index repositories.ResourceIndexDO) error {
	return updateResourceDownloadNum(col.collectionName, &index, 1)
}

func (col model) AddLike(r repositories.ResourceIndexDO) error {
	return updateResourceLikeNum(col.collectionName, &r, 1)
}

func (col model) RemoveLike(r repositories.ResourceIndexDO) error {
	return updateResourceLikeNum(col.collectionName, &r, -1)
}

func (col model) AddRelatedDataset(do *repositories.RelatedResourceDO) error {
	return updateRelatedResource(col.collectionName, fieldDatasets, true, do)
}

func (col model) RemoveRelatedDataset(do *repositories.RelatedResourceDO) error {
	return updateRelatedResource(col.collectionName, fieldDatasets, false, do)
}

func (col model) AddRelatedProject(do *repositories.ReverselyRelatedResourceInfoDO) error {
	return updateReverselyRelatedResource(col.collectionName, fieldProjects, true, do)
}

func (col model) RemoveRelatedProject(do *repositories.ReverselyRelatedResourceInfoDO) error {
	return updateReverselyRelatedResource(col.collectionName, fieldProjects, false, do)
}

func (col model) ListAndSortByUpdateTime(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ModelSummaryDO, int, error) {

	f := func(items []modelItem) []modelItem {
		v := make([]updateAtSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = updateAtSortData{
				id:       item.Id,
				level:    item.Level,
				index:    i,
				updateAt: item.UpdatedAt,
			}
		}

		v = updateAtSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]modelItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col model) ListAndSortByFirstLetter(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ModelSummaryDO, int, error) {

	f := func(items []modelItem) []modelItem {
		v := make([]firstLetterSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = firstLetterSortData{
				index:    i,
				level:    item.Level,
				letter:   item.FL,
				updateAt: item.UpdatedAt,
			}
		}

		v = firstLetterSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]modelItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col model) ListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ModelSummaryDO, int, error) {

	f := func(items []modelItem) []modelItem {
		v := make([]downloadSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = downloadSortData{
				index:    i,
				level:    item.Level,
				download: item.DownloadCount,
				updateAt: item.UpdatedAt,
			}
		}

		v = downloadSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]modelItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col model) listResource(
	owner string,
	do *repositories.ResourceListDO,
	sortAndPagination func(items []modelItem) []modelItem,
) (r []repositories.ModelSummaryDO, total int, err error) {
	var v []dModel

	err = listResourceWithoutSort(
		col.collectionName, owner, do, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 || len(v[0].Items) == 0 {
		return
	}

	items := v[0].Items
	total = len(items)

	items = sortAndPagination(items)
	if len(items) == 0 {
		return
	}

	r = make([]repositories.ModelSummaryDO, len(items))
	for i := range items {
		col.toModelSummaryDO(owner, &items[i], &r[i])
	}

	return
}

func (col model) summaryFields() []string {
	return []string{
		fieldId, fieldName, fieldDesc, fieldTitle, fieldTags, fieldFirstLetter,
		fieldUpdatedAt, fieldLikeCount, fieldDownloadCount, fieldLevel,
	}
}

func (col model) toModelSummaryDO(owner string, item *modelItem, do *repositories.ModelSummaryDO) {
	*do = repositories.ModelSummaryDO{
		Id:            item.Id,
		Owner:         owner,
		Name:          item.Name,
		Desc:          item.Desc,
		Tags:          item.Tags,
		Title:         item.Title,
		UpdatedAt:     item.UpdatedAt,
		LikeCount:     item.LikeCount,
		DownloadCount: item.DownloadCount,
	}
}

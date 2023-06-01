package mongodb

import (
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func (col dataset) IncreaseDownload(index repositories.ResourceIndexDO) error {
	return updateResourceDownloadNum(col.collectionName, &index, 1)
}

func (col dataset) AddLike(r repositories.ResourceIndexDO) error {
	return updateResourceLikeNum(col.collectionName, &r, 1)
}

func (col dataset) RemoveLike(r repositories.ResourceIndexDO) error {
	return updateResourceLikeNum(col.collectionName, &r, -1)
}

func (col dataset) AddRelatedProject(do *repositories.ReverselyRelatedResourceInfoDO) error {
	return updateReverselyRelatedResource(col.collectionName, fieldProjects, true, do)
}

func (col dataset) RemoveRelatedProject(do *repositories.ReverselyRelatedResourceInfoDO) error {
	return updateReverselyRelatedResource(col.collectionName, fieldProjects, false, do)
}

func (col dataset) AddRelatedModel(do *repositories.ReverselyRelatedResourceInfoDO) error {
	return updateReverselyRelatedResource(col.collectionName, fieldModels, true, do)
}

func (col dataset) RemoveRelatedModel(do *repositories.ReverselyRelatedResourceInfoDO) error {
	return updateReverselyRelatedResource(col.collectionName, fieldModels, false, do)
}

func (col dataset) ListAndSortByUpdateTime(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.DatasetSummaryDO, int, error) {

	f := func(items []datasetItem) []datasetItem {
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

		r := make([]datasetItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col dataset) ListAndSortByFirstLetter(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.DatasetSummaryDO, int, error) {

	f := func(items []datasetItem) []datasetItem {
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

		r := make([]datasetItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col dataset) ListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.DatasetSummaryDO, int, error) {

	f := func(items []datasetItem) []datasetItem {
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

		r := make([]datasetItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col dataset) listResource(
	owner string,
	do *repositories.ResourceListDO,
	sortAndPagination func(items []datasetItem) []datasetItem,
) (r []repositories.DatasetSummaryDO, total int, err error) {
	var v []dDataset

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

	r = make([]repositories.DatasetSummaryDO, len(items))
	for i := range items {
		col.toDatasetSummaryDO(owner, &items[i], &r[i])
	}

	return
}

func (col dataset) summaryFields() []string {
	return []string{
		fieldId, fieldName, fieldDesc, fieldTitle, fieldTags, fieldFirstLetter,
		fieldUpdatedAt, fieldLikeCount, fieldDownloadCount, fieldLevel,
	}
}

func (col dataset) toDatasetSummaryDO(owner string, item *datasetItem, do *repositories.DatasetSummaryDO) {
	*do = repositories.DatasetSummaryDO{
		Id:            item.Id,
		Owner:         owner,
		Name:          item.Name,
		Desc:          item.Desc,
		Title:         item.Title,
		Tags:          item.Tags,
		UpdatedAt:     item.UpdatedAt,
		LikeCount:     item.LikeCount,
		DownloadCount: item.DownloadCount,
	}
}

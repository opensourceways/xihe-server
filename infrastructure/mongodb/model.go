package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewModelMapper(name string) repositories.ModelMapper {
	return model{name}
}

type model struct {
	collectionName string
}

func (col model) newDoc(owner string) error {
	docFilter := resourceOwnerFilter(owner)

	doc := bson.M{
		fieldOwner: owner,
		fieldItems: bson.A{},
	}

	f := func(ctx context.Context) error {
		_, err := cli.newDocIfNotExist(
			ctx, col.collectionName, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil && isDBError(err) {
		return err
	}

	return nil
}

func (col model) Insert(do repositories.ModelDO) (identity string, err error) {
	if identity, err = col.insert(do); err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert

	if err = col.newDoc(do.Owner); err == nil {
		if identity, err = col.insert(do); err != nil && isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return
}

func (col model) insert(do repositories.ModelDO) (identity string, err error) {
	identity = newId()

	do.Id = identity
	doc, err := col.toModelDoc(&do)
	if err != nil {
		return
	}
	doc[fieldVersion] = 0
	doc[fieldLikeCount] = 0
	doc[fieldDownloadCount] = 0
	doc[fieldDatasets] = bson.A{}
	doc[fieldProjects] = bson.A{}

	err = insertResource(col.collectionName, do.Owner, do.Name, doc)

	return
}

func (col model) Delete(do *repositories.ResourceIndexDO) error {
	return deleteResource(col.collectionName, do)
}

func (col model) UpdateProperty(do *repositories.ModelPropertyDO) error {
	p := &ModelPropertyItem{
		FL:       do.FL,
		Name:     do.Name,
		Desc:     do.Desc,
		RepoType: do.RepoType,
		Tags:     do.Tags,
		TagKinds: do.TagKinds,
	}

	p.setDefault()

	return updateResourceProperty(col.collectionName, &do.ResourceToUpdateDO, p)
}

func (col model) Get(owner, identity string) (do repositories.ModelDO, err error) {
	var v []dModel

	if err = getResourceById(col.collectionName, owner, identity, &v); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toModelDO(owner, &v[0].Items[0], &do)

	return
}

func (col model) GetByName(owner, name string) (do repositories.ModelDO, err error) {
	var v []dModel

	if err = getResourceByName(col.collectionName, owner, name, &v); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toModelDO(owner, &v[0].Items[0], &do)

	return
}

func (col model) GetSummaryByName(owner, name string) (
	do repositories.ResourceSummaryDO, err error,
) {
	var v []dModel

	err = getResourceSummaryByName(col.collectionName, owner, name, &v)
	if err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	item := &v[0].Items[0]
	do.Id = item.Id
	do.Name = name
	do.Owner = owner
	do.RepoId = item.RepoId
	do.RepoType = item.RepoType

	return
}

func (col model) ListUsersModels(opts map[string][]string) (
	r []repositories.ModelSummaryDO, err error,
) {
	var v []dModel

	err = listUsersResources(col.collectionName, opts, &v)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]repositories.ModelSummaryDO, 0, len(v))

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		dos := make([]repositories.ModelSummaryDO, len(items))

		for j := range items {
			col.toModelSummaryDO(owner, &items[j], &dos[j])
		}

		r = append(r, dos...)
	}

	return
}

func (col model) ListSummary(opts map[string][]string) (
	r []repositories.ResourceSummaryDO, err error,
) {
	var v []dModel

	err = listUsersResourcesSummary(col.collectionName, opts, &v)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]repositories.ResourceSummaryDO, 0, len(v))

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		dos := make([]repositories.ResourceSummaryDO, len(items))

		for j := range items {
			item := &items[j]

			dos[j] = repositories.ResourceSummaryDO{
				Id:       item.Id,
				Name:     item.Name,
				Owner:    owner,
				RepoId:   item.RepoId,
				RepoType: item.RepoType,
			}
		}

		r = append(r, dos...)
	}

	return
}

func (col model) toModelDoc(do *repositories.ModelDO) (bson.M, error) {
	docObj := modelItem{
		Id:        do.Id,
		RepoId:    do.RepoId,
		Protocol:  do.Protocol,
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
		ModelPropertyItem: ModelPropertyItem{
			FL:       do.FL,
			Name:     do.Name,
			Desc:     do.Desc,
			Title:    do.Title,
			RepoType: do.RepoType,
			Tags:     do.Tags,
			TagKinds: do.TagKinds,
		},
	}

	docObj.ModelPropertyItem.setDefault()

	return genDoc(docObj)
}

func (col model) toModelDO(owner string, item *modelItem, do *repositories.ModelDO) {
	*do = repositories.ModelDO{
		Id:            item.Id,
		Owner:         owner,
		Name:          item.Name,
		Desc:          item.Desc,
		Title:         item.Title,
		Protocol:      item.Protocol,
		RepoType:      item.RepoType,
		RepoId:        item.RepoId,
		Tags:          item.Tags,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		Version:       item.Version,
		LikeCount:     item.LikeCount,
		DownloadCount: item.DownloadCount,

		RelatedDatasets: toResourceIndexDO(item.RelatedDatasets),
		RelatedProjects: toResourceIndexDO(item.RelatedProjects),
	}
}

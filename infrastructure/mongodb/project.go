package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func projectDocFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func projectItemFilter(name string) bson.M {
	return bson.M{
		fieldName: name,
	}
}

func arrayFilterById(identity string) bson.M {
	return bson.M{
		fieldId: identity,
	}
}

func NewProjectMapper(name string) repositories.ProjectMapper {
	return project{name}
}

type project struct {
	collectionName string
}

func (col project) New(owner string) error {
	docFilter := projectDocFilter(owner)

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

	if err := withContext(f); err != nil && !errors.Is(err, errDocExists) {
		return err
	}

	return nil
}

func (col project) Insert(do repositories.ProjectDO) (identity string, err error) {
	identity = newId()

	do.Id = identity
	doc, err := col.toProjectDoc(&do)
	if err != nil {
		return
	}
	doc[fieldVersion] = 0

	docFilter := projectDocFilter(do.Owner)

	appendElemMatchToFilter(
		fieldItems, false,
		projectItemFilter(do.Name), docFilter,
	)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName,
			fieldItems, docFilter, doc,
		)
	}

	err = withContext(f)

	if errors.Is(err, errDocNotExists) {
		err = repositories.NewErrorDuplicateCreating(err)
	}

	return
}

func (col project) Update(do repositories.ProjectDO) error {
	doc, err := col.toProjectDoc(&do)
	if err != nil {
		return err
	}

	docFilter := projectDocFilter(do.Owner)

	updated := false

	f := func(ctx context.Context) error {
		b, err := cli.updateArrayElem(
			ctx, col.collectionName, fieldItems,
			docFilter, arrayFilterById(do.Id), doc, do.Version,
		)

		updated = b

		return err
	}

	if err := withContext(f); err != nil {
		return err
	}

	if !updated {
		return repositories.NewErrorConcurrentUpdating(errors.New("no update"))
	}

	return nil
}

func (col project) Get(owner, identity string) (do repositories.ProjectDO, err error) {
	var v []dProject

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			projectDocFilter(owner), arrayFilterById(identity),
			bson.M{fieldItems: 1}, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toProjectDO(owner, &v[0].Items[0], &do)

	return
}

func (col project) List(owner string, do repositories.ProjectListDO) (
	r []repositories.ProjectDO, err error,
) {
	var v []dProject

	f := func(ctx context.Context) error {
		return cli.getArraysElemsByCustomizedCond(
			ctx, col.collectionName, projectDocFilter(owner),
			map[string]func() bson.M{
				fieldItems: func() bson.M {
					conds := bson.A{}

					if do.RepoType != "" {
						conds = append(conds, eqCondForArrayElem(
							fieldRepoType, do.RepoType,
						))
					}

					if do.Name != "" {
						conds = append(conds, matchCondForArrayElem(
							fieldName, do.Name,
						))
					}

					return condForArrayElem(conds)
				},
			},
			bson.M{fieldItems: 1}, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		return
	}

	items := v[0].Items
	r = make([]repositories.ProjectDO, len(items))
	for i := range items {
		col.toProjectDO(owner, &items[i], &r[i])
	}

	return
}

func (col project) toProjectDoc(do *repositories.ProjectDO) (bson.M, error) {
	docObj := projectItem{
		Id:       do.Id,
		Name:     do.Name,
		Desc:     do.Desc,
		Type:     do.Type,
		CoverId:  do.CoverId,
		Protocol: do.Protocol,
		Training: do.Training,
		RepoType: do.RepoType,
		RepoId:   do.RepoId,
		Tags:     do.Tags,
	}

	return genDoc(docObj)
}

func (col project) toProjectDO(owner string, item *projectItem, do *repositories.ProjectDO) {
	*do = repositories.ProjectDO{
		Id:       item.Id,
		Owner:    owner,
		Name:     item.Name,
		Desc:     item.Desc,
		Type:     item.Type,
		CoverId:  item.CoverId,
		Protocol: item.Protocol,
		Training: item.Training,
		RepoType: item.RepoType,
		RepoId:   item.RepoId,
		Tags:     item.Tags,
		Version:  item.Version,
	}
}

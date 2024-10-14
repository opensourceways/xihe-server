package mongodb

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"github.com/opensourceways/xihe-server/spaceapp/infrastructure/repositoryimpl"
)

func NewInferenceMapper(name string) repositoryimpl.InferenceMapper {
	return inference{name}
}

func inferenceDocFilter(owner, projectId, commit string) bson.M {
	return bson.M{
		fieldPId:    projectId,
		fieldOwner:  owner,
		fieldCommit: commit,
	}
}

type inference struct {
	collectionName string
}

func (col inference) newDoc(do *repositoryimpl.InferenceDO) error {
	docFilter := inferenceDocFilter(do.ProjectOwner, do.ProjectId, do.LastCommit)

	doc := bson.M{
		fieldOwner:   do.ProjectOwner,
		fieldPId:     do.ProjectId,
		fieldName:    do.ProjectName,
		fieldCommit:  do.LastCommit,
		fieldItems:   bson.A{},
		fieldVersion: 0,
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

func (col inference) Insert(do *repositoryimpl.InferenceDO, version int) (
	identity string, err error,
) {
	identity, err = col.insert(do, version)
	if err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert

	if err = col.newDoc(do); err == nil {
		identity, err = col.insert(do, version)
		if err != nil && isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return
}

func (col inference) insert(do *repositoryimpl.InferenceDO, version int) (identity string, err error) {
	identity = newId()

	doc := bson.M{fieldId: identity}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, col.collectionName,
			inferenceDocFilter(do.ProjectOwner, do.ProjectId, do.LastCommit),
			bson.M{fieldItems: doc}, mongoCmdPush, version,
		)
	}

	err = withContext(f)

	return
}

func (col inference) UpdateDetail(
	index *repositoryimpl.InferenceIndexDO,
	detail *repositoryimpl.InferenceDetailDO,
) error {
	data := inferenceItem{
		Expiry:    detail.Expiry,
		Error:     detail.Error,
		AccessURL: detail.AccessURL,
	}

	doc, err := genDoc(data)
	if err != nil {
		return err
	}

	logrus.Debugf(
		"update inference(%s/%s) to %v",
		index.Project.Id, index.Id, doc,
	)

	f := func(ctx context.Context) error {
		_, err := cli.modifyArrayElemWithoutVersion(
			ctx, col.collectionName, fieldItems,
			col.docFilter(index),
			resourceIdFilter(index.Id), doc, mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

func (col inference) docFilter(index *repositoryimpl.InferenceIndexDO) bson.M {
	return inferenceDocFilter(
		index.Project.Owner, index.Project.Id, index.LastCommit,
	)
}

func (col inference) Get(index *repositoryimpl.InferenceIndexDO) (
	r repositoryimpl.InferenceSummaryDO, err error,
) {
	var v []dInference

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			col.docFilter(index),
			resourceIdFilter(index.Id),
			nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toInferenceSummaryDO(&v[0].Items[0], &r)

	return
}

func (col inference) List(index *repositories.ResourceIndexDO, lastCommit string) (
	[]repositoryimpl.InferenceSummaryDO, int, error,
) {
	var v dInference

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			inferenceDocFilter(index.Owner, index.Id, lastCommit),
			bson.M{
				fieldVersion: 1,
				fieldItems:   1,
			}, &v)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	t := v.Items
	r := make([]repositoryimpl.InferenceSummaryDO, len(t))

	for i := range t {
		col.toInferenceSummaryDO(&t[i], &r[i])
	}

	return r, v.Version, nil
}

func (col inference) toInferenceSummaryDO(doc *inferenceItem, r *repositoryimpl.InferenceSummaryDO) {
	r.Id = doc.Id
	r.Error = doc.Error
	r.Expiry = doc.Expiry
	r.AccessURL = doc.AccessURL
}

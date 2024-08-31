package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewFinetuneMapper(name string) repositories.FinetuneMapper {
	return finetuneCol{name}
}

type finetuneCol struct {
	collectionName string
}

func (col finetuneCol) Insert(do *repositories.UserFinetuneDO, version int) (
	identity string, err error,
) {
	identity, err = col.insert(do, version)
	if err != nil && isDocNotExists(err) {
		err = repositories.NewErrorDuplicateCreating(err)
	}

	return
}

func (col finetuneCol) insert(do *repositories.UserFinetuneDO, version int) (identity string, err error) {
	identity = newId()
	do.Id = identity

	doc, err := col.toFinetuneDoc(do)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, col.collectionName,
			resourceOwnerFilter(do.Owner),
			bson.M{fieldItems: doc}, mongoCmdPush, version,
		)
	}

	err = withContext(f)

	return
}

func (col finetuneCol) Delete(index *repositories.FinetuneIndexDO) error {
	f := func(ctx context.Context) error {
		return cli.pullArrayElem(
			ctx, col.collectionName, fieldItems,
			resourceOwnerFilter(index.Owner),
			resourceIdFilter(index.Id),
		)
	}

	return withContext(f)
}

func (col finetuneCol) Get(index *repositories.FinetuneIndexDO) (
	do repositories.FinetuneDetailDO, err error,
) {
	var v []dFinetune

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			resourceOwnerFilter(index.Owner),
			resourceIdFilter(index.Id),
			nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)
	} else {
		col.toFinetuneDetailDO(&v[0].Items[0], &do)
	}

	return
}

func (col finetuneCol) List(user string) (do repositories.UserFinetunesDO, version int, err error) {
	var v dFinetune

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			resourceOwnerFilter(user),
			bson.M{
				fieldExpiry:                     1,
				fieldVersion:                    1,
				subfieldOfItems(fieldId):        1,
				subfieldOfItems(fieldName):      1,
				subfieldOfItems(fieldDetail):    1,
				subfieldOfItems(fieldCreatedAt): 1,
			}, &v)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return
	}

	version = v.Version
	do.Expiry = v.Expiry

	t := v.Items
	if len(t) == 0 {
		return
	}

	r := make([]repositories.FinetuneSummaryDO, len(t))
	for i := range t {
		col.toFinetuneSummaryDO(&t[i], &r[i])
	}
	do.Data = r

	return
}

func (col finetuneCol) GetJob(index *repositories.FinetuneIndexDO) (
	do repositories.FinetuneJobDO, err error,
) {
	var v []dFinetune

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			resourceOwnerFilter(index.Owner),
			resourceIdFilter(index.Id),
			bson.M{
				subfieldOfItems(fieldJob):    1,
				subfieldOfItems(fieldDetail): 1,
			},
			&v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	item := &v[0].Items[0]
	do.JobId = item.Job.JobId
	do.Endpoint = item.Job.Endpoint
	do.Status = item.JobDetail.Status

	return
}

func (col finetuneCol) UpdateJobInfo(
	index *repositories.FinetuneIndexDO, do *repositories.FinetuneJobInfoDO,
) error {
	v := dFinetuneJobInfo{
		Endpoint: do.Endpoint,
		JobId:    do.JobId,
	}

	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := cli.modifyArrayElemWithoutVersion(
			ctx, col.collectionName, fieldItems,
			resourceOwnerFilter(index.Owner),
			resourceIdFilter(index.Id),
			bson.M{fieldJob: doc}, mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

func (col finetuneCol) UpdateJobDetail(
	index *repositories.FinetuneIndexDO, do *repositories.FinetuneJobDetailDO,
) error {
	v := dFinetuneJobDetail{
		Duration: do.Duration,
		Error:    do.Error,
		Status:   do.Status,
	}

	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := cli.modifyArrayElemWithoutVersion(
			ctx, col.collectionName, fieldItems,
			resourceOwnerFilter(index.Owner),
			resourceIdFilter(index.Id),
			bson.M{fieldDetail: doc}, mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

const (
	lastArrayElemIndex int = -1
)

func NewTrainingMapper(name string) repositories.TrainingMapper {
	return training{name}
}

func trainingDocFilter(owner, projectId string) bson.M {
	return bson.M{
		fieldOwner: owner,
		fieldPId:   projectId,
	}
}

type training struct {
	collectionName string
}

func (col training) newDoc(do *repositories.UserTrainingDO) error {
	docFilter := trainingDocFilter(do.Owner, do.ProjectId)

	doc := bson.M{
		fieldOwner:   do.Owner,
		fieldPId:     do.ProjectId,
		fieldName:    do.ProjectName,
		fieldRId:     do.ProjectRepoId,
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

func (col training) Insert(do *repositories.UserTrainingDO, version int) (
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

func (col training) insert(do *repositories.UserTrainingDO, version int) (identity string, err error) {
	identity = newId()
	do.Id = identity

	doc, err := col.toTrainingDoc(do)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, col.collectionName,
			trainingDocFilter(do.Owner, do.ProjectId),
			bson.M{fieldItems: doc}, mongoCmdPush, version,
		)
	}

	err = withContext(f)

	return
}

func (col training) List(user, projectId string) ([]repositories.TrainingSummaryDO, int, error) {
	var v dTraining

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			trainingDocFilter(user, projectId),
			bson.M{
				fieldVersion:                    1,
				subfieldOfItems(fieldId):        1,
				subfieldOfItems(fieldName):      1,
				subfieldOfItems(fieldDesc):      1,
				subfieldOfItems(fieldDetail):    1,
				subfieldOfItems(fieldCreatedAt): 1,
			}, &v)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	t := v.Items
	r := make([]repositories.TrainingSummaryDO, len(t))

	for i := range t {
		col.toTrainingSummary(&t[i], &r[i])
	}

	return r, v.Version, nil
}

func (col training) Delete(info *repositories.TrainingIndexDO) error {
	f := func(ctx context.Context) error {
		return cli.pullArrayElem(
			ctx, col.collectionName, fieldItems,
			trainingDocFilter(info.User, info.ProjectId),
			resourceIdFilter(info.TrainingId),
		)
	}

	return withContext(f)
}

func (col training) Get(info *repositories.TrainingIndexDO) (
	do repositories.TrainingDetailDO, err error,
) {
	var v []dTraining

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			trainingDocFilter(info.User, info.ProjectId),
			resourceIdFilter(info.TrainingId),
			bson.M{
				fieldRId:   1,
				fieldName:  1,
				fieldItems: 1,
			},
			&v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)
	} else {
		col.toTrainingDetailDO(&v[0], &do)
	}

	return
}

func (col training) GetTrainingConfig(info *repositories.TrainingIndexDO) (
	do repositories.TrainingConfigDO, err error,
) {
	var v []dTraining

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			trainingDocFilter(info.User, info.ProjectId),
			resourceIdFilter(info.TrainingId),
			bson.M{
				fieldName:  1,
				fieldRId:   1,
				fieldItems: 1,
			},
			&v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)
	} else {
		col.toTrainingConfigDO(&v[0], &do)
	}

	return
}

func (col training) GetLastTrainingConfig(res *repositories.ResourceIndexDO) (
	do repositories.TrainingConfigDO, err error,
) {
	var v *dTraining

	f := func(ctx context.Context) error {
		return cli.getArrayElemSlice(
			ctx, col.collectionName,
			trainingDocFilter(res.Owner, res.Id),
			fieldItems, lastArrayElemIndex,
			v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if v == nil || len(v.Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)
	} else {
		col.toTrainingConfigDO(v, &do)
	}

	return
}

func (col training) GetJobInfo(info *repositories.TrainingIndexDO) (
	do repositories.TrainingJobInfoDO, err error,
) {
	var v []dTraining

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			trainingDocFilter(info.User, info.ProjectId),
			resourceIdFilter(info.TrainingId),
			bson.M{
				subfieldOfItems(fieldJob): 1,
			},
			&v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)
	} else {
		col.toTrainingJobInfoDO(&v[0].Items[0].Job, &do)
	}

	return
}

func (col training) UpdateJobInfo(info *repositories.TrainingIndexDO, job *repositories.TrainingJobInfoDO) error {
	v := dJobInfo{
		Endpoint:  job.Endpoint,
		JobId:     job.JobId,
		LogDir:    job.LogDir,
		AimDir:    job.AimDir,
		OutputDir: job.OutputDir,
	}

	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := cli.modifyArrayElemWithoutVersion(
			ctx, col.collectionName, fieldItems,
			trainingDocFilter(info.User, info.ProjectId),
			resourceIdFilter(info.TrainingId),
			bson.M{fieldJob: doc}, mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

func (col training) UpdateJobDetail(info *repositories.TrainingIndexDO, detail *repositories.TrainingJobDetailDO) error {
	v := dJobDetail{
		Duration:   detail.Duration,
		Error:      detail.Error,
		Status:     detail.Status,
		LogPath:    detail.LogPath,
		AimPath:    detail.AimPath,
		OutputPath: detail.OutputPath,
	}

	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := cli.modifyArrayElemWithoutVersion(
			ctx, col.collectionName, fieldItems,
			trainingDocFilter(info.User, info.ProjectId),
			resourceIdFilter(info.TrainingId),
			bson.M{fieldDetail: doc}, mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

func (col training) GetJobDetail(info *repositories.TrainingIndexDO) (
	do repositories.TrainingJobDetailDO, endpoint string, err error,
) {
	var v []dTraining

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			trainingDocFilter(info.User, info.ProjectId),
			resourceIdFilter(info.TrainingId),
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
	} else {
		item := &v[0].Items[0]
		endpoint = item.Job.Endpoint
		col.toTrainingJobDetailDO(&item.JobDetail, &do)
	}

	return
}

func (col training) toTrainingSummary(t *trainingItem, s *repositories.TrainingSummaryDO) {
	*s = repositories.TrainingSummaryDO{
		Id:        t.Id,
		Name:      t.Name,
		Desc:      t.Desc,
		Error:     t.JobDetail.Error,
		Status:    t.JobDetail.Status,
		Duration:  t.JobDetail.Duration,
		CreatedAt: t.CreatedAt,
	}
}

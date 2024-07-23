package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

func NewAICCFinetuneRepo(m mongodbClient) repository.AICCFinetune {
	return aiccFinetuneRepoImpl{m}
}

type aiccFinetuneRepoImpl struct {
	cli mongodbClient
}

func aiccFinetuneDocFilter(user, model string) bson.M {
	return bson.M{
		fieldUser:  user,
		fieldModel: model,
	}
}

func (impl aiccFinetuneRepoImpl) newDoc(a *domain.AICCFinetune) error {
	docFilter := aiccFinetuneDocFilter(a.User.Account(), a.Model.ModelName())

	doc := bson.M{
		fieldUser:    a.User,
		fieldModel:   a.Model,
		fieldItems:   bson.A{},
		fieldVersion: 0,
	}

	f := func(ctx context.Context) error {
		_, err := impl.cli.NewDocIfNotExist(
			ctx, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil {
		if !impl.cli.IsDocNotExists(err) {
			return nil
		}
		return err
	}

	return nil

}

func (impl aiccFinetuneRepoImpl) Save(a *domain.AICCFinetune, version int) (id string, err error) {
	id = primitive.NewObjectID().Hex()
	a.Id = id
	if err = impl.newDoc(a); err != nil {
		return
	}
	docFilter := aiccFinetuneDocFilter(a.User.Account(), a.Model.ModelName())
	doc, err := impl.genAICCFinetuneDoc(a)
	if err != nil {
		return
	}

	doc[fieldVersion] = 0

	f := func(ctx context.Context) error {
		err = impl.cli.UpdateDoc(
			ctx, docFilter,
			bson.M{fieldItems: doc}, mongoCmdPush, version,
		)
		return err
	}

	if err = withContext(f); err != nil {
		return
	}

	return
}

func (impl aiccFinetuneRepoImpl) Delete(info *domain.AICCFinetuneIndex) error {
	f := func(ctx context.Context) error {
		return impl.cli.PullArrayElem(
			ctx, fieldItems,
			aiccFinetuneDocFilter(info.User.Account(), info.Model.ModelName()),
			bson.M{fieldId: info.FinetuneId},
		)
	}

	return withContext(f)
}

func (impl aiccFinetuneRepoImpl) Get(info *domain.AICCFinetuneIndex) (obj domain.AICCFinetune, err error) {
	var v []dAICCFinetune

	f := func(ctx context.Context) error {
		return impl.cli.GetArrayElem(
			ctx,
			fieldItems,
			aiccFinetuneDocFilter(info.User.Account(), info.Model.ModelName()),
			bson.M{fieldId: info.FinetuneId},
			bson.M{
				fieldUser:  1,
				fieldModel: 1,
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
		err = v[0].toAICCFinetuneDO(&obj)
	}

	return
}

func subfieldOfItems(k string) string {
	return fieldItems + "." + k
}

func (impl aiccFinetuneRepoImpl) List(user types.Account, Model domain.ModelName) (
	r []domain.AICCFinetuneSummary, version int, err error,
) {
	var v dAICCFinetune

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx,
			aiccFinetuneDocFilter(user.Account(), Model.ModelName()),
			bson.M{
				fieldVersion:                    1,
				subfieldOfItems(fieldId):        1,
				subfieldOfItems(fieldName):      1,
				subfieldOfItems(fieldDesc):      1,
				subfieldOfItems(fieldCreatedAt): 1,
				subfieldOfItems(fieldDetail):    1,
				subfieldOfItems(fieldTask):      1,
			}, &v)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	t := v.Items
	r = make([]domain.AICCFinetuneSummary, len(t))

	for i := range t {
		if err = t[i].toAICCFinetuneSummary(&r[i]); err != nil {
			return r, v.Version, err
		}
	}

	return r, v.Version, nil
}

func (impl aiccFinetuneRepoImpl) SaveJob(info *domain.AICCFinetuneIndex, job *domain.JobInfo) error {
	v := dJobInfo{

		Endpoint:  job.Endpoint,
		JobId:     job.JobId,
		LogDir:    job.LogDir,
		OutputDir: job.OutputDir,
	}

	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := impl.cli.ModifyArrayElem(
			ctx,
			fieldItems,
			aiccFinetuneDocFilter(info.User.Account(), info.Model.ModelName()),
			bson.M{fieldId: info.FinetuneId},
			bson.M{fieldJob: doc},
			mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

func (impl aiccFinetuneRepoImpl) GetJob(info *domain.AICCFinetuneIndex) (job domain.JobInfo, err error) {
	var v []dAICCFinetune

	f := func(ctx context.Context) error {
		return impl.cli.GetArrayElem(
			ctx,
			fieldItems,
			aiccFinetuneDocFilter(info.User.Account(), info.Model.ModelName()),
			bson.M{fieldId: info.FinetuneId},
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
		err = v[0].Items[0].Job.toAICCFinetuneJobInfo(&job)
	}

	return
}

func (impl aiccFinetuneRepoImpl) GetJobDetail(info *domain.AICCFinetuneIndex) (
	job domain.JobDetail, endpoint string, err error,
) {
	var v []dAICCFinetune

	f := func(ctx context.Context) error {
		return impl.cli.GetArrayElem(
			ctx,
			fieldItems,
			aiccFinetuneDocFilter(info.User.Account(), info.Model.ModelName()),
			bson.M{fieldId: info.FinetuneId},
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
		err = v[0].Items[0].JobDetail.toAICCFinetuneJobDetail(&job)
		endpoint = v[0].Items[0].Job.Endpoint
	}

	return
}

func (impl aiccFinetuneRepoImpl) UpdateJobDetail(info *domain.AICCFinetuneIndex, detail *domain.JobDetail) error {
	v := dJobDetail{
		Duration:   detail.Duration,
		Error:      detail.Error,
		Status:     detail.Status,
		LogPath:    detail.LogPath,
		OutputPath: detail.OutputPath,
	}

	doc, err := genDoc(v)
	if err != nil {
		return err
	}
	filter := bson.M{
		fieldUser:  info.User.Account(),
		fieldModel: info.Model.ModelName(),
	}

	f := func(ctx context.Context) error {
		_, err := impl.cli.ModifyArrayElem(
			ctx,
			fieldItems,
			filter,
			bson.M{fieldId: info.FinetuneId},
			bson.M{fieldDetail: doc},
			mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

func (repo aiccFinetuneRepoImpl) genAICCFinetuneDoc(p *domain.AICCFinetune) (bson.M, error) {

	c := p.AICCFinetuneConfig

	docObj := aiccFinetuneItem{
		Id:              p.Id,
		Name:            c.Name.FinetuneName(),
		Desc:            c.Desc.FinetuneDesc(),
		Task:            p.Task.FinetuneTask(),
		CreatedAt:       p.CreatedAt,
		Model:           p.Model.ModelName(),
		Env:             repo.toKeyValueDoc(p.Env),
		Hyperparameters: repo.toKeyValueDoc(p.Hyperparameters),
	}
	return genDoc(docObj)
}

func (col aiccFinetuneRepoImpl) toKeyValueDoc(kv []domain.KeyValue) []dKeyValue {
	n := len(kv)
	if n == 0 {
		return nil
	}

	r := make([]dKeyValue, n)

	for i := range kv {
		r[i].Key = kv[i].Key.CustomizedKey()
		r[i].Value = kv[i].Value.CustomizedValue()
	}

	return r
}

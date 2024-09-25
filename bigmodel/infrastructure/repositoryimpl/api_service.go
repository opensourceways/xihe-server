package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func NewApiService(m mongodbClient) repository.ApiService {
	return &apiServiceRepoImpl{m}
}

type apiServiceRepoImpl struct {
	cli mongodbClient
}

func (impl *apiServiceRepoImpl) ApplyApi(d *domain.UserApiRecord) error {
	doc, err := toApiApplyDoc(d)
	doc[fieldVersion] = 1
	filter := bson.M{fieldModelName: d.ModelName, fieldUserName: d.User.Account()}

	f := func(ctx context.Context) error {
		_, err = impl.cli.NewDocIfNotExist(
			ctx,
			filter,
			doc,
		)
		return err
	}

	if err = withContext(f); err != nil {
		return err
	}

	return nil
}

func (impl *apiServiceRepoImpl) GetApiByUserModel(user types.Account, model domain.ModelName) (
	u domain.UserApiRecord, err error,
) {
	v := new(dApiApply)

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx,
			bson.M{fieldUser: user.Account(), fieldModelName: model.ModelName()},
			nil,
			&v,
		)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}
		return
	}

	err = v.toUserApiRecord(&u)
	return
}

func (impl *apiServiceRepoImpl) GetApiByUser(user types.Account) (u []domain.UserApiRecord, err error) {
	var v []dApiApply

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx,
			bson.M{fieldUser: user.Account()},
			nil,
			&v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	u = make([]domain.UserApiRecord, len(v))
	for i := range v {
		if err = v[i].toUserApiRecord(&u[i]); err != nil {
			return
		}
	}

	return
}

func (impl *apiServiceRepoImpl) AddApiCallCount(user types.Account, model domain.ModelName, version int) error {
	f := func(ctx context.Context) error {
		return impl.cli.UpdateIncDoc(
			ctx,
			bson.M{fieldUser: user.Account(), fieldModelName: model.ModelName()},
			bson.M{fieldCallCount: 1},
			version,
		)
	}

	if err := withContext(f); err != nil {
		return err
	}

	return nil
}

func (a *dApiApply) toUserApiRecord(d *domain.UserApiRecord) (err error) {
	if d.User, err = types.NewAccount(a.User); err != nil {
		return
	}

	if d.ModelName, err = domain.NewModelName(a.ModelName); err != nil {
		return
	}

	d.ApplyAt = a.ApplyAt
	d.Enabled = a.Enabled
	d.Token = a.Token
	d.UpdateAt = a.UpdateAt
	d.Version = a.Version
	return
}

func (impl *apiServiceRepoImpl) UpdateToken(user types.Account,
	model domain.ModelName, token string, date string, version int) (err error) {
	filter := bson.M{fieldUser: user.Account(), fieldModelName: model.ModelName()}

	f := func(ctx context.Context) error {
		err = impl.cli.UpdateDoc(
			ctx,
			filter,
			bson.M{fieldToken: token, fieldUpdateAt: date},
			mongoCmdSet,
			version,
		)
		return err
	}

	if err = withContext(f); err != nil {
		return err
	}

	return nil
}

func toApiApplyDoc(d *domain.UserApiRecord) (bson.M, error) {
	return genDoc(dApiApply{
		User:      d.User.Account(),
		ModelName: d.ModelName.ModelName(),
		ApplyAt:   d.ApplyAt,
		UpdateAt:  d.UpdateAt,
		Token:     d.Token,
		Enabled:   d.Enabled,
		CallCount: 0,
	})
}

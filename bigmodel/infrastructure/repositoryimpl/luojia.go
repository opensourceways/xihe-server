package repositoryimpl

import (
	"context"
	"errors"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	commoninfra "github.com/opensourceways/xihe-server/common/infrastructure"
	types "github.com/opensourceways/xihe-server/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func NewLuoJiaRepo(m mongodbClient) repository.LuoJia {
	return luojiaRepoImpl{m}
}

type luojiaRepoImpl struct {
	cli mongodbClient
}

func (impl luojiaRepoImpl) Save(ur *domain.UserLuoJiaRecord) (r domain.LuoJiaRecord, err error) {
	if ur.Id != "" {
		err = errors.New("must be a new luojia")

		return
	}

	ur.Id = newId()

	luojiaDoc, err := impl.genLuoJiaDoc(ur)
	if err != nil {
		err = commoninfra.ConvertError(err)

		return
	}

	luojiaItemDoc, err := impl.genLuoJiaItemDoc(&ur.LuoJiaRecord)
	if err != nil {
		err = commoninfra.ConvertError(err)

		return
	}

	var id string
	f := func(ctx context.Context) (err error) {
		// 1. create owner
		filterOwner := bson.M{
			fieldOwner: ur.User,
		}
		if id, err = impl.cli.NewDocIfNotExist(ctx, filterOwner, luojiaDoc); err != nil {
			if !impl.cli.IsDocExists(err) {
				return
			}
		}

		// 2. insert into array, owner filter
		if err = impl.cli.PushElemToLimitedArray(ctx, fieldItems, 10, filterOwner, luojiaItemDoc); err != nil {
			return
		}

		return
	}

	if err = withContext(f); err != nil {
		return
	}

	r = ur.LuoJiaRecord
	r.Id = id

	return
}

func (impl luojiaRepoImpl) List(user types.Account) (r []domain.LuoJiaRecord, err error) {
	var v dLuoJia
	f := func(ctx context.Context) (err error) {
		filterOwner := bson.M{
			fieldOwner: user.Account(),
		}

		if err = impl.cli.GetDoc(ctx, filterOwner, bson.M{fieldItems: 1}, &v); err != nil {
			if impl.cli.IsDocNotExists(err) {
				err = nil

				return
			}

			return
		}

		return
	}

	if err = withContext(f); err != nil {
		return
	}

	return impl.toLuoJiaRecordList(v.Items), nil
}

func (impl luojiaRepoImpl) genLuoJiaDoc(d *domain.UserLuoJiaRecord) (bson.M, error) {
	return genDoc(dLuoJia{
		Owner: d.User.Account(),
		Items: []luojiaItem{
			{
				Id:        d.Id,
				CreatedAt: d.CreatedAt,
			},
		},
	})
}

func (impl luojiaRepoImpl) genLuoJiaItemDoc(d *domain.LuoJiaRecord) (bson.M, error) {
	return genDoc(luojiaItem{
		Id:        d.Id,
		CreatedAt: d.CreatedAt,
	})

}

func (impl luojiaRepoImpl) toLuoJiaRecordList(v []luojiaItem) (r []domain.LuoJiaRecord) {
	r = make([]domain.LuoJiaRecord, len(v))

	for i := range v {
		v[i].toLuoJiaRecord(&r[i])
	}

	return
}

func (v *luojiaItem) toLuoJiaRecord(d *domain.LuoJiaRecord) {
	*d = domain.LuoJiaRecord{
		Id:        v.Id,
		CreatedAt: v.CreatedAt,
	}
}

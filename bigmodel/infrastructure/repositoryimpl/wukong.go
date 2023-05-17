package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func wukongIdFilter(identity string) bson.M {
	return bson.M{
		fieldId: identity,
	}
}

func wukongOwnerFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

type WuKongPictureListOptionDO = repository.WuKongPictureListOption

func NewWuKongRepo(m mongodbClient) repository.WuKong {
	return &wukongRepoImpl{m}
}

type wukongRepoImpl struct {
	cli mongodbClient
}

func (impl *wukongRepoImpl) ListSamples(sid string, nums []int) ([]string, error) {
	project := bson.M{
		fieldSamples: bson.M{"$filter": bson.M{
			"input": "$" + fieldSamples,
			"cond": func() bson.M {
				return impl.cli.InCondForArrayElem(fieldNum, nums)
			}(),
		}},
	}

	pipeline := bson.A{
		bson.M{"$match": wukongIdFilter(sid)},
		bson.M{"$project": project},
	}

	var v []dWuKong

	err := withContext(func(ctx context.Context) error {
		col := impl.cli.Collection()
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	})
	if err != nil || len(v) == 0 {
		return nil, err
	}

	doc := v[0]

	r := make([]string, len(doc.Samples))
	for i := range doc.Samples {
		r[i] = doc.Samples[i].Name
	}

	return r, nil
}

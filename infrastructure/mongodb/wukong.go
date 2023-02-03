package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewWuKongMapper(name string) repositories.WuKongMapper {
	return wukong{
		collectionName: name,
	}
}

type wukong struct {
	collectionName string
}

func (col wukong) ListSamples(docId string, nums []int) ([]string, error) {
	project := bson.M{
		fieldSamples: bson.M{"$filter": bson.M{
			"input": "$" + fieldSamples,
			"cond": func() bson.M {
				return inCondForArrayElem(fieldNum, nums)
			}(),
		}},
	}

	pipeline := bson.A{
		bson.M{"$match": resourceIdFilter(docId)},
		bson.M{"$project": project},
	}

	var v []dWuKong

	err := withContext(func(ctx context.Context) error {
		col := cli.collection(col.collectionName)
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

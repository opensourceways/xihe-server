package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func NewCompetitionRepo(m mongodbClient) repository.Competition {
	return competitionRepoImpl{m}
}

func condFieldOfArrayElem(key string) string {
	return "$$this." + key
}

func valueInCondForArrayElem(key string, value interface{}) bson.M {
	return bson.M{"$in": bson.A{value, condFieldOfArrayElem(key)}}
}

func condForArrayElem(conds bson.A) bson.M {
	n := len(conds)
	if n > 1 {
		return bson.M{"$and": conds}
	}

	if n == 1 {
		return conds[0].(bson.M)
	}

	return bson.M{
		"$toBool": 1,
	}
}

type competitionRepoImpl struct {
	cli mongodbClient
}

func (impl competitionRepoImpl) docFilter(cid string) bson.M {
	return bson.M{
		fieldId: cid,
	}
}

func (impl competitionRepoImpl) FindCompetition(opt *repository.CompetitionGetOption) (
	c domain.Competition, err error,
) {
	var v dCompetition

	f := func(ctx context.Context) error {
		filter := bson.M{}

		filter[fieldId] = opt.CompetitionId

		if opt.Lang != nil {
			filter[fieldLanguage] = opt.Lang.Language()
		}
		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return
	}

	err = v.toCompetition(&c)

	return
}

func (impl competitionRepoImpl) FindScoreOrder(cid string) (
	domain.CompetitionScoreOrder, error,
) {
	var v dCompetition

	f := func(ctx context.Context) error {
		filter := impl.docFilter(cid)

		return impl.cli.GetDoc(ctx, filter, bson.M{"order": 1}, &v)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return nil, err
	}

	return domain.NewCompetitionScoreOrder(v.SmallerOk), nil
}

func (impl competitionRepoImpl) FindCompetitions(opt *repository.CompetitionListOption) (
	[]domain.CompetitionSummary, error,
) {
	var v []dCompetition

	f := func(ctx context.Context) error {
		filter := bson.M{}
		if opt.Status != nil {
			filter[fieldStatus] = opt.Status.CompetitionStatus()
		}

		if len(opt.CompetitionIds) > 0 {
			filter[fieldId] = bson.M{
				"$in": opt.CompetitionIds,
			}
		}

		if opt.Tag != nil {
			filter[fieldTags] = opt.Tag.CompetitionTag()
		}

		if opt.Lang != nil {
			filter[fieldLanguage] = opt.Lang.Language()
		}

		return impl.cli.GetDocs(ctx, filter, nil, &v)
	}

	if err := withContext(f); err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]domain.CompetitionSummary, len(v))
	for i := range v {
		if err := v[i].toCompetitionSummary(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

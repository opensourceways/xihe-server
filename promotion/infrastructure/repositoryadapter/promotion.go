package repositoryadapter

import (
	"context"

	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	fieldType      = "type"
	fieldWay       = "way"
	fieldStartTime = "start_time"
	fieldEndTime   = "end_time"
	fieldTags      = "tags"

	operatorGreaterThanEq = "$gte"
	operatorGreaterThan   = "$gt"
	operatorLessThanEq    = "$lte"
	operatorLessThan      = "$lt"
	operatorMatch         = "$match"
	operatorAddFields     = "$addFields"
	operatorSort          = "$sort"
	operatorSkip          = "$skip"
	operatorLimit         = "$limit"
	operatorProject       = "$project"
	operatorIn            = "$in"
)

func PromotionAdapter(cli mongodbClient) repository.Promotion {
	return &promotionAdapter{cli}
}

type promotionAdapter struct {
	cli mongodbClient
}

func (impl *promotionAdapter) FindById(id string) (domain.Promotion, error) {
	var do promotionDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx, docIdFilter(id),
			nil, &do,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return domain.Promotion{}, err
	}

	return do.toPromotion()
}

func (impl *promotionAdapter) FindAll() (prs []domain.Promotion, err error) {
	var dos []promotionDO

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(
			ctx, nil,
			nil, &dos,
		)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return nil, err
	}

	prs = make([]domain.Promotion, len(dos))
	for i := range dos {
		if prs[i], err = dos[i].toPromotion(); err != nil {
			return
		}
	}

	return
}

func (impl *promotionAdapter) UserRegister(promotionid string, user types.Account, origin domain.Origin,
	version int) error {
	regUserDO := RegUserDO{
		User:      user.Account(),
		CreatedAt: utils.Now(),
		Origin:    origin.Oringn(),
	}
	doc, err := regUserDO.doc()
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		return impl.cli.PushElemArrayWithVersion(
			ctx, fieldRegUsers,
			docIdFilter(promotionid),
			doc, version, nil,
		)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}

		return err
	}

	return nil
}

func promotionsQueryToFilter(query *repository.PromotionsQuery) primitive.M {
	filter := primitive.M{}

	if query.Type != nil {
		filter[fieldType] = query.Type.PromotionType()
	}

	if query.Status != nil {
		now := utils.Now()
		switch query.Status.PromotionStatus() {
		case domain.PromotionStatusInProgress:
			filter[fieldStartTime] = primitive.M{
				operatorLessThanEq: now,
			}
			filter[fieldEndTime] = primitive.M{
				operatorGreaterThanEq: now,
			}
		case domain.PromotionStatusOver:
			filter[fieldStartTime] = primitive.M{
				operatorLessThan: now,
			}
			filter[fieldEndTime] = primitive.M{
				operatorLessThan: now,
			}
		default:
		}
	}

	if query.Way != nil {
		filter[fieldWay] = query.Way.PromotionWay()
	}

	if len(query.Tags) > 0 {
		filter[fieldTags] = primitive.M{
			operatorIn: query.Tags,
		}
	}

	return filter
}

func promotionsQueryToPipeline(query *repository.PromotionsQuery) mongo.Pipeline {
	/*
		promotionsQueryToPipeline convert a query to mongodb pipeline. There are 5 steps will be executed.

		step 1: $match give the conditions which the documents should match
		step 2: $addFields generate a temporary field `over` to indicate the status of promotion is over
		step 3: $sort show the sorting fields and sorting direction
		step 4: $skip show the number of documents to skip
		step 5: $limit show the number of documents to return

		Entire MQL is:
		db.collection.aggregate([
			{
				$match: {
					type: "xxx",
					start_time: {$lt: 0123456789},
					end_time: {$gt: 0123456789},
					way: "xxx",
					tags: {$in: ["xxx"]} // optional
				}
			},
			{
				$addFields: {
					over: {
						$gt: [0123456789, "$end_time"],
					}
				}
			},
			{
				$sort: {
					over: 1,
					priority: -1,
					start_time: -1,
				}
			},
			{ $skip: 0 },
			{ $limit: 12 }
		])
	*/

	match := promotionsQueryToFilter(query)

	const fieldOver = "over"
	fields := bson.M{
		fieldOver: bson.M{operatorGreaterThanEq: bson.A{utils.Now(), "$" + fieldEndTime}},
	}

	sort := bson.D{{Key: fieldOver, Value: 1}}
	for _, item := range query.Sort {
		switch item[1] {
		case repository.SortAsc:
			sort = append(sort, bson.E{Key: item[0], Value: 1})
		case repository.SortDesc:
			sort = append(sort, bson.E{Key: item[0], Value: -1})
		}
	}

	return mongo.Pipeline{
		{{Key: operatorMatch, Value: match}},
		{{Key: operatorAddFields, Value: fields}},
		{{Key: operatorSort, Value: sort}},
		{{Key: operatorSkip, Value: query.Offset}},
		{{Key: operatorLimit, Value: query.Limit}},
	}
}

func (impl *promotionAdapter) FindByCustom(query *repository.PromotionsQuery) ([]domain.Promotion, error) {
	var promotionsDO []promotionDO

	pipeline := promotionsQueryToPipeline(query)

	f := func(ctx context.Context) error {
		return impl.cli.Aggregate(
			ctx, pipeline, nil,
			&promotionsDO,
		)
	}
	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return nil, err
	}

	promotions := make([]domain.Promotion, len(promotionsDO))
	for i := range promotionsDO {
		promotion, err := promotionsDO[i].toPromotion()
		if err != nil {
			return promotions, err
		}
		promotions[i] = promotion
	}

	return promotions, nil
}

func (impl *promotionAdapter) Count(query *repository.PromotionsQuery) (int64, error) {
	filter := promotionsQueryToFilter(query)
	var count int64
	f := func(ctx context.Context) (err error) {
		count, err = impl.cli.Count(ctx, filter, nil)
		return
	}
	err := withContext(f)

	return count, err
}

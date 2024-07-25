package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"github.com/opensourceways/xihe-server/utils"
)

func NewAIQuestionMapper(name, pool string) repositories.AIQuestionMapper {
	return aiquestion{
		collectionName: name,
		poolCollection: pool,
	}
}

type aiquestion struct {
	collectionName string
	poolCollection string
}

func (col aiquestion) GetCompetitorAndSubmission(qid, competitor string) (
	isCompetitor bool, score int,
	do repositories.QuestionSubmissionDO,
	err error,
) {
	docFilter, err := objectIdFilter(qid)
	if err != nil {
		return
	}

	var v []dAIQuestion

	f := func(ctx context.Context) error {
		filter := bson.M{fieldAccount: competitor}

		return cli.getArraysElem(
			ctx, col.collectionName,
			docFilter,
			map[string]bson.M{
				fieldCompetitors: filter,
				fieldSubmissions: filter,
			},
			bson.M{
				fieldCompetitors + "." + fieldAccount: 1,
				fieldSubmissions:                      1,
			},
			&v,
		)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	doc := &v[0]
	if len(doc.Competitors) == 0 {
		return
	}

	isCompetitor = true

	if len(doc.Submissions) == 0 {
		return
	}

	today := utils.Date()
	items := doc.Submissions
	for i := range items {
		item := &items[i]

		if score < item.Score {
			score = item.Score
		}

		if item.Date == today {
			col.toQuestionSubmissionDo(&do, item)
		}
	}

	return
}

func (col aiquestion) GetResult(qid string) (
	do []repositories.QuestionSubmissionInfoDO, err error,
) {
	docFilter, err := objectIdFilter(qid)
	if err != nil {
		return
	}

	var v dAIQuestion

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, docFilter,
			bson.M{
				fieldSubmissions + "." + fieldAccount: 1,
				fieldSubmissions + "." + fieldScore:   1,
			},
			&v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	items := v.Submissions
	do = make([]repositories.QuestionSubmissionInfoDO, len(items))

	for i := range items {
		col.toQuestionSubmissionInfoDo(&do[i], &items[i])
	}

	return
}

func (col aiquestion) SaveCompetitor(qid string, do *repositories.CompetitorInfoDO) error {
	v := new(DCompetitorInfo)
	toCompetitorInfoDOC(v, do)

	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	docFilter, err := objectIdFilter(qid)
	if err != nil {
		return err
	}

	appendElemMatchToFilter(
		fieldCompetitors, false, bson.M{fieldAccount: do.Account}, docFilter,
	)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName, fieldCompetitors, docFilter, doc,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			return repositories.NewErrorDuplicateCreating(err)
		}
	}

	return err
}

func (col aiquestion) InsertSubmission(qid string, do *repositories.QuestionSubmissionDO) (
	sid string, err error,
) {
	sid = newId()

	do.Id = sid
	v := new(dQuestionSubmission)
	col.toQuestionSubmissionDoc(do, v)

	doc, err := genDoc(v)
	if err != nil {
		return
	}
	doc[fieldVersion] = 0

	docFilter, err := objectIdFilter(qid)
	if err != nil {
		return
	}

	appendElemMatchToFilter(
		fieldSubmissions, false,
		bson.M{
			fieldAccount: do.Account,
			fieldDate:    do.Date,
		},
		docFilter,
	)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName, fieldSubmissions, docFilter, doc,
		)
	}

	err = withContext(f)

	return
}

func (col aiquestion) UpdateSubmission(qid string, do *repositories.QuestionSubmissionDO) error {
	v := new(dQuestionSubmission)
	col.toQuestionSubmissionDoc(do, v)

	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	docFilter, err := objectIdFilter(qid)
	if err != nil {
		return err
	}

	updated := false

	f := func(ctx context.Context) error {
		updated, err = cli.updateArrayElem(
			ctx, col.collectionName, fieldSubmissions,
			docFilter,
			resourceIdFilter(do.Id),
			doc, do.Version, 0,
		)

		return err
	}

	// #nosec G104 -- this is a false positive
	if err = withContext(f); err != nil {
		return err
	}

	if !updated {
		return repositories.NewErrorConcurrentUpdating(
			errors.New("no update"),
		)
	}

	return nil
}

func (col aiquestion) GetSubmission(qid, competitor, date string) (
	do repositories.QuestionSubmissionDO, err error,
) {
	docFilter, err := objectIdFilter(qid)
	if err != nil {
		return
	}

	var v []dAIQuestion

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldSubmissions,
			docFilter,
			bson.M{
				fieldAccount: competitor,
				fieldDate:    date,
			},
			nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		err = errDocNotExists

		return
	}

	if len(v[0].Submissions) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toQuestionSubmissionDo(&do, &v[0].Submissions[0])

	return
}

func (col aiquestion) GetQuestions(poolId string, choice, completion []int) (
	choices []repositories.ChoiceQuestionDO,
	completions []repositories.CompletionQuestionDO, err error,
) {
	docFilter, err := objectIdFilter(poolId)
	if err != nil {
		return
	}

	project := bson.M{
		fieldChoices: bson.M{"$filter": bson.M{
			"input": "$" + fieldChoices,
			"cond": func() bson.M {
				return inCondForArrayElem(fieldNum, choice)
			}(),
		}},
		fieldCompletions: bson.M{"$filter": bson.M{
			"input": "$" + fieldCompletions,
			"cond": func() bson.M {
				return inCondForArrayElem(fieldNum, completion)
			}(),
		}},
	}

	pipeline := bson.A{
		bson.M{"$match": docFilter},
		bson.M{"$project": project},
	}

	var v []dQuestionPool

	err = withContext(func(ctx context.Context) error {
		col := cli.collection(col.poolCollection)
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	})
	if err != nil || len(v) == 0 {
		return
	}

	questions := v[0]

	choices = make([]repositories.ChoiceQuestionDO, len(questions.Choices))
	for i := range questions.Choices {
		col.toChoiceQuestionDO(&choices[i], &questions.Choices[i])
	}

	completions = make([]repositories.CompletionQuestionDO, len(questions.Completions))
	for i := range questions.Completions {
		col.toCompletionQuestionDO(&completions[i], &questions.Completions[i])
	}

	return
}

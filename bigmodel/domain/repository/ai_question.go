package repository

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type AIQuestion interface {
	GetResult(string) ([]domain.QuestionSubmissionInfo, error)

	GetCompetitorAndSubmission(string, types.Account) (
		isCompetitor bool, highestScore int,
		submission domain.QuestionSubmission,
		err error,
	)

	GetQuestions(pool string, choice, completion []int) (
		[]domain.ChoiceQuestion, []domain.CompletionQuestion, error,
	)

	GetSubmission(qid string, user types.Account, date string) (
		domain.QuestionSubmission, error,
	)

	SaveSubmission(qid string, v *domain.QuestionSubmission) (string, error)
}

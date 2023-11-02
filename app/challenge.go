package app

import (
	"encoding/base64"
	"errors"
	"sort"
	"strings"

	"github.com/opensourceways/xihe-server/agreement/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/challenge"
	"github.com/opensourceways/xihe-server/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

type ChallengeService interface {
	Apply(*CompetitorApplyCmd) error
	GetCompetitor(domain.Account) (ChallengeCompetitorInfoDTO, error)
	GetAIQuestions(domain.Account) (AIQuestionDTO, string, error)
	SubmitAIQuestionAnswer(domain.Account, *AIQuestionAnswerSubmitCmd) (int, string, error)
	GetRankingList() ([]ChallengeRankingDTO, error)
}

type challengeService struct {
	comptitions []domain.CompetitionIndex
	aiQuestion  challenge.AIQuestionInfo
	delimiter   string

	competitionRepo repository.Competition
	aiQuestionRepo  repository.AIQuestion
	helper          challenge.Challenge
	encryption      utils.SymmetricEncryption
	userRepo        userrepo.User
}

func NewChallengeService(
	competitionRepo repository.Competition,
	aiQuestionRepo repository.AIQuestion,
	helper challenge.Challenge,
	encryption utils.SymmetricEncryption,
	userRepo userrepo.User,
) ChallengeService {
	v := helper.GetChallenge()

	s := &challengeService{
		competitionRepo: competitionRepo,
		aiQuestionRepo:  aiQuestionRepo,
		encryption:      encryption,
		helper:          helper,
		delimiter:       ",-;",
		userRepo:        userRepo,
	}

	s.comptitions = make([]domain.CompetitionIndex, len(v.Competition))

	for i, cid := range v.Competition {
		s.comptitions[i] = domain.CompetitionIndex{
			Id:    cid,
			Phase: domain.CompetitionPhasePreliminary,
		}
	}

	s.aiQuestion = v.AIQuestionInfo

	return s
}

func (s *challengeService) Apply(cmd *CompetitorApplyCmd) error {
	c := cmd.toCompetitor()

	for i := range s.comptitions {
		err := s.competitionRepo.SaveCompetitor(&s.comptitions[i], c)
		if err != nil && !repository.IsErrorDuplicateCreating(err) {
			return err
		}
	}

	err := s.aiQuestionRepo.SaveCompetitor(s.aiQuestion.AIQuestionId, c)
	if err != nil && !repository.IsErrorDuplicateCreating(err) {
		return err
	}

	ver := app.GetCurrentCourseAgree()
	user, err := s.userRepo.GetByAccount(cmd.Account)
	if err != nil {
		return err
	}

	if user.CourseAgreement != ver {
		user.CourseAgreement = ver
		logrus.Debugf("User %s challenge agreement updated from %s to %s ",
			user.Account.Account(), user.CourseAgreement, ver)
		if _, err = s.userRepo.Save(&user); err != nil {
			return err
		}
	}

	return nil
}

func (s *challengeService) GetCompetitor(user domain.Account) (
	ChallengeCompetitorInfoDTO, error,
) {
	dto := ChallengeCompetitorInfoDTO{}
	if user == nil {
		return dto, nil
	}

	score := 0
	for i := range s.comptitions {
		isCompetitor, v, err := s.getCompetitorOfCompetition(
			&s.comptitions[i], user,
		)

		if err != nil || !isCompetitor {
			return dto, err
		}

		score += v
	}

	isCompetitor, v, submission, err := s.aiQuestionRepo.GetCompetitorAndSubmission(
		s.aiQuestion.AIQuestionId, user,
	)
	if err == nil && isCompetitor {
		dto.IsCompetitor = true
		dto.Score = score + v

		dto.AIQuestionInfo.RemainingTimes = s.aiQuestion.RetryTimes - submission.Times
		if submission.Id != "" {
			dto.AIQuestionInfo.InProgress = submission.Status == domain.AIQuestionStatusStart
		}
	}

	return dto, err
}

func (s *challengeService) getCompetitorOfCompetition(
	index *domain.CompetitionIndex, user domain.Account,
) (isCompetitor bool, score int, err error) {

	isCompetitor, submissions, err := s.competitionRepo.GetCompetitorAndSubmission(
		index, user,
	)
	if err != nil || !isCompetitor {
		return
	}

	score = s.helper.CalcCompetitionScore(submissions)

	return
}

func (s *challengeService) SubmitAIQuestionAnswer(competitor domain.Account, cmd *AIQuestionAnswerSubmitCmd) (
	score int, code string, err error,
) {
	now := utils.Now()

	v, err := s.aiQuestionRepo.GetSubmission(
		s.aiQuestion.AIQuestionId, competitor, utils.ToDate(now),
	)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			err = errors.New("illegal submission")
		}
		return
	}

	if v.Status != domain.AIQuestionStatusStart {
		err = errors.New("can't submit")

		return
	}

	if now > v.Expiry {
		code = ErrorCodeAIQuestionSubmissionExpiry
		err = errors.New("it is timeout")

		return
	}

	if cmd.Times != v.Times {
		code = ErrorCodeAIQuestionSubmissionUnmatchedTimes
		err = errors.New("unmatched times")

		return
	}

	answer, err := s.decryptAnswer(cmd.Answer)
	if err != nil {
		return
	}

	if len(cmd.Result) != len(answer) {
		err = errors.New("invalid result")

		return
	}

	score = s.helper.CalcAIQuestionScore(cmd.Result, answer)
	if score > v.Score {
		v.Score = score
	}

	v.Status = domain.AIQuestionStatusEnd

	_, err = s.aiQuestionRepo.SaveSubmission(
		s.aiQuestion.AIQuestionId, &v,
	)

	return
}

func (s *challengeService) GetAIQuestions(competitor domain.Account) (
	dto AIQuestionDTO, code string, err error,
) {
	now := utils.Now()
	date := utils.ToDate(now)
	expiry := now + int64((s.aiQuestion.Timeout+10)*60)

	v, err := s.aiQuestionRepo.GetSubmission(
		s.aiQuestion.AIQuestionId, competitor, date,
	)
	if err != nil {
		if !repository.IsErrorResourceNotExists(err) {
			return
		}

		// gen question first to avoid occupying a times.
		if err = s.genAIQuestions(&dto); err != nil {
			return
		}

		v = domain.QuestionSubmission{
			Account: competitor,
			Date:    date,
			Status:  domain.AIQuestionStatusStart,
			Expiry:  expiry,
			Times:   1,
		}

		_, err = s.aiQuestionRepo.SaveSubmission(s.aiQuestion.AIQuestionId, &v)

		dto.Times = v.Times

		return
	}

	if v.Times >= s.aiQuestion.RetryTimes {
		code = ErrorCodeAIQuestionExceedMaxTimes
		err = errors.New("exceed max times")

		return
	}

	// gen question first to avoid occupying a times.
	if err = s.genAIQuestions(&dto); err != nil {
		return
	}

	v.Status = domain.AIQuestionStatusStart
	v.Expiry = expiry
	v.Times++

	_, err = s.aiQuestionRepo.SaveSubmission(s.aiQuestion.AIQuestionId, &v)

	dto.Times = v.Times

	return
}

func (s *challengeService) genAIQuestions(dto *AIQuestionDTO) (err error) {
	choice, completion := s.helper.GenAIQuestionNums()
	choices, completions, err := s.aiQuestionRepo.GetQuestions(
		s.aiQuestion.QuestionPoolId, choice, completion,
	)
	if err != nil {
		return
	}

	n := len(choice)
	answers := make([]string, n+len(completion))
	dto.Choices = make([]ChoiceQuestionDTO, n)

	for i := range choices {
		item := &choices[i]

		dto.Choices[i] = ChoiceQuestionDTO{
			Desc:    item.Desc,
			Options: item.Options,
		}

		answers[i] = item.Answer
	}

	dto.Completions = make([]CompletionQuestionDTO, len(completion))

	for i := range completions {
		item := &completions[i]

		dto.Completions[i] = CompletionQuestionDTO{
			Desc: item.Desc,
			Info: item.Info,
		}

		answers[i+n] = item.Answer
	}

	str, err := s.encryptAnswer(answers)
	if err == nil {
		dto.Answer = str
	}

	return
}

func (s *challengeService) GetRankingList() ([]ChallengeRankingDTO, error) {
	r := map[string]int{}

	for i := range s.comptitions {
		_, _, results, err := s.competitionRepo.GetResult(&s.comptitions[i])
		if err != nil {
			return nil, err
		}

		scores := s.helper.CalcCompetitionScoreForAll(results)

		for k, v := range scores {
			if _, ok := r[k]; !ok {
				r[k] = v
			}
		}
	}

	scores, err := s.getResultOfAIQuestion()
	if err != nil {
		return nil, err
	}

	for k, v := range scores {
		r[k] += v
	}

	dto := make([]ChallengeRankingDTO, len(r))
	i := 0
	for k, v := range r {
		dto[i] = ChallengeRankingDTO{
			Competitor: k,
			Score:      v,
		}

		i++
	}

	sort.Slice(dto, func(i, j int) bool {
		return dto[i].Score >= dto[j].Score
	})

	return dto, nil
}

func (s *challengeService) getResultOfAIQuestion() (map[string]int, error) {
	results, err := s.aiQuestionRepo.GetResult(s.aiQuestion.AIQuestionId)
	if err != nil {
		return nil, err
	}

	r := map[string]int{}

	for i := range results {
		name := results[i].Account.Account()
		score := results[i].Score

		if r[name] < score {
			r[name] = score
		}
	}

	return r, nil
}

func (s *challengeService) encryptAnswer(answers []string) (string, error) {
	str := strings.Join(answers, s.delimiter)

	v, err := s.encryption.Encrypt([]byte(str))
	if err == nil {
		return base64.StdEncoding.EncodeToString(v), nil
	}

	return "", err
}

func (s *challengeService) decryptAnswer(str string) ([]string, error) {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	v, err := s.encryption.Decrypt(b)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(v), s.delimiter), nil
}

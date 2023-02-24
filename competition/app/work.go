package app

import (
	"errors"
	"sort"

	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

func (s *competitionService) GetRankingList(cid string) (
	dto CompetitonRankingDTO, err error,
) {
	order, err := s.repo.FindScoreOrder(cid)
	if err != nil {
		return
	}

	results, err := s.workRepo.FindWorks(cid)
	if err != nil || len(results) == 0 {
		return
	}

	dto.Final = s.getRankingList(
		results, domain.CompetitionPhaseFinal, order,
	)

	dto.Preliminary = s.getRankingList(
		results, domain.CompetitionPhasePreliminary, order,
	)

	return
}

func (s *competitionService) getRankingList(
	ws []domain.Work,
	phase domain.CompetitionPhase,
	order domain.CompetitionScoreOrder,
) []RankingDTO {
	dtos := make([]RankingDTO, 0, len(ws))
	for i := range ws {
		if v := ws[i].BestOne(phase, order); v != nil {
			dtos = append(dtos, RankingDTO{
				Score:    v.Score,
				TeamName: ws[i].PlayerName,
				SubmitAt: utils.ToDate(v.SubmitAt),
			})
		}
	}

	// sort
	sort.Slice(dtos, func(i, j int) bool {
		return order.IsBetterThanB(dtos[i].Score, dtos[j].Score)
	})

	return dtos
}

func (s *competitionService) GetSubmissions(cid string, user types.Account) (
	dto CompetitionSubmissionsDTO, err error,
) {
	competition, err := s.repo.FindCompetition(cid)
	if err != nil {
		return
	}

	p, _, err := s.playerRepo.FindPlayer(cid, user)
	if err != nil {
		return
	}

	w, _, err := s.workRepo.FindWork(
		domain.NewWorkIndex(cid, p.Id), competition.Phase,
	)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	dto.RelatedProject = w.Repo

	results := w.Submissions(competition.Phase)
	if len(results) == 0 {
		return
	}

	v := make([]*domain.Submission, len(results))
	for i := range results {
		v[i] = &results[i]
	}

	sort.Slice(v, func(i, j int) bool {
		return v[i].SubmitAt >= v[j].SubmitAt
	})

	items := make([]CompetitionSubmissionDTO, len(v))
	for i := range v {
		s.toCompetitionSubmissionDTO(v[i], &items[i])
	}

	dto.Details = items

	return
}

func (s *competitionService) AddRelatedProject(cmd *CompetitionAddReleatedProjectCMD) (
	code string, err error,
) {
	competition, err := s.repo.FindCompetition(cmd.Id)
	if err != nil {
		return
	}

	if !competition.IsPreliminary() {
		err = errors.New("it can only change the related project on preliminary phase")

		return
	}

	p, _, err := s.playerRepo.FindPlayer(cmd.Id, cmd.User)
	if err != nil {
		return
	}

	if !p.IsIndividualOrLeader() {
		code = errorNoPermission
		err = errors.New("no permission to submit")

		return
	}

	w, version, err := s.workRepo.FindWork(
		domain.NewWorkIndex(cmd.Id, p.Id), competition.Phase,
	)
	if err != nil {
		if !repoerr.IsErrorResourceNotExists(err) {
			return
		}

		w = domain.NewWork(cmd.Id, &p)
		if err = s.workRepo.SaveWork(&w); err != nil {
			return
		}
	}

	w.Repo = cmd.repo()
	err = s.workRepo.SaveRepo(&w, version)

	return
}

func (s *competitionService) Submit(cmd *CompetitionSubmitCMD) (
	dto CompetitionSubmissionDTO, code string, err error,
) {
	competition, err := s.repo.FindCompetition(cmd.CompetitionId)
	if err != nil {
		return
	}

	if competition.IsOver() {
		err = errors.New("competition is over")

		return
	}

	p, _, err := s.playerRepo.FindPlayer(cmd.CompetitionId, cmd.User)
	if err != nil {
		return
	}

	if !p.IsIndividualOrLeader() {
		code = errorNoPermission
		err = errors.New("no permission to submit")

		return
	}

	if competition.IsFinal() && !p.IsFinalist {
		code = errorNotFinalist
		err = errors.New("you are not finalist")

		return
	}

	// work
	phase := competition.Phase
	w, version, err := s.workRepo.FindWork(
		domain.NewWorkIndex(competition.Id, p.Id), phase,
	)
	if err != nil {
		if !repoerr.IsErrorResourceNotExists(err) {
			return
		}

		w = domain.NewWork(competition.Id, &p)
		if err = s.workRepo.SaveWork(&w); err != nil {
			return
		}
	}

	if w.HasSubmittedToday(phase) {
		code = errorSubmitTooMany
		err = errors.New("submit more than one time per day")

		return
	}

	// submit
	ps, err := s.submissionServie.Submit(
		&w, phase, cmd.FileName, cmd.Data,
	)
	if err != nil {
		return
	}

	if err = s.workRepo.AddSubmission(&w, &ps, version); err != nil {
		return
	}

	// notify
	info := w.NewSubmissionMessage(&ps)
	if err = s.producer.NotifyCalcScore(&info); err != nil {
		return
	}

	dto.FileName = cmd.FileName
	dto.SubmitAt = utils.ToDate(ps.SubmitAt)
	dto.Status = ps.Status

	return
}

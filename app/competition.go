package app

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/competition"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type CompetitionIndex = domain.CompetitionIndex
type CompetitionListCMD = repository.CompetitionListOption
type CompetitionSubmissionInfo = domain.CompetitionSubmissionInfo
type CompetitionService interface {
	Apply(string, *CompetitorApplyCmd) error

	Get(cid string, competitor domain.Account) (UserCompetitionDTO, error)
	List(*CompetitionListCMD) ([]CompetitionSummaryDTO, error)

	Submit(*CompetitionSubmitCMD) (CompetitionSubmissionDTO, string, error)

	GetSubmissions(*CompetitionIndex, domain.Account) (CompetitionSubmissionsDTO, error)

	GetTeam(cid string, competitor domain.Account) (CompetitionTeamDTO, error)

	GetRankingList(cid string, phase domain.CompetitionPhase) ([]RankingDTO, error)

	AddRelatedProject(*CompetitionAddReleatedProjectCMD) error
}

func NewCompetitionService(
	repo repository.Competition,
	sender message.Sender,
	uploader competition.Competition,
) CompetitionService {
	return competitionService{
		repo:     repo,
		sender:   sender,
		uploader: uploader,
	}
}

type competitionService struct {
	repo     repository.Competition
	sender   message.Sender
	uploader competition.Competition
}

func (s competitionService) Apply(cid string, cmd *CompetitorApplyCmd) error {
	err := s.repo.SaveCompetitor(
		&domain.CompetitionIndex{
			Id:    cid,
			Phase: domain.CompetitionPhasePreliminary,
		},
		cmd.toCompetitor(),
	)
	if err != nil && !repository.IsErrorDuplicateCreating(err) {
		return err
	}

	return nil
}

func (s competitionService) Get(cid string, competitor domain.Account) (
	dto UserCompetitionDTO, err error,
) {
	index := domain.CompetitionIndex{
		Id:    cid,
		Phase: domain.CompetitionPhasePreliminary,
	}

	v, b, err := s.repo.Get(&index, competitor)
	if err != nil {
		return
	}

	s.toCompetitionDTO(&v.Competition, &dto.CompetitionDTO)

	dto.CompetitorCount = v.CompetitorCount

	dto.IsCompetitor = b.IsCompetitor
	if !b.IsCompetitor {
		dto.DatasetURL = ""
	}

	dto.TeamId = b.TeamId
	if b.TeamRole != nil {
		dto.TeamRole = b.TeamRole.TeamRole()
	}

	// Only the normal competition can change the phase
	if v.Type.CompetitionType() == "" && !v.Enabled {
		dto.Phase = domain.CompetitionPhaseFinal.CompetitionPhase()
	}

	return
}

func (s competitionService) List(cmd *CompetitionListCMD) (
	dtos []CompetitionSummaryDTO, err error,
) {
	cmd.Phase = domain.CompetitionPhasePreliminary

	v, err := s.repo.List(cmd)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]CompetitionSummaryDTO, len(v))

	for i := range v {
		s.toCompetitionSummaryDTO(&v[i].CompetitionSummary, &dtos[i])

		dtos[i].CompetitorCount = v[i].CompetitorCount
	}

	return
}

func (s competitionService) GetTeam(cid string, competitor domain.Account) (
	dto CompetitionTeamDTO, err error,
) {
	index := domain.CompetitionIndex{
		Id:    cid,
		Phase: domain.CompetitionPhasePreliminary,
	}

	v, err := s.repo.GetTeam(&index, competitor)
	if err != nil {
		return
	}

	if name := v[0].Team.Name; name != nil {
		dto.Name = name.TeamName()
	}

	members := make([]CompetitionTeamMemberDTO, len(v))
	for i := range v {
		item := &v[i]

		members[i] = CompetitionTeamMemberDTO{
			Name:  item.Name.CompetitorName(),
			Email: item.Email.Email(),
		}

		if item.TeamRole != nil {
			members[i].Role = item.TeamRole.TeamRole()
		}
	}

	dto.Members = members

	return
}

func (s competitionService) GetRankingList(cid string, phase domain.CompetitionPhase) (
	dtos []RankingDTO, err error,
) {
	index := domain.CompetitionIndex{
		Id:    cid,
		Phase: phase,
	}

	order, teams, results, err := s.repo.GetResult(&index)
	if err != nil || len(results) == 0 {
		return
	}

	rs := map[string]*domain.CompetitionSubmission{}

	for i := range results {
		item := &results[i]

		if !item.IsSuccess() {
			continue
		}

		k := item.Key()
		if v, ok := rs[k]; !ok || order.IsBetterThanB(item.Score, v.Score) {
			rs[k] = item
		}
	}

	// sort
	i := 0
	rl := make([]*domain.CompetitionSubmission, len(rs))
	for _, v := range rs {
		rl[i] = v
		i++
	}

	sort.Slice(rl, func(i, j int) bool {
		return order.IsBetterThanB(rl[i].Score, rl[j].Score)
	})

	// result
	tm := map[string]string{}
	for i := range teams {
		tm[teams[i].Id] = teams[i].Name.TeamName()
	}

	dtos = make([]RankingDTO, len(rl))
	for i := range rl {
		item := rl[i]

		dtos[i] = RankingDTO{
			Score:    item.Score,
			SubmitAt: utils.ToDate(item.SubmitAt),
		}

		if item.IsTeamWork() {
			dtos[i].TeamName = tm[item.TeamId]
		} else {
			// If it is individual, just show its account instead of name.
			// Because the name maybe duplicate, but the account will not.
			dtos[i].TeamName = item.Individual.Account()
		}
	}

	return
}

func (s competitionService) GetSubmissions(index *CompetitionIndex, competitor domain.Account) (
	dto CompetitionSubmissionsDTO, err error,
) {
	repo, results, err := s.repo.GetSubmisstions(index, competitor)
	if err != nil {
		return
	}

	if repo.Owner != nil {
		dto.RelatedProject = repo.Owner.Account() + "/" + repo.Repo.ResourceName()
	}

	if len(results) == 0 {
		return
	}

	v := make([]*domain.CompetitionSubmission, len(results))
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

func (s competitionService) AddRelatedProject(cmd *CompetitionAddReleatedProjectCMD) (
	err error,
) {
	if cmd.Index.Phase.IsFinal() {
		err = errors.New("can't change the related project on final phase")

		return
	}

	// check permission
	v, b, err := s.repo.Get(&cmd.Index, cmd.Competitor)
	if err != nil {
		return
	}

	if !b.IsCompetitor || (b.TeamId != "" && !b.TeamRole.IsLeader()) {
		err = errors.New("no permission to submit")

		return
	}

	if !v.Enabled {
		err = errors.New("competition is over for this phase")

		return
	}

	project := &cmd.Project

	if cmd.Competitor.Account() != project.Owner.Account() {
		err = errors.New("can't add project which is not your's")

		return
	}

	repo := domain.CompetitionRepo{
		Owner: project.Owner,
		Repo:  project.Name,
	}
	if b.TeamId == "" {
		repo.Individual = cmd.Competitor
	} else {
		repo.TeamId = b.TeamId
	}

	return s.repo.AddRelatedProject(&cmd.Index, &repo)
}

func (s competitionService) Submit(cmd *CompetitionSubmitCMD) (
	dto CompetitionSubmissionDTO, code string, err error,
) {
	index := &cmd.Index

	// check permission
	v, b, err := s.repo.Get(index, cmd.Competitor)
	if err != nil {
		return
	}

	if !b.IsCompetitor || (b.TeamId != "" && !b.TeamRole.IsLeader()) {
		err = errors.New("no permission to submit")

		return
	}

	if !v.Enabled {
		err = errors.New("competition is over for this phase")

		return
	}

	// upload file
	user := b.TeamId
	if b.TeamId == "" {
		user = cmd.Competitor.Account()
	}
	now := utils.Now()
	obspath := fmt.Sprintf(
		"%s/%s/%s/%s_%s",
		index.Id, index.Phase.CompetitionPhase(),
		user, strconv.FormatInt(now, 10), cmd.FileName,
	)
	if err = s.uploader.UploadSubmissionFile(cmd.Data, obspath); err != nil {
		return
	}

	// save
	submission := domain.CompetitionSubmission{
		SubmitAt: now,
		OBSPath:  obspath,
		Status:   "calculating",
	}
	if b.TeamId == "" {
		submission.Individual = cmd.Competitor
	} else {
		submission.TeamId = b.TeamId
	}

	sid, err := s.repo.SaveSubmission(index, &submission)
	if err != nil {
		if repository.IsErrorDuplicateCreating(err) {
			code = ErrorCompetitionDuplicateSubmission
		}

		return
	}

	// send mq
	info := message.SubmissionInfo{
		Id:      sid,
		Index:   *index,
		OBSPath: obspath,
	}

	if err = s.sender.CalcScore(&info); err != nil {
		return
	}

	dto.FileName = cmd.FileName
	dto.SubmitAt = utils.ToDate(now)
	dto.Status = submission.Status

	return
}

// Internal Service
type CompetitionInternalService interface {
	UpdateSubmission(*CompetitionIndex, *CompetitionSubmissionInfo) error
}

func NewCompetitionInternalService(repo repository.Competition) CompetitionInternalService {
	return competitionInternalService{
		repo: repo,
	}
}

type competitionInternalService struct {
	repo repository.Competition
}

func (s competitionInternalService) UpdateSubmission(
	index *CompetitionIndex, info *CompetitionSubmissionInfo,
) error {
	return s.repo.UpdateSubmission(index, info)
}

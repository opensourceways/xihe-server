package app

import (
	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	"github.com/opensourceways/xihe-server/competition/domain/uploader"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

type CompetitionService interface {
	// player
	Apply(string, *CompetitorApplyCmd) (string, error)
	CreateTeam(cid string, cmd *CompetitionTeamCreateCmd) (string, error)
	JoinTeam(cid string, cmd *CompetitionTeamJoinCmd) (string, error)
	GetMyTeam(cid string, competitor types.Account) (CompetitionTeamDTO, string, error)
	ChangeTeamName(cid string, cmd *CmdToChangeCompetitionTeamName) error
	TransferLeader(cid string, cmd *CmdToTransferTeamLeader) error
	QuitTeam(cid string, competitor types.Account) error
	DeleteMember(cid string, cmd *CmdToDeleteTeamMember) error

	// competition
	Get(cid string, competitor types.Account) (UserCompetitionDTO, error)
	List(*CompetitionListCMD) ([]CompetitionSummaryDTO, error)

	// work
	Submit(*CompetitionSubmitCMD) (CompetitionSubmissionDTO, string, error)
	GetSubmissions(string, types.Account) (CompetitionSubmissionsDTO, error)
	GetRankingList(string) (CompetitonRankingDTO, error)
	AddRelatedProject(*CompetitionAddReleatedProjectCMD) (string, error)
}

var _ CompetitionService = (*competitionService)(nil)

func NewCompetitionService(
	repo repository.Competition,
	workRepo repository.Work,
	playerRepo repository.Player,
	producer message.CalcScoreMessageProducer,
	uploader uploader.SubmissionFileUploader,
) *competitionService {
	return &competitionService{
		repo:             repo,
		workRepo:         workRepo,
		playerRepo:       playerRepo,
		producer:         producer,
		submissionServie: domain.NewSubmissionService(uploader),
	}
}

type competitionService struct {
	repo             repository.Competition
	workRepo         repository.Work
	playerRepo       repository.Player
	producer         message.CalcScoreMessageProducer
	submissionServie domain.SubmissionService
}

// show competition detail
func (s *competitionService) Get(cid string, user types.Account) (
	dto UserCompetitionDTO, err error,
) {
	c, err := s.repo.FindCompetition(cid)
	if err != nil {
		return
	}
	s.toCompetitionDTO(&c, &dto.CompetitionDTO)

	// competitors count
	dto.CompetitorCount, err = s.playerRepo.CompetitorsCount(cid)
	if err != nil {
		return
	}

	// competitor info
	if user == nil {
		return
	}

	p, _, err := s.playerRepo.FindPlayer(cid, user)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			err = nil
			dto.DatasetURL = ""
		}

		return
	}

	dto.IsFinalist = p.IsFinalist
	dto.IsCompetitor = true
	if p.IsATeam() {
		dto.TeamId = p.Id
		dto.TeamRole = p.RoleOfCurrentCompetitor()
	}

	return
}

func (s *competitionService) List(cmd *CompetitionListCMD) (
	dtos []CompetitionSummaryDTO, err error,
) {
	if cmd.User != nil {
		return s.getCompetitionsUserApplied(cmd)
	}

	return s.listCompetitions(&repository.CompetitionListOption{
		Status: cmd.Status,
	})
}

func (s *competitionService) listCompetitions(opt *repository.CompetitionListOption) (
	dtos []CompetitionSummaryDTO, err error,
) {
	v, err := s.repo.FindCompetitions(opt)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]CompetitionSummaryDTO, len(v))
	for i := range v {
		dtos[i].CompetitorCount, err = s.playerRepo.CompetitorsCount(v[i].Id)
		if err != nil {
			return
		}

		s.toCompetitionSummaryDTO(&v[i], &dtos[i])
	}

	return
}

func (s *competitionService) getCompetitionsUserApplied(cmd *CompetitionListCMD) (
	[]CompetitionSummaryDTO, error,
) {
	cs, err := s.playerRepo.FindCompetitionsUserApplied(cmd.User)
	if err != nil {
		return nil, err
	}

	return s.listCompetitions(&repository.CompetitionListOption{
		Status:         cmd.Status,
		CompetitionIds: cs,
	})
}

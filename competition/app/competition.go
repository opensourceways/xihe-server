package app

import (
	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	"github.com/opensourceways/xihe-server/competition/domain/uploader"
	"github.com/opensourceways/xihe-server/competition/domain/user"
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
	DissolveTeam(cid string, leader types.Account) error

	// competition
	Get(*CompetitionGetCmd) (UserCompetitionDTO, error)
	List(*CompetitionListCMD) ([]CompetitionSummaryDTO, error)

	// work
	Submit(*CompetitionSubmitCMD) (CompetitionSubmissionDTO, string, error)
	GetSubmissions(*CompetitionGetCmd) (CompetitionSubmissionsDTO, error)
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
	userCli user.User,
) *competitionService {
	return &competitionService{
		repo:             repo,
		workRepo:         workRepo,
		playerRepo:       playerRepo,
		producer:         producer,
		submissionServie: domain.NewSubmissionService(uploader),
		userCli:          userCli,
	}
}

type competitionService struct {
	repo             repository.Competition
	workRepo         repository.Work
	playerRepo       repository.Player
	producer         message.CalcScoreMessageProducer
	submissionServie domain.SubmissionService
	userCli          user.User
}

// show competition detail
func (s *competitionService) Get(cmd *CompetitionGetCmd) (
	dto UserCompetitionDTO, err error,
) {
	c, err := s.repo.FindCompetition(&repository.CompetitionGetOption{
		CompetitionId: cmd.CompetitionId,
		Lang:          cmd.Lang,
	})
	if err != nil {
		return
	}

	// competitors count
	n, err := s.playerRepo.CompetitorsCount(cmd.CompetitionId)
	if err != nil {
		return
	}

	s.toCompetitionDTO(&c, n, &dto.CompetitionDTO)

	// competitor info
	if cmd.User == nil {
		return
	}

	p, _, err := s.playerRepo.FindPlayer(cmd.CompetitionId, cmd.User)
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
		Tag:    cmd.Tag,
		Lang:   cmd.Lang,
	})
}

func (s *competitionService) listCompetitions(opt *repository.CompetitionListOption) (
	[]CompetitionSummaryDTO, error,
) {
	v, err := s.repo.FindCompetitions(opt)
	if err != nil || len(v) == 0 {
		return nil, err
	}

	dtos := make([]CompetitionSummaryDTO, len(v))
	for i := range v {
		n, err := s.playerRepo.CompetitorsCount(v[i].Id)
		if err != nil {
			return nil, err
		}

		s.toCompetitionSummaryDTO(&v[i], n, &dtos[i])
	}

	return dtos, nil
}

func (s *competitionService) getCompetitionsUserApplied(cmd *CompetitionListCMD) (
	[]CompetitionSummaryDTO, error,
) {
	cs, err := s.playerRepo.FindCompetitionsUserApplied(cmd.User)
	if err != nil || len(cs) == 0 {
		return nil, err
	}

	return s.listCompetitions(&repository.CompetitionListOption{
		Status:         cmd.Status,
		CompetitionIds: cs,
	})
}

package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

func (s *competitionService) Apply(cid string, cmd *CompetitorApplyCmd) (code string, err error) {
	competition, err := s.repo.FindCompetition(&repository.CompetitionGetOption{
		CompetitionId: cid,
	})
	if err != nil {
		return
	}

	if competition.IsOver() {
		err = errors.New("competition is over")

		return
	}

	if competition.IsFinal() {
		err = errors.New("apply on final phase")

		return
	}

	p := cmd.toPlayer(cid)
	if err = s.playerRepo.AddPlayer(&p); err != nil {
		if repoerr.IsErrorDuplicateCreating(err) {
			code = errorCompetitorExists
		}
	}

	if err = s.userCli.AddUserRegInfo(&p.Leader); err != nil {
		return
	}

	msg := new(message.CompetitionMsg)
	msg.GenApplyMsg(cmd.Account, competition.CompetitionSummary.Name.CompetitionName())
	s.producer.SendCompetitionMsg(msg)

	return
}

func (s *competitionService) CreateTeam(cid string, cmd *CompetitionTeamCreateCmd) (
	code string, err error,
) {
	p, version, err := s.playerRepo.FindPlayer(cid, cmd.User)
	if err != nil {
		return
	}

	if err = p.CreateTeam(cmd.Name); err != nil {
		return
	}

	if err = s.playerRepo.DeletePlayer(&p, version); err != nil {
		return
	}

	if err = s.playerRepo.AddPlayer(&p); err != nil {
		if repoerr.IsErrorDuplicateCreating(err) {
			code = errorTeamExists
		}

		utils.RetryThreeTimes(func() error {
			return s.playerRepo.ResumePlayer(cid, p.Leader.Account)
		})
	}

	return
}

func (s *competitionService) JoinTeam(cid string, cmd *CompetitionTeamJoinCmd) (code string, err error) {
	me, pv, err := s.playerRepo.FindPlayer(cid, cmd.User)
	if err != nil {
		return
	}

	team, version, err := s.playerRepo.FindPlayer(cid, cmd.Leader)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			code = errorNoCorrespondingTeam
		}

		return
	}

	if err = me.JoinTo(&team); err != nil {
		if domain.IsErrorTeamMembersEnough(err) {
			code = errorTeamMembersEnough
		}

		return
	}

	if err = s.playerRepo.DeletePlayer(&me, pv); err != nil {
		return
	}

	err = s.playerRepo.AddMember(
		repository.PlayerVersion{
			Player:  &team,
			Version: version,
		},
		repository.PlayerVersion{
			Player:  &me,
			Version: pv,
		},
	)
	if err != nil {
		utils.RetryThreeTimes(func() error {
			return s.playerRepo.ResumePlayer(cid, me.Leader.Account)
		})
	}

	return
}

func (s *competitionService) GetMyTeam(cid string, user types.Account) (
	dto CompetitionTeamDTO, code string, err error,
) {
	p, _, err := s.playerRepo.FindPlayer(cid, user)
	if err != nil {
		return
	}

	if !p.IsATeam() {
		code = errorNotATeam
		err = errors.New("not a team")

		return
	}

	dto.Name = p.Name()

	m := p.Members()
	members := make([]CompetitionTeamMemberDTO, p.CompetitorsCount())
	for i := range m {
		item := &m[i]
		members[i+1] = CompetitionTeamMemberDTO{
			Name:    item.Name.CompetitorName(),
			Email:   item.Email.Email(),
			Account: item.Account.Account(),
		}
	}

	leader := &p.Leader
	members[0] = CompetitionTeamMemberDTO{
		Name:    leader.Name.CompetitorName(),
		Email:   leader.Email.Email(),
		Role:    domain.TeamLeaderRole(),
		Account: leader.Account.Account(),
	}

	dto.Members = members

	return
}

func (s *competitionService) ChangeTeamName(cid string, cmd *CmdToChangeCompetitionTeamName) error {
	p, version, err := s.playerRepo.FindPlayer(cid, cmd.User)
	if err != nil {
		return err
	}

	if err = p.ChangeTeamName(cmd.Name); err != nil {
		return err
	}

	return s.playerRepo.SaveTeamName(&p, version)
}

func (s *competitionService) TransferLeader(cid string, cmd *CmdToTransferTeamLeader) error {
	p, version, err := s.playerRepo.FindPlayer(cid, cmd.Leader)
	if err != nil {
		return err
	}

	if err = p.TransferLeader(cmd.User); err != nil {
		return err
	}

	return s.playerRepo.SavePlayer(&p, version)
}

func (s *competitionService) QuitTeam(cid string, competitor types.Account) error {
	p, version, err := s.playerRepo.FindPlayer(cid, competitor)
	if err != nil {
		return err
	}

	if err = p.Quit(); err != nil {
		return err
	}

	if err = s.playerRepo.SavePlayer(&p, version); err != nil {
		return err
	}

	return s.playerRepo.ResumePlayer(cid, competitor)
}

func (s *competitionService) DeleteMember(cid string, cmd *CmdToDeleteTeamMember) error {
	p, version, err := s.playerRepo.FindPlayer(cid, cmd.Leader)
	if err != nil {
		return err
	}

	if err = p.Delete(cmd.User); err != nil {
		return err
	}

	if err = s.playerRepo.ResumePlayer(cid, cmd.User); err != nil {
		return err
	}

	return s.playerRepo.SavePlayer(&p, version)
}

func (s *competitionService) DissolveTeam(cid string, leader types.Account) error {
	p, version, err := s.playerRepo.FindPlayer(cid, leader)
	if err != nil {
		return err
	}

	for _, m := range p.Members() {
		if err = p.Delete(m.Account); err != nil {
			return err
		}

		if err = s.playerRepo.ResumePlayer(cid, m.Account); err != nil {
			return err
		}

		if err = s.playerRepo.SavePlayer(&p, version); err != nil {
			return err
		}

		version++
	}

	if err = s.playerRepo.ResumePlayer(cid, leader); err != nil {
		return err
	}

	return s.playerRepo.DeletePlayer(&p, version)
}

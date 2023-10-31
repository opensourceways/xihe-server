package controller

import (
	"github.com/opensourceways/xihe-server/competition/app"
	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type AddRelatedProjectRequest struct {
	Owner string `json:"owner"`
	Name  string `json:"project_name"`
}

func (req *AddRelatedProjectRequest) ToInfo() (
	owner types.Account, name types.ResourceName, err error,
) {
	if owner, err = types.NewAccount(req.Owner); err != nil {
		return
	}

	name, err = types.NewResourceName(req.Name)

	return
}

type CompetitorApplyRequest struct {
	Name      string            `json:"name"`
	City      string            `json:"city"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Identity  string            `json:"identity"`
	Province  string            `json:"province"`
	Detail    map[string]string `json:"detail"`
	Agreement bool              `json:"agreement"`
}

func (req *CompetitorApplyRequest) ToCmd(user types.Account) (cmd app.CompetitorApplyCmd, err error) {
	if cmd.Name, err = domain.NewCompetitorName(req.Name); err != nil {
		return
	}

	if cmd.City, err = domain.NewCity(req.City); err != nil {
		return
	}

	if cmd.Email, err = types.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Phone, err = domain.NewPhone(req.Phone); err != nil {
		return
	}

	if cmd.Identity, err = domain.NewcompetitionIdentity(req.Identity); err != nil {
		return
	}

	if cmd.Province, err = domain.NewProvince(req.Province); err != nil {
		return
	}

	cmd.Detail = req.Detail
	cmd.Account = user

	err = cmd.Validate()

	return
}

type CreateTeamRequest struct {
	Name string `json:"team_name"`
}

func (req *CreateTeamRequest) ToCmd(user types.Account) (
	cmd app.CompetitionTeamCreateCmd, err error,
) {
	if cmd.Name, err = domain.NewTeamName(req.Name); err != nil {
		return
	}

	cmd.User = user

	return
}

type JoinTeamRequest struct {
	Account string `json:"leader_account"`
}

func (req *JoinTeamRequest) ToCmd(user types.Account) (
	cmd app.CompetitionTeamJoinCmd, err error,
) {
	if cmd.Leader, err = types.NewAccount(req.Account); err != nil {
		return
	}

	cmd.User = user

	return
}

type ChangeTeamNameRequest = CreateTeamRequest

type TransferLeaderRequest struct {
	Account string `json:"competitor_account"`
}

func (req *TransferLeaderRequest) ToCmd(leader types.Account) (
	cmd app.CmdToTransferTeamLeader, err error,
) {
	if cmd.User, err = types.NewAccount(req.Account); err != nil {
		return
	}

	cmd.Leader = leader

	return

}

type DeleteMemberRequest = TransferLeaderRequest

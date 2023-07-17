package app

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type CompetitionListCMD struct {
	Status domain.CompetitionStatus
	User   types.Account
	Tag    domain.CompetitionTag
}

type CompetitorApplyCmd domain.Competitor

func (cmd *CompetitorApplyCmd) Validate() error {
	b := cmd.Account != nil &&
		cmd.Name != nil &&
		cmd.Email != nil &&
		cmd.Identity != nil

	if !b {
		return errors.New("invalid cmd")
	}

	for i := range cmd.Detail {
		if utils.StrLen(cmd.Detail[i]) > 20 {
			return errors.New("invalid detail")
		}
	}

	return nil
}

func (cmd *CompetitorApplyCmd) toPlayer(cid string) (p domain.Player) {
	p.CompetitionId = cid
	p.Leader = *(*domain.Competitor)(cmd)

	return
}

type CompetitionSubmitCMD struct {
	CompetitionId string
	FileName      string
	Data          io.Reader
	User          types.Account
}

func (cmd *CompetitionSubmitCMD) Validate() error {
	if len(cmd.FileName) > 0 && !utils.IsSafeFileName(cmd.FileName) {
		return errors.New("invalid file name")
	}

	return nil
}

type CompetitionAddReleatedProjectCMD struct {
	Id      string
	User    types.Account
	Project types.ResourceSummary
}

func (cmd *CompetitionAddReleatedProjectCMD) repo() string {
	return cmd.User.Account() + "/" + cmd.Project.Name.ResourceName()
}

type CompetitionSummaryDTO struct {
	CompetitorCount int      `json:"count"`
	Bonus           int      `json:"bonus"`
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Host            string   `json:"host"`
	Desc            string   `json:"desc"`
	Status          string   `json:"status"`
	Poster          string   `json:"poster"`
	Duration        string   `json:"duration"`
	Tags            []string `json:"tags"`
}

type CompetitionDTO struct {
	CompetitionSummaryDTO

	Type       string `json:"type"`
	Phase      string `json:"phase"`
	Doc        string `json:"doc"`
	Forum      string `json:"forum"`
	Winners    string `json:"winners"`
	DatasetDoc string `json:"dataset_doc"`
	DatasetURL string `json:"dataset_url"`
}

type UserCompetitionDTO struct {
	TeamId       string `json:"team_id"`
	TeamRole     string `json:"team_role"`
	IsFinalist   bool   `json:"is_finalist"`
	IsCompetitor bool   `json:"is_competitor"`

	CompetitionDTO
}

// ranking
type CompetitonRankingDTO struct {
	Final       []RankingDTO `json:"final"`
	Preliminary []RankingDTO `json:"preliminary"`
}

type RankingDTO struct {
	Score    float32 `json:"score"`
	TeamName string  `json:"team_name"`
	SubmitAt string  `json:"submit_at"`
}

// team
type CompetitionTeamCreateCmd struct {
	User types.Account
	Name domain.TeamName
}

type CompetitionTeamJoinCmd struct {
	User   types.Account
	Leader types.Account
}

type CompetitionTeamDTO struct {
	Name    string                     `json:"name"`
	Members []CompetitionTeamMemberDTO `json:"members"`
}

type CompetitionTeamMemberDTO struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Email   string `json:"email"`
	Account string `json:"account"`
}

// result
type CompetitionSubmissionsDTO struct {
	RelatedProject string                     `json:"project"`
	Details        []CompetitionSubmissionDTO `json:"details"`
}

type CompetitionSubmissionDTO struct {
	SubmitAt string  `json:"submit_at"`
	FileName string  `json:"file_name"`
	Status   string  `json:"status"`
	Score    float32 `json:"score"`
}

func (s competitionService) toCompetitionSubmissionDTO(
	v *domain.Submission, dto *CompetitionSubmissionDTO,
) {
	*dto = CompetitionSubmissionDTO{
		SubmitAt: utils.ToDate(v.SubmitAt),
		FileName: filepath.Base(v.OBSPath),
		Status:   v.Status,
		Score:    v.Score,
	}
}

func (s competitionService) toCompetitionSummaryDTO(
	c *domain.CompetitionSummary, competitorsCount int,
	dto *CompetitionSummaryDTO,
) {
	*dto = CompetitionSummaryDTO{
		Bonus:           c.Bonus.CompetitionBonus(),
		Id:              c.Id,
		Name:            c.Name.CompetitionName(),
		Host:            c.Host.CompetitionHost(),
		Desc:            c.Desc.CompetitionDesc(),
		Status:          c.Status.CompetitionStatus(),
		Poster:          c.Poster.URL(),
		Duration:        c.Duration.CompetitionDuration(),
		CompetitorCount: competitorsCount,
	}

	for _, v := range c.Tags {
		dto.Tags = append(dto.Tags, v.CompetitionTag())
	}
}

func (s competitionService) toCompetitionDTO(
	c *domain.Competition, competitorsCount int, dto *CompetitionDTO,
) {
	s.toCompetitionSummaryDTO(
		&c.CompetitionSummary, competitorsCount,
		&dto.CompetitionSummaryDTO,
	)

	dto.Type = c.Type.CompetitionType()
	dto.Phase = c.Phase.CompetitionPhase()
	dto.Doc = c.Doc.URL()
	dto.Forum = c.Forum.Forum()
	dto.Winners = c.Winners.Winners()
	dto.DatasetDoc = c.DatasetDoc.URL()
	dto.DatasetURL = c.DatasetURL.URL()
}

type CmdToChangeCompetitionTeamName = CompetitionTeamCreateCmd

type CmdToTransferTeamLeader = CompetitionTeamJoinCmd

type CmdToDeleteTeamMember = CompetitionTeamJoinCmd

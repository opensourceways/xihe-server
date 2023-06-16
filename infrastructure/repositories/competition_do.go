package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type CompetitionListOptionDO struct {
	Phase      string
	Status     string
	Competitor string
}

func (impl competition) toCompetitionListOptionDO(
	opt *repository.CompetitionListOption, do *CompetitionListOptionDO,
) {
	if opt.Phase != nil {
		do.Phase = opt.Phase.CompetitionPhase()
	}

	if opt.Status != nil {
		do.Status = opt.Status.CompetitionStatus()
	}

	if opt.Competitor != nil {
		do.Competitor = opt.Competitor.Account()
	}
}

type CompetitionSummaryDO struct {
	Bonus    int
	Id       string
	Name     string
	Desc     string
	Host     string
	Status   string
	Poster   string
	Duration string

	CompetitorsCount int
}

func (do *CompetitionSummaryDO) toCompetitionSummary(
	c *domain.CompetitionSummary,
) (err error) {
	c.Id = do.Id

	if c.Bonus, err = domain.NewCompetitionBonus(do.Bonus); err != nil {
		return
	}

	if c.Name, err = domain.NewCompetitionName(do.Name); err != nil {
		return
	}

	if c.Desc, err = domain.NewCompetitionDesc(do.Desc); err != nil {
		return
	}

	if c.Host, err = domain.NewCompetitionHost(do.Host); err != nil {
		return
	}

	if c.Status, err = domain.NewCompetitionStatus(do.Status); err != nil {
		return
	}

	if c.Poster, err = domain.NewURL(do.Poster); err != nil {
		return
	}

	c.Duration, err = domain.NewCompetitionDuration(do.Duration)

	return
}

type CompetitionDO struct {
	CompetitionSummaryDO

	Enabled    bool
	Type       string
	Phase      string
	Doc        string
	Forum      string
	Winners    string
	DatasetDoc string
	DatasetURL string
}

func (do *CompetitionDO) toCompetition(
	c *domain.Competition,
) (err error) {
	c.Enabled = do.Enabled

	if c.Type, err = domain.NewCompetitionType(do.Type); err != nil {
		return
	}

	if c.Phase, err = domain.NewCompetitionPhase(do.Phase); err != nil {
		return
	}

	if c.Doc, err = domain.NewURL(do.Doc); err != nil {
		return
	}

	if c.Forum, err = domain.NewForum(do.Forum); err != nil {
		return
	}

	if c.Winners, err = domain.NewWinners(do.Winners); err != nil {
		return
	}

	if c.DatasetDoc, err = domain.NewURL(do.DatasetDoc); err != nil {
		return
	}

	if c.DatasetURL, err = domain.NewURL(do.DatasetURL); err != nil {
		return
	}

	if err = do.CompetitionSummaryDO.toCompetitionSummary(&c.CompetitionSummary); err != nil {
		return
	}

	return
}

type CompetitorInfoDO struct {
	Account  string
	Name     string
	City     string
	Email    string
	Phone    string
	Identity string
	Province string
	Detail   map[string]string
}

func toCompetitorInfoDO(c *domain.CompetitorInfo, do *CompetitorInfoDO) {
	*do = CompetitorInfoDO{
		Account:  c.Account.Account(),
		Name:     c.Name.CompetitorName(),
		Email:    c.Email.Email(),
		Identity: c.Identity.CompetitionIdentity(),
		Detail:   c.Detail,
	}

	if c.City != nil {
		do.City = c.City.City()
	}

	if c.Phone != nil {
		do.Phone = c.Phone.Phone()
	}

	if c.Province != nil {
		do.Province = c.Province.Province()
	}
}

func (do *CompetitorInfoDO) toCompetitorInfo(c *domain.CompetitorInfo) (err error) {
	if c.Account, err = domain.NewAccount(do.Account); err != nil {
		return
	}

	if c.Name, err = domain.NewCompetitorName(do.Name); err != nil {
		return
	}

	if c.City, err = domain.NewCity(do.City); err != nil {
		return
	}

	if c.Email, err = domain.NewEmail(do.Email); err != nil {
		return
	}

	if c.Phone, err = domain.NewPhone(do.Phone); err != nil {
		return
	}

	if c.Identity, err = domain.NewcompetitionIdentity(do.Identity); err != nil {
		return
	}

	if c.Province, err = domain.NewProvince(do.Province); err != nil {
		return
	}

	c.Detail = do.Detail

	return
}

type CompetitorDO struct {
	CompetitorInfoDO

	TeamId   string
	TeamRole string
	TeamName string
}

func (do *CompetitorDO) toCompetitor(c *domain.Competitor) (err error) {
	if err = do.CompetitorInfoDO.toCompetitorInfo(&c.CompetitorInfo); err != nil {
		return
	}

	if do.TeamId == "" {
		return
	}

	c.Team.Id = do.TeamId

	if c.Team.Name, err = domain.NewTeamName(do.TeamName); err != nil {
		return
	}

	if c.TeamRole, err = domain.NewTeamRole(do.TeamRole); err != nil {
		return
	}

	return
}

type CompetitionSubmissionDO struct {
	Id         string
	TeamId     string
	Individual string

	SubmitAt int64
	OBSPath  string
	Status   string
	Score    float32
}

func (do *CompetitionSubmissionDO) toCompetitionSubmission(r *domain.CompetitionSubmission) (err error) {
	*r = domain.CompetitionSubmission{
		Id:       do.Id,
		SubmitAt: do.SubmitAt,
		OBSPath:  do.OBSPath,
		Status:   do.Status,
		Score:    do.Score,
		TeamId:   do.TeamId,
	}

	if do.Individual != "" {
		r.Individual, err = domain.NewAccount(do.Individual)
	}

	return
}

func (impl competition) toCompetitionSubmissionDO(
	r *domain.CompetitionSubmission,
	do *CompetitionSubmissionDO,
) {
	*do = CompetitionSubmissionDO{
		TeamId:   r.TeamId,
		SubmitAt: r.SubmitAt,
		OBSPath:  r.OBSPath,
		Status:   r.Status,
		Score:    r.Score,
	}

	if r.Individual != nil {
		do.Individual = r.Individual.Account()
	}
}

type CompetitionTeamDO struct {
	Id   string
	Name string
}

func (do *CompetitionTeamDO) toTeam(r *domain.CompetitionTeam) (err error) {
	r.Id = do.Id

	if do.Name != "" {
		r.Name, err = domain.NewTeamName(do.Name)
	}

	return
}

type CompetitionIndexDO struct {
	Id    string
	Phase string
}

func (impl competition) toCompetitionIndexDO(index *domain.CompetitionIndex) CompetitionIndexDO {
	return CompetitionIndexDO{
		Id:    index.Id,
		Phase: index.Phase.CompetitionPhase(),
	}
}

func (impl competition) toCompetitionRepoDO(
	repo *domain.CompetitionRepo, do *CompetitionRepoDO,
) {
	*do = CompetitionRepoDO{
		TeamId: repo.TeamId,
		Owner:  repo.Owner.Account(),
		Repo:   repo.Repo.ResourceName(),
	}

	if repo.Individual != nil {
		do.Individual = repo.Individual.Account()
	}
}

type CompetitionRepoDO struct {
	TeamId     string
	Individual string

	Owner string
	Repo  string
}

func (do *CompetitionRepoDO) toCompetitionRepo(r *domain.CompetitionRepo) (err error) {
	r.TeamId = do.TeamId

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Repo, err = domain.NewResourceName(do.Repo); err != nil {
		return
	}

	if do.Individual != "" {
		r.Individual, err = domain.NewAccount(do.Individual)
	}

	return
}

type CompetitorSummaryDO struct {
	IsCompetitor bool
	TeamId       string
	TeamRole     string
}

func (do *CompetitorSummaryDO) toCompetitorSummary(r *domain.CompetitorSummary) (err error) {
	r.IsCompetitor = do.IsCompetitor
	r.TeamId = do.TeamId
	r.TeamRole, err = domain.NewTeamRole(do.TeamRole)

	return
}

type CompetitionSubmissionInfoDO = domain.CompetitionSubmissionInfo

package repositoryimpl

import (
	"errors"

	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

func (doc *dCompetition) toCompetitionSummary(c *domain.CompetitionSummary) (err error) {
	c.Id = doc.Id

	if c.Bonus, err = domain.NewCompetitionBonus(doc.Bonus); err != nil {
		return
	}

	if c.Name, err = domain.NewCompetitionName(doc.Name); err != nil {
		return
	}

	if c.Desc, err = domain.NewCompetitionDesc(doc.Desc); err != nil {
		return
	}

	if c.Host, err = domain.NewCompetitionHost(doc.Host); err != nil {
		return
	}

	if c.Status, err = domain.NewCompetitionStatus(doc.Status); err != nil {
		return
	}

	if c.Poster, err = domain.NewURL(doc.Poster); err != nil {
		return
	}

	if c.Lang, err = domain.NewLanguage(doc.Language); err != nil {
		return
	}

	var t domain.CompetitionTag
	for _, v := range doc.Tags {
		if t, err = domain.NewCompetitionTag(v); err != nil {
			return
		}
		c.Tags = append(c.Tags, t)
	}

	c.Duration, err = domain.NewCompetitionDuration(doc.Duration)

	return
}

func (doc *dCompetition) toCompetition(c *domain.Competition) (err error) {
	if c.Type, err = domain.NewCompetitionType(doc.Type); err != nil {
		return
	}

	if c.Phase, err = domain.NewCompetitionPhase(doc.Phase); err != nil {
		return
	}

	c.Order = domain.NewCompetitionScoreOrder(doc.SmallerOk)

	if c.Doc, err = domain.NewURL(doc.Doc); err != nil {
		return
	}

	if c.Forum, err = domain.NewForum(doc.Forum); err != nil {
		return
	}

	if c.Winners, err = domain.NewWinners(doc.Winners); err != nil {
		return
	}

	if c.DatasetDoc, err = domain.NewURL(doc.DatasetDoc); err != nil {
		return
	}

	if c.DatasetURL, err = domain.NewURL(doc.DatasetURL); err != nil {
		return
	}

	err = doc.toCompetitionSummary(&c.CompetitionSummary)

	return
}

func (doc *dWork) toWork(w *domain.Work) {
	w.CompetitionId = doc.CompetitionId
	w.PlayerName = doc.PlayerName
	w.PlayerId = doc.PlayerId
	w.Repo = doc.Repo

	if r := doc.toSubmissions(doc.Preliminary); len(r) != 0 {
		w.Preliminary = r
	}

	if r := doc.toSubmissions(doc.Final); len(r) != 0 {
		w.Final = r
	}
}

func (doc *dWork) toSubmissions(v []dSubmission) []domain.Submission {
	if len(v) == 0 {
		return nil
	}

	r := make([]domain.Submission, len(v))
	for i := range v {
		v[i].toSubmission(&r[i])
	}

	return r
}

func (doc *dSubmission) toSubmission(s *domain.Submission) {
	*s = domain.Submission{
		Id:       doc.Id,
		Status:   doc.Status,
		OBSPath:  doc.OBSPath,
		SubmitAt: doc.SubmitAt,
		Score:    float32(doc.Score),
	}
}

func (doc *dPlayer) toPlayer(p *domain.Player) error {
	cs, err := doc.toCompetitors()
	if err != nil {
		return err
	}

	p.Leader = cs[0]

	if doc.TeamName != "" {
		if p.Team.Name, err = domain.NewTeamName(doc.TeamName); err != nil {
			return err
		}

		if len(cs) > 1 {
			p.Team.Members = cs[1:]
		}
	}

	p.Id = doc.Id.Hex()
	p.IsFinalist = doc.IsFinalist
	p.CompetitionId = doc.CompetitionId

	return nil
}

func (doc *dPlayer) toCompetitors() ([]domain.Competitor, error) {
	if len(doc.Competitors) == 0 {
		return nil, errors.New("imporsible, no competitors")
	}

	lp := -1
	for i := range doc.Competitors {
		if doc.Leader == doc.Competitors[i].Account {
			lp = i
		}
	}
	if lp < 0 {
		return nil, errors.New("imporsible, no leader")
	}

	r := make([]domain.Competitor, len(doc.Competitors))
	j := 1
	for i := range doc.Competitors {
		k := 0
		if i != lp {
			k = j
			j++
		}

		if err := doc.Competitors[i].toCompetitor(&r[k]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (doc *dCompetitor) toCompetitor(c *domain.Competitor) (err error) {
	if c.Account, err = types.NewAccount(doc.Account); err != nil {
		return
	}

	if c.Name, err = domain.NewCompetitorName(doc.Name); err != nil {
		return
	}

	if c.City, err = domain.NewCity(doc.City); err != nil {
		return
	}

	if c.Email, err = types.NewEmail(doc.Email); err != nil {
		return
	}

	if c.Phone, err = domain.NewPhone(doc.Phone); err != nil {
		return
	}

	if c.Identity, err = domain.NewcompetitionIdentity(doc.Identity); err != nil {
		return
	}

	if c.Province, err = domain.NewProvince(doc.Province); err != nil {
		return
	}

	c.Detail = doc.Detail

	return
}

func toCompetitorDoc(c *domain.Competitor) dCompetitor {
	doc := dCompetitor{
		Account:  c.Account.Account(),
		Name:     c.Name.CompetitorName(),
		Email:    c.Email.Email(),
		Identity: c.Identity.CompetitionIdentity(),
		Detail:   c.Detail,
	}

	if c.City != nil {
		doc.City = c.City.City()
	}

	if c.Phone != nil {
		doc.Phone = c.Phone.Phone()
	}

	if c.Province != nil {
		doc.Province = c.Province.Province()
	}

	return doc
}

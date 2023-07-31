package domain

type CompetitionSummary struct {
	Id       string
	Name     CompetitionName
	Desc     CompetitionDesc
	Host     CompetitionHost
	Bonus    CompetitionBonus
	Status   CompetitionStatus
	Duration CompetitionDuration
	Poster   URL
	Tags     []CompetitionTag
	Lang     Language
}

type Competition struct {
	CompetitionSummary

	Doc        URL
	Forum      Forum
	Winners    Winners
	DatasetDoc URL
	DatasetURL URL

	Type  CompetitionType
	Phase CompetitionPhase
	Order CompetitionScoreOrder
}

func (c *Competition) IsOver() bool {
	return c.Status != nil && c.Status.IsOver()
}

func (c *Competition) IsPreliminary() bool {
	return c.Phase.IsPreliminary()
}

func (c *Competition) IsFinal() bool {
	return c.Phase.IsFinal()
}

// CompetitionScoreOrder
type CompetitionScoreOrder interface {
	IsBetterThanB(a, b float32) bool
}

func NewCompetitionScoreOrder(b bool) CompetitionScoreOrder {
	return smallerIsBetter(b)
}

type smallerIsBetter bool

func (order smallerIsBetter) IsBetterThanB(a, b float32) bool {
	if order {
		return a <= b
	}

	return a >= b
}

package domain

type Task struct {
	Id    string
	Names Sentence
	Rule  Rule
}

type Rule struct {
	Descs     Sentence
	CreatedAt string
	MaxPoints int
}

func (r *Rule) IsValidPoint(point int) bool {
	return point <= r.MaxPoints
}

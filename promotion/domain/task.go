package domain

type Task struct {
	Id    string
	Names Sentence
	Rule  Rule
}

type Rule struct {
	Descs     Sentence
	CreatedAt string
	Points    int
}

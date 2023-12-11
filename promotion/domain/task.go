package domain

import (
	"errors"

	"github.com/opensourceways/xihe-server/utils"
)

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

func (r *Task) ToItem(point int) (Item, error) {
	if point > r.Rule.MaxPoints {
		return Item{}, errors.New("point over allowed")
	}

	return Item{
		TaskId:   r.Id,
		TaskName: r.Names,
		Descs:    r.Rule.Descs,
		Date:     utils.Now(),
		Points:   point,
	}, nil
}

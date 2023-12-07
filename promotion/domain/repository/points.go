package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type Points interface {
	Save(*Point) error
	Update(user string, done string) error
	Find(types.Account) (Point, error)
	FindAll() ([]Point, error)
}

type Point struct {
	User    types.Account
	Dones   []Done
	Version int
}

type Done struct {
	TaskId string
	Date   string
}

func (r Point) HasDone(taskid string) bool {
	for i := range r.Dones {
		if r.Dones[i].TaskId == taskid {
			return true
		}
	}

	return false
}

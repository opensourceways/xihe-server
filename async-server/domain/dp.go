package domain

import "errors"

const (
	taskStatusWaiting  = "waiting"
	taskStatusFinished = "finished"
	taskStatusError    = "error"
)

type TaskStatus interface {
	TaskStatus() string
	IsWaiting() bool
	IsFinished() bool
	IsError() bool
}

func NewTaskStatus(v string) (TaskStatus, error) {
	b := v == taskStatusWaiting ||
		v == taskStatusFinished ||
		v == taskStatusError

	if !b {
		return nil, errors.New("invalid value")
	}

	return dptaskstatus(v), nil
}

type dptaskstatus string

func (r dptaskstatus) TaskStatus() string {
	return string(r)
}

func (r dptaskstatus) IsWaiting() bool {
	return r.TaskStatus() == taskStatusWaiting
}

func (r dptaskstatus) IsFinished() bool {
	return r.TaskStatus() == taskStatusFinished
}

func (r dptaskstatus) IsError() bool {
	return r.TaskStatus() == taskStatusError
}

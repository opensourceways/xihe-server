package domain

import (
	"errors"
	"strings"
)

const (
	taskStatusWaiting  = "waiting"
	taskStatusRunning  = "running"
	taskStatusFinished = "finished"
	taskStatusError    = "error"

	taskTypeWuKong     = "wukong"
	taskTypeWuKong4Img = "wukong_4img"
)

// taskStatus
type TaskStatus interface {
	TaskStatus() string
	IsWaiting() bool
	IsRunning() bool
	IsFinished() bool
	IsError() bool
}

func NewTaskStatus(v string) (TaskStatus, error) {
	b := v == taskStatusWaiting ||
		v == taskStatusRunning ||
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

func (r dptaskstatus) IsRunning() bool {
	return r.TaskStatus() == taskStatusRunning
}

func (r dptaskstatus) IsFinished() bool {
	return r.TaskStatus() == taskStatusFinished
}

func (r dptaskstatus) IsError() bool {
	return r.TaskStatus() == taskStatusError
}

// Task Type
type TaskType interface {
	TaskType() string
	IsWuKong() bool
	IsWuKong4Img() bool
}

type dptasktype string

func NewTaskType(v string) (TaskType, error) {
	b := v == taskTypeWuKong ||
		v == taskTypeWuKong4Img

	if !b {
		return nil, errors.New("invalid value")
	}

	return dptasktype(v), nil
}

func (r dptasktype) TaskType() string {
	return string(r)
}

func (r dptasktype) IsWuKong() bool {
	return r.TaskType() == taskTypeWuKong
}

func (r dptasktype) IsWuKong4Img() bool {
	return r.TaskType() == taskTypeWuKong4Img
}

// Links
type Links interface {
	Links() []string
	StringLinks() string
}

func NewLinks(v string) (Links, error) {
	return dplinks(strings.Split(v, ",")), nil
}

func NewLinksFromMap(v map[string]string) (Links, error) {
	if len(v) == 0 {
		return nil, errors.New("invalid value")
	}

	a := make([]string, len(v))
	var i int
	for _, val := range v {
		a[i] = val
		i++
	}

	return dplinks(a), nil
}

type dplinks []string

func (r dplinks) Links() []string {
	return ([]string)(r)
}

func (r dplinks) StringLinks() string {
	s := ""

	for _, v := range r.Links() {
		s += v + ","
	}

	return strings.TrimRight(s, ",")
}

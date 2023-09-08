package domain

import (
	common "github.com/opensourceways/xihe-server/domain"
)

// UserPoints
type UserPoints struct {
	User    common.Account
	Total   int
	Date    string
	Items   []PointsItem // items of corresponding date
	Dones   []string     // tasks that user has done
	Version int
}

func (entity *UserPoints) AddPointsItem(task *Task, time, desc string) *PointsItem {
	newItem := false

	item := entity.poitsItem(task.Name)
	if item == nil {
		item = &PointsItem{Task: task.Name}
		newItem = true
	}

	v := entity.calc(task, item)
	if v == 0 {
		return nil
	}

	entity.Total += v

	item.add(&PointsDetail{
		Id:     "", // TODO uuid
		Time:   time,
		Desc:   desc,
		Points: v,
	})

	if newItem {
		entity.Items = append(entity.Items, *item)

		item = &entity.Items[len(entity.Items)-1]
	}

	if !entity.hasDone(task.Name) {
		entity.Dones = append(entity.Dones, task.Name)
	}

	return item
}

func (entity *UserPoints) IsCompleted(task *Task) bool {
	item := entity.poitsItem(task.Name)
	if item == nil {
		return false
	}

	v := task.Rule.calc(item.points(), !entity.hasDone(task.Name))

	return v == 0
}

func (entity *UserPoints) calc(task *Task, item *PointsItem) int {
	pointsOfDay := entity.pointsOfDay()

	if pointsOfDay >= config.MaxPointsOfDay {
		return 0
	}

	v := task.Rule.calc(item.points(), !entity.hasDone(task.Name))
	if v == 0 {
		return 0
	}

	if n := config.MaxPointsOfDay - pointsOfDay; v >= n {
		return n
	}

	return v
}

func (entity *UserPoints) hasDone(t string) bool {
	for _, i := range entity.Dones {
		if i == t {
			return true
		}
	}

	return false
}

func (entity *UserPoints) pointsOfDay() int {
	n := 0
	for i := range entity.Items {
		n += entity.Items[i].points()
	}

	return n
}

func (entity *UserPoints) poitsItem(t string) *PointsItem {
	items := entity.Items

	for i := range items {
		if items[i].Task == t {
			return &items[i]
		}
	}

	return nil
}

// PointsItem
type PointsItem struct {
	Task    string
	Details []PointsDetail
}

func (item *PointsItem) points() int {
	if item == nil {
		return 0
	}

	n := 0
	for i := range item.Details {
		n += item.Details[i].Points
	}

	return n
}

func (item *PointsItem) add(p *PointsDetail) {
	item.Details = append(item.Details, *p)
}

// PointsDetail
type PointsDetail struct {
	Id     string `json:"id"`
	Time   string `json:"time"`
	Desc   string `json:"desc"`
	Points int    `json:"points"`
}

// Task
type Task struct {
	Name string
	Kind string //Novice, EveryDay, Activity
	Rule Rule
}

// Rule
type Rule struct {
	OnceOnly       bool // only can do once
	Desc           string
	CreatedAt      string
	PointsPerOnce  int
	MaxPointsOfDay int
}

// points is the one that user has got on this task today
// firstTime is that user is first time to do this task
func (r *Rule) calc(points int, firstTime bool) int {
	if r.OnceOnly {
		if firstTime {
			return 0
		}

		return r.PointsPerOnce
	}

	if r.MaxPointsOfDay > 0 && points >= r.MaxPointsOfDay {
		return 0
	}

	return r.PointsPerOnce
}

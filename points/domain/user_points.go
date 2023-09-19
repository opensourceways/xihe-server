package domain

import (
	"strconv"

	"github.com/sirupsen/logrus"

	common "github.com/opensourceways/xihe-server/common/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

// UserPoints
type UserPoints struct {
	User    types.Account
	Total   int
	Items   []PointsItem // items of day or all the items
	Dones   []string     // tasks that user has done
	Version int
}

func (entity *UserPoints) DetailsNum() int {
	n := 0
	for i := range entity.Items {
		n += entity.Items[i].detailsNum()
	}

	return n
}

func (entity *UserPoints) IsFirstPointsDetailOfDay() bool {
	return len(entity.Items) == 1 && entity.Items[0].isFirstDetail()
}

func (entity *UserPoints) AddPointsItem(task *Task, date string, detail *PointsDetail) *PointsItem {
	item := entity.pointsItem(task.Id)

	v := entity.calc(task, item)
	if v == 0 {
		return nil
	}

	entity.Total += v

	detail.Id = date + "_" + strconv.Itoa(entity.DetailsNum()+1)
	detail.Points = v

	if !entity.hasDone(task.Id) {
		entity.Dones = append(entity.Dones, task.Id)
	}

	if item != nil {
		item.add(detail)

		return item
	}

	entity.Items = append(entity.Items, PointsItem{
		TaskId:  task.Id,
		Date:    date,
		Details: []PointsDetail{*detail},
	})

	return &entity.Items[len(entity.Items)-1]
}

func (entity *UserPoints) IsCompleted(task *Task) bool {
	item := entity.pointsItem(task.Id)

	v := task.Rule.calcPoints(item.points(), entity.hasDone(task.Id))

	return v == 0
}

func (entity *UserPoints) calc(task *Task, item *PointsItem) int {
	pointsOfDay := entity.pointsOfDay()

	if pointsOfDay >= config.MaxPointsOfDay {
		logrus.Warnf("Points of day reached the limit: %d, points: %d", config.MaxPointsOfDay, pointsOfDay)

		return 0
	}

	v := task.Rule.calcPoints(item.points(), entity.hasDone(task.Id))
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

func (entity *UserPoints) pointsItem(t string) *PointsItem {
	items := entity.Items

	for i := range items {
		if items[i].TaskId == t {
			return &items[i]
		}
	}

	return nil
}

// PointsItem
type PointsItem struct {
	Date    string
	TaskId  string
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

func (item *PointsItem) detailsNum() int {
	return len(item.Details)
}

func (item *PointsItem) isFirstDetail() bool {
	return item != nil && len(item.Details) == 1
}

func (item *PointsItem) LatestDetail() *PointsDetail {
	if item == nil || len(item.Details) == 0 {
		return nil
	}

	return &item.Details[len(item.Details)-1]
}

// PointsDetail
type PointsDetail struct {
	Id      string `json:"id"` // serial number
	Desc    string `json:"desc"`
	TimeStr string `json:"time_str"`
	Time    int64  `json:"time"`
	Points  int    `json:"points"`
}

// Task
type Task struct {
	Id    string            `json:"id"`
	Names map[string]string `json:"names"`
	Kind  string            `json:"kind"` // Novice, EveryDay, Activity, PassiveItem
	Addr  string            `json:"addr"` // The website address of task
	Rule  Rule              `json:"rule"`
}

func (t *Task) Name(lang common.Language) string {
	return t.Names[lang.Language()]
}

func (t *Task) IsPassiveTask() bool {
	return t.Kind == "PassiveItem"
}

// Rule
type Rule struct {
	Descs          map[string]string `json:"descs"`
	CreatedAt      string            `json:"created_at"`
	OnceOnly       bool              `json:"once_only"` // only can do once
	PointsPerOnce  int               `json:"points_per_once"`
	MaxPointsOfDay int               `json:"max_points_of_day"`
}

// points is the one that user has got on this task today
func (r *Rule) calcPoints(points int, hasDone bool) int {
	if r.OnceOnly {
		if hasDone {
			logrus.Warn("Rule has been done today, will not calc again.")

			return 0
		}

		return r.PointsPerOnce
	}

	if r.MaxPointsOfDay > 0 && points >= r.MaxPointsOfDay {
		logrus.Warnf("Points of today: %d, exceed the limit: %d", points, r.MaxPointsOfDay)

		return 0
	}

	return r.PointsPerOnce
}

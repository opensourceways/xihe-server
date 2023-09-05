package domain

import (
	common "github.com/opensourceways/xihe-server/domain"
)

type PointsCalculator interface {
	// pointsOfDay is the total points that user has got that day
	// pointsOfItem is the points that user has got on the item that day
	Calc(pointsOfDay, pointsOfItem int) int
}

// UserPoints
type UserPoints struct {
	User    common.Account
	Total   int
	Date    string
	Items   []PointsItem // items of corresponding date
	Version int
}

func (entity *UserPoints) AddPointsItem(t string, detail *PointsDetail, r PointsCalculator) *PointsItem {
	item := entity.poitsItem(t)

	v := r.Calc(entity.pointsOfDay(), item.points())
	if v == 0 {
		return nil
	}

	entity.Total += v
	detail.Points = v

	item.add(detail)

	return item
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
		if items[i].Type == t {
			return &items[i]
		}
	}

	entity.Items = append(items, PointsItem{Type: t})

	return &entity.Items[len(entity.Items)-1]
}

// PointsItem
type PointsItem struct {
	Type    string
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

// PointsItemRule
type PointsItemRule struct {
	Type           string
	Desc           string
	CreatedAt      string
	PointsPerOnce  int
	MaxPointsOfDay int
}

// points is the one that user has got on this item
func (r *PointsItemRule) Calc(points int) int {
	if r.MaxPointsOfDay > 0 && points >= r.MaxPointsOfDay {
		return 0
	}

	return r.PointsPerOnce
}

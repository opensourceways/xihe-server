package domain

import (
	common "github.com/opensourceways/xihe-server/domain"
)

type PointsCalculator interface {
	Calc(total, pointsOfDay, pointsOfItem int) int
}

// UserPoint
type UserPoint struct {
	User    common.Account
	Total   int
	Date    string
	Items   []PointItem // items of corresponding date
	Version int
}

func (entity *UserPoint) AddPointItem(t string, detail *PointDetail, r PointsCalculator) *PointItem {
	item := entity.poitItem(t)

	v := r.Calc(entity.Total, entity.pointsOfDay(), item.points())
	if v == 0 {
		return nil
	}

	entity.Total += v
	detail.Point = v

	item.add(detail)

	return item
}

func (entity *UserPoint) pointsOfDay() int {
	n := 0
	for i := range entity.Items {
		n += entity.Items[i].points()
	}

	return n
}

func (entity *UserPoint) poitItem(t string) *PointItem {
	items := entity.Items

	for i := range items {
		if items[i].Type == t {
			return &items[i]
		}
	}

	entity.Items = append(items, PointItem{Type: t})

	return &entity.Items[len(entity.Items)-1]
}

// PointItem
type PointItem struct {
	Type    string
	Details []PointDetail
}

func (item *PointItem) points() int {
	if item == nil {
		return 0
	}

	n := 0
	for i := range item.Details {
		n += item.Details[i].Point
	}

	return n
}

func (item *PointItem) add(p *PointDetail) {
	item.Details = append(item.Details, *p)
}

// PointDetail
type PointDetail struct {
	Id    string `json:"id"`
	Time  string `json:"time"`
	Desc  string `json:"desc"`
	Point int    `json:"point"`
}

// PointItemRule
type PointItemRule struct {
	Type           string
	Desc           string
	CreatedAt      string
	PointsPerOnce  int
	MaxPointsOfDay int
}

// points is the one that user has got on this item
func (r *PointItemRule) Calc(points int) int {
	if points >= r.MaxPointsOfDay {
		return 0
	}

	return r.PointsPerOnce
}

/*
type PointItemRules []PointRule

// times is the count which the action happened and it starts from 1.
func (val PointItemRules) Calc(t string, times int) int {
	for i := range val {
		if val[i].Type == t {
			return val[i].calc(times)
		}
	}

	return 0
}
*/

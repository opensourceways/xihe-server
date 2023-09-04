package domain

import (
	common "github.com/opensourceways/xihe-server/domain"
)

// UserPoint
type UserPoint struct {
	User    common.Account
	Total   int
	Date    string
	Items   []PointItem // items of corresponding date
	Version int
}

func (entity *UserPoint) AddPointItem(t string, detail *PointDetail, r *PointRule) *PointItem {
	item := entity.poitItem(t)

	v := r.calc(item.count() + 1)
	if v == 0 {
		return nil
	}

	entity.Total += v
	detail.Point = v

	item.add(detail)

	return item
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

func (item *PointItem) count() int {
	if item == nil {
		return 0
	}

	return len(item.Details)
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

// PointRule
type PointRule struct {
	Type           string
	Desc           string
	CreatedAt      string
	Point          int
	MaxTimesPerDay int
}

// times is the count which the action happened and it starts from 1.
func (r *PointRule) calc(times int) int {
	if times > r.MaxTimesPerDay {
		return 0
	}

	return r.Point
}

/*
type PointRules []PointRule

// times is the count which the action happened and it starts from 1.
func (val PointRules) Calc(t string, times int) int {
	for i := range val {
		if val[i].Type == t {
			return val[i].calc(times)
		}
	}

	return 0
}
*/

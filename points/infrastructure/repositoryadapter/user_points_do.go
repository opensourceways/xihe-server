package repositoryadapter

import (
	"sort"

	"go.mongodb.org/mongo-driver/bson"

	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/domain"
)

const (
	fieldUser    = "user"
	fieldDate    = "date"
	fieldDays    = "days"
	fieldTotal   = "total"
	fieldDones   = "dones"
	fieldDetails = "details"
	fieldVersion = "version"
)

type userPointsDO struct {
	User    string                 `bson:"user"     json:"user"`
	Days    []pointsDetailsOfDayDO `bson:"days"     json:"days"`
	Dones   []string               `bson:"dones"    json:"dones"`
	Total   int                    `bson:"total"    json:"total"`
	Version int                    `bson:"version"  json:"version"`
}

func (do *userPointsDO) doc() (bson.M, error) {
	return genDoc(do)
}

func (do *userPointsDO) toUserPoints() (domain.UserPoints, error) {
	u, err := common.NewAccount(do.User)
	if err != nil {
		return domain.UserPoints{}, err
	}

	return domain.UserPoints{
		User:    u,
		Total:   do.Total,
		Items:   do.toPointsItems(),
		Dones:   do.Dones,
		Version: do.Version,
	}, nil
}

func (do *userPointsDO) toPointsItems() []domain.PointsItem {
	r := []domain.PointsItem{}

	sort.Slice(do.Days, func(i, j int) bool {
		return do.Days[i].Date < do.Days[j].Date
	})

	for i := len(do.Days) - 1; i >= 0; i-- {
		item := &do.Days[i]

		r = append(r, item.toPointsItems()...)
	}

	return r
}

// pointsDetailsOfDayDO
type pointsDetailsOfDayDO struct {
	Date    string           `bson:"date"     json:"date"`
	Details []pointsDetailDO `bson:"details"  json:"details"`
}

func (do *pointsDetailsOfDayDO) doc() (bson.M, error) {
	return genDoc(do)
}

func (do pointsDetailsOfDayDO) toPointsItems() []domain.PointsItem {
	m := map[string]int{}
	r := []domain.PointsItem{}

	for i := range do.Details {
		item := &do.Details[i]

		j, ok := m[item.TaskId]
		if !ok {
			j = len(r)
			m[item.TaskId] = j

			r = append(r, domain.PointsItem{
				Date:   do.Date,
				TaskId: item.TaskId,
			})
		}

		r[j].Details = append(r[j].Details, item.toPointsDetail())
	}

	return r
}

func toPointsItemsOfDayDO(item *domain.PointsItem) pointsDetailsOfDayDO {
	return pointsDetailsOfDayDO{
		Date: item.Date,
		Details: []pointsDetailDO{
			topointsDetailDO(item.TaskId, item.LatestDetail()),
		},
	}

}

// pointsDetailDO
type pointsDetailDO struct {
	Id     string `bson:"id"       json:"id"`
	Desc   string `bson:"desc"     json:"desc"`
	Time   string `bson:"time"     json:"time"`
	TaskId string `bson:"task_id"  json:"task_id"`
	Points int    `bson:"points"   json:"points"`
}

func (do *pointsDetailDO) toPointsDetail() domain.PointsDetail {
	return domain.PointsDetail{
		Id:     do.Id,
		Desc:   do.Desc,
		Time:   do.Time,
		Points: do.Points,
	}
}

func (do *pointsDetailDO) doc() (bson.M, error) {
	return genDoc(do)
}

func topointsDetailDO(taskId string, detail *domain.PointsDetail) pointsDetailDO {
	return pointsDetailDO{
		Id:     detail.Id,
		Desc:   detail.Desc,
		TaskId: taskId,
		Time:   detail.Time,
		Points: detail.Points,
	}
}

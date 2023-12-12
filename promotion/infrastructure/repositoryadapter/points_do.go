package repositoryadapter

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	fieldUser        = "user"
	fieldItems       = "items"
	fieldTotal       = "total"
	fieldPromotionId = "promotion_id"
)

type pointsDO struct {
	User        string   `bson:"user"           json:"user"`
	PromotionId string   `bson:"promotion_id"   json:"promotion_id"`
	Total       int      `bson:"total"          json:"total"`
	Items       []itemDO `bson:"items"          json:"items"`
	Version     int      `bson:"version"        json:"version"`
}

func (do *pointsDO) doc() (bson.M, error) {
	return genDoc(do)
}

func (do *pointsDO) toUserPoints() (ups domain.UserPoints, err error) {
	if ups.User, err = types.NewAccount(do.User); err != nil {
		return
	}

	ups.Items = make([]domain.Item, len(do.Items))
	for i := range do.Items {
		if ups.Items[i], err = do.Items[i].toItem(); err != nil {
			return
		}
	}

	ups.PromotionId = do.PromotionId
	ups.Total = do.Total
	ups.Version = do.Version

	return
}

func toPointsDO(ups *domain.UserPoints) pointsDO {
	do := pointsDO{
		User:        ups.User.Account(),
		PromotionId: ups.PromotionId,
		Total:       ups.Total,
		Version:     ups.Version,
	}

	items := make([]itemDO, len(ups.Items))
	for i := range items {
		items[i] = toItemDO(&ups.Items[i])
	}

	return do
}

type itemDO struct {
	TaskId   string            `bson:"task_id"   json:"task_id"`
	TaskName map[string]string `bson:"task_name" json:"task_name"`
	Descs    map[string]string `bson:"descs"     json:"descs"`
	Date     int64             `bson:"date"      json:"date"`
	Points   int               `bson:"points"    json:"points"`
}

func (do *itemDO) doc() (bson.M, error) {
	return genDoc(do)
}

func (do *itemDO) toItem() (domain.Item, error) {
	taskname, err := domain.NewSentence(
		do.TaskName[domain.FieldEN], do.TaskName[domain.FieldZH])
	if err != nil {
		return domain.Item{}, err
	}

	descs, err := domain.NewSentence(
		do.Descs[domain.FieldEN], do.Descs[domain.FieldZH])
	if err != nil {
		return domain.Item{}, err
	}

	return domain.Item{
		TaskId:   do.TaskId,
		TaskName: taskname,
		Descs:    descs,
		Date:     do.Date,
		Points:   do.Points,
	}, nil
}

func toItemDO(item *domain.Item) itemDO {
	return itemDO{
		TaskId:   item.TaskId,
		TaskName: item.TaskName.SentenceMap(),
		Descs:    item.Descs.SentenceMap(),
		Date:     item.Date,
		Points:   item.Points,
	}
}

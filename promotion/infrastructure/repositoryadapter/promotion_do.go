package repositoryadapter

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"github.com/opensourceways/xihe-server/promotion/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	fieldRegUsers = "reg_users"
)

type promotionDO struct {
	Id       string      `bson:"id"        json:"id"`
	Name     string      `bson:"name"      json:"name"`
	Desc     string      `bson:"desc"      json:"desc"`
	Poster   string      `bson:"poster"    json:"poster"`
	RegUsers []RegUserDO `bson:"reg_users" json:"reg_users"`
	Duration string      `bson:"duration"  json:"duration"`
	Version  int         `bson:"version"   json:"version"`
}

func (do *promotionDO) toPromotionRepo() (p repository.PromotionRepo, err error) {
	if p.Name, err = domain.NewPromotionName(do.Name); err != nil {
		return
	}

	if p.Desc, err = domain.NewPromotionDesc(do.Desc); err != nil {
		return
	}

	if p.Duration, err = domain.NewPromotionDuration(do.Duration); err != nil {
		return
	}

	p.RegUsers = make([]domain.RegUser, len(do.RegUsers))
	for i := range do.RegUsers {
		if p.RegUsers[i].User, err = types.NewAccount(do.RegUsers[i].User); err != nil {
			return
		}

		if do.RegUsers[i].Origin != "" {
			if p.RegUsers[i].Origin, err = domain.NewOrigin(do.RegUsers[i].Origin); err != nil {
				return
			}
		}

		p.RegUsers[i].CreatedAt = do.RegUsers[i].CreatedAt
	}

	p.Id = do.Id
	p.Poster = do.Poster
	p.Version = do.Version

	return
}

type RegUserDO struct {
	User      string `bson:"user"       json:"user"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	Origin    string `bson:"origin" json:"origin"`
}

func (do *RegUserDO) doc() (bson.M, error) {
	return genDoc(do)
}

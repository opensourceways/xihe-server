package repositoryadapter

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
	"go.mongodb.org/mongo-driver/bson"
)

const fieldRegUsers = "reg_users"

type promotionDO struct {
	Id        string      `bson:"id"         json:"id"`
	Name      string      `bson:"name"       json:"name"`
	Desc      string      `bson:"desc"       json:"desc"`
	Poster    string      `bson:"poster"     json:"poster"`
	RegUsers  []RegUserDO `bson:"reg_users"  json:"reg_users"`
	StartTime int64       `bson:"start_time" json:"start_time"`
	EndTime   int64       `bson:"end_time"   json:"end_time"`
	Version   int         `bson:"version"    json:"version"`
	Way       string      `bson:"way"        json:"way"`
	Host      string      `bson:"host"       json:"host"`
	Type      string      `bson:"type"       json:"type"`
	Intro     string      `bson:"intro"      json:"intro"`
	IsStatic  bool        `bson:"is_static"  json:"is_static"`
	Priority  int         `bson:"priority"   json:"priority"`
}

func (do *promotionDO) toPromotion() (p domain.Promotion, err error) {
	if p.Name, err = domain.NewPromotionName(do.Name); err != nil {
		return
	}

	if p.Desc, err = domain.NewPromotionDesc(do.Desc); err != nil {
		return
	}

	if do.Way != "" {
		if p.Way, err = domain.NewPromotionWay(do.Way); err != nil {
			return
		}
	}

	if do.Type != "" {
		if p.Type, err = domain.NewPromotionType(do.Type); err != nil {
			return
		}
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
	p.StartTime = do.StartTime
	p.EndTime = do.EndTime
	p.Poster = do.Poster
	p.Host = do.Host
	p.Intro = do.Intro
	p.IsStatic = do.IsStatic
	p.Version = do.Version

	return
}

type RegUserDO struct {
	User      string `bson:"user"       json:"user"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	Origin    string `bson:"origin"     json:"origin"`
}

func (do *RegUserDO) doc() (bson.M, error) {
	return genDoc(do)
}

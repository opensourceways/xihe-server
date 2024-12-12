package domain

import (
	"errors"
	"fmt"
	"time"

	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	layout   = "2006.01.02"
	location = "Asia/Shanghai"
)

type Promotion struct {
	Id        string
	Name      PromotionName
	Type      PromotionType
	Desc      PromotionDesc
	Way       PromotionWay
	Tags      []string
	RegUsers  []RegUser
	StartTime int64
	EndTime   int64
	Poster    string
	Host      string
	Intro     string
	Version   int
	IsStatic  bool
}

type RegUser struct {
	User      types.Account
	CreatedAt int64
	Origin    Origin
}

func (r *Promotion) HasRegister(u types.Account) bool {
	for i := range r.RegUsers {
		if u != nil && u.Account() == r.RegUsers[i].User.Account() {
			return true
		}
	}

	return false
}

func (r *Promotion) Status() (string, error) {
	if r.StartTime <= r.EndTime {
		now := utils.Now()
		if now < r.StartTime {
			return PromotionStatusPreparing, nil
		} else if now > r.EndTime {
			return PromotionStatusOver, nil
		}

		return PromotionStatusInProgress, nil
	}

	return "", errors.New("invalid promotion status")
}

func (r *Promotion) Duration() (string, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", time.Unix(r.StartTime, 0).In(loc).Format(layout),
		time.Unix(r.EndTime, 0).In(loc).Format(layout)), nil
}

func (r *Promotion) CountRegUsers() int64 {
	return int64(len(r.RegUsers))
}

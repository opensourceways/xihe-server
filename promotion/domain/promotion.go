package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/sirupsen/logrus"
)

type Promotion struct {
	Id       string
	Name     PromotionName
	Desc     PromotionDesc
	Poster   string
	RegUsers []RegUser
	Duration PromotionDuration
}

type RegUser struct {
	User      types.Account
	CreatedAt int64
	Origin    Origin
}

func (r *Promotion) HasRegister(u types.Account) bool {
	for i := range r.RegUsers {
		if u.Account() == r.RegUsers[i].User.Account() {
			return true
		}
	}

	return false
}

func (r *Promotion) Status() string {
	status, err := r.Duration.PromotionStatus()
	if err != nil {
		logrus.Warnf("get promotion status error: %s", err.Error())
		status = promotionStatusPreparing
	}

	return status
}

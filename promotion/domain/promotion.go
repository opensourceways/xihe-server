package domain

import types "github.com/opensourceways/xihe-server/domain"

type Promotion struct {
	Id       string
	Name     PromotionName
	Desc     PromotionDesc
	RegUsers []RegUser
	Duration PromotionDuration
}

type RegUser struct {
	User      types.Account
	CreatedAt int64
}

func (r *Promotion) HasRegister(u types.Account) bool {
	for i := range r.RegUsers {
		if u.Account() == r.RegUsers[i].User.Account() {
			return true
		}
	}

	return false
}

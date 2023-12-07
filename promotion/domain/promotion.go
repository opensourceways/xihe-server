package domain

import types "github.com/opensourceways/xihe-server/domain"

type Promotion struct {
	Id       string
	Name     PromotionName
	Desc     PromotionDesc
	RegUsers []types.Account
	Duration PromotionDuration
}

func (r *Promotion) HasRegister(u types.Account) bool {
	for i := range r.RegUsers {
		if u.Account() == r.RegUsers[i].Account() {
			return true
		}
	}

	return false
}

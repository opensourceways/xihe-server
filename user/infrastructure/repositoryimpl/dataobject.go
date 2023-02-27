package repositoryimpl

import "github.com/opensourceways/xihe-server/user/domain"

func toUserRegInfoDoc(u *domain.UserRegInfo, doc *DUserRegInfo) {
	*doc = DUserRegInfo{
		Account:  u.Account.Account(),
		Name:     u.Name.Name(),
		City:     u.City.City(),
		Email:    u.Email.Email(),
		Phone:    u.Phone.Phone(),
		Identity: u.Identity.Identity(),
		Province: u.Province.Province(),
		Detail:   u.Detail,
	}
}

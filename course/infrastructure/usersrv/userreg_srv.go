package usersrv

import userdomain "github.com/opensourceways/xihe-server/user/domain"

type userClient interface {
	AddUserRegInfo(u *userdomain.UserRegInfo) error
}

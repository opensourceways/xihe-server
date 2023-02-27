package repository

import "github.com/opensourceways/xihe-server/user/domain"

type UserReg interface {
	AddUserRegInfo(*domain.UserRegInfo) error
}

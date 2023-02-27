package repository

import "github.com/opensourceways/xihe-server/user/domain"

type User interface {
	AddUserRegInfo(*domain.UserRegInfo) error
}

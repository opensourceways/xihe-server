package user

import "github.com/opensourceways/xihe-server/course/domain"

type User interface {
	AddUserRegInfo(*domain.Student) error
}

package user

import "github.com/opensourceways/xihe-server/promotion/domain"

type User interface {
	UpdateRegister(*domain.UserRegistration) error
}

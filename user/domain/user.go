package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type UserRegInfo struct {
	Account  types.Account
	Name     Name
	City     City
	Email    Email
	Phone    Phone
	Identity Identity
	Province Province
	Detail   map[string]string
}

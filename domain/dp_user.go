package domain

import (
	"errors"
	"strings"

	libutil "github.com/opensourceways/community-robot-lib/utils"
	"github.com/opensourceways/xihe-server/utils"
)

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" || strings.ToLower(v) == "root" || !utils.IsUserName(v) {
		return nil, errors.New("invalid user name")
	}

	return dpAccount(v), nil
}

// CreateAccount is usually called internally, such as repository.
func CreateAccount(v string) Account {
	return dpAccount(v)
}

type dpAccount string

func (r dpAccount) Account() string {
	return string(r)
}

// Password
type Password interface {
	Password() string
}

func NewPassword(s string) (Password, error) {
	if n := len(s); n < 8 || n > 20 {
		return nil, errors.New("invalid password")
	}

	part := make([]bool, 4)

	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			part[0] = true
		} else if c >= 'A' && c <= 'Z' {
			part[1] = true
		} else if c >= '0' && c <= '9' {
			part[2] = true
		} else {
			part[3] = true
		}
	}

	i := 0
	for _, b := range part {
		if b {
			i++
		}
	}
	if i < 3 {
		return nil, errors.New(
			"the password must includes three of lowercase, uppercase, digital and special character",
		)
	}

	return dpPassword(s), nil
}

type dpPassword string

func (r dpPassword) Password() string {
	return string(r)
}

// Email
type Email interface {
	Email() string
}

func NewEmail(v string) (Email, error) {
	if v != "" && !libutil.IsValidEmail(v) {
		return nil, errors.New("invalid email")
	}

	return dpEmail(v), nil
}

type dpEmail string

func (r dpEmail) Email() string {
	return string(r)
}

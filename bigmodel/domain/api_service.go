package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type UserApiRecord struct {
	User      types.Account
	ModelName ModelName
	Enabled   bool
	ApplyAt   string
	UpdateAt  string
	Token     string
	Version   int
}

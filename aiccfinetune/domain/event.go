package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type AICCFinetuneCreateEvent struct {
	User  types.Account
	Id    string
	Model string
	Task  string
}

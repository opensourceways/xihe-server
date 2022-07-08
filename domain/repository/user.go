package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type User interface {
	Get(string) (domain.User, error)
	Save(domain.User) error
}

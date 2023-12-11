package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/domain"
)

type Points interface {
	Save(*domain.UserPoints) error
	Update(user types.Account, taskid string, version int) error
	Find(types.Account) (domain.UserPoints, error)
	FindAll() ([]domain.UserPoints, error)
}

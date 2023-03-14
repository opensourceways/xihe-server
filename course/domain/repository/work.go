package repository

import (
	"github.com/opensourceways/xihe-server/course/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type WorkVersion struct {
	Work    domain.Work
	Version int
}

type Work interface {
	GetWork(cid string, account types.Account, asgId string, status domain.WorkStatus) (domain.Work, error)
}

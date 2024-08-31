package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserFinetunes struct {
	domain.FinetuneUserInfo

	Version int
	Data    []domain.FinetuneSummary
}

type Finetune interface {
	Save(domain.Account, *domain.Finetune, int) (string, error)
	Get(*domain.FinetuneIndex) (domain.Finetune, error)
	Delete(*domain.FinetuneIndex) error
	List(user domain.Account) (UserFinetunes, error)

	GetJob(*domain.FinetuneIndex) (domain.FinetuneJob, error)
	SaveJob(*domain.FinetuneIndex, *domain.FinetuneJobInfo) error

	UpdateJobDetail(*domain.FinetuneIndex, *domain.FinetuneJobDetail) error
}

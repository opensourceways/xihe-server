package repository

import (
	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type AICCFinetune interface {
	Save(*domain.AICCFinetune, int) (string, error)
	Get(*domain.AICCFinetuneIndex) (domain.AICCFinetune, error)
	Delete(*domain.AICCFinetuneIndex) error
	List(user types.Account, model domain.ModelName) ([]domain.AICCFinetuneSummary, int, error)

	SaveJob(*domain.AICCFinetuneIndex, *domain.JobInfo) error
	GetJob(*domain.AICCFinetuneIndex) (domain.JobInfo, error)

	UpdateJobDetail(*domain.AICCFinetuneIndex, *domain.JobDetail) error
	GetJobDetail(*domain.AICCFinetuneIndex) (domain.JobDetail, string, error)
}

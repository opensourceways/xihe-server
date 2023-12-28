package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Training interface {
	Save(*domain.UserTraining, int) (string, error)
	Get(*domain.TrainingIndex) (domain.UserTraining, error)
	Delete(*domain.TrainingIndex) error
	List(user domain.Account, projectId string) ([]domain.TrainingSummary, int, error)

	GetTrainingConfig(*domain.TrainingIndex) (domain.TrainingConfig, error)
	GetLastTrainingConfig(*domain.ResourceIndex) (domain.TrainingConfig, error)

	SaveJob(*domain.TrainingIndex, *domain.JobInfo) error
	GetJob(*domain.TrainingIndex) (domain.JobInfo, error)

	UpdateJobDetail(*domain.TrainingIndex, *domain.JobDetail) error
	GetJobDetail(*domain.TrainingIndex) (domain.JobDetail, string, error)
}

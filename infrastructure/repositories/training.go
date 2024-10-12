package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type TrainingMapper interface {
	Insert(*UserTrainingDO, int) (string, error)
	Delete(*TrainingIndexDO) error
	Get(*TrainingIndexDO) (TrainingDetailDO, error)
	GetTrainingConfig(*TrainingIndexDO) (TrainingConfigDO, error)
	GetLastTrainingConfig(*ResourceIndexDO) (TrainingConfigDO, error)
	List(user, projectId string) ([]TrainingSummaryDO, int, error)
	UpdateJobInfo(*TrainingIndexDO, *TrainingJobInfoDO) error
	GetJobInfo(*TrainingIndexDO) (TrainingJobInfoDO, error)
	UpdateJobDetail(*TrainingIndexDO, *TrainingJobDetailDO) error
	GetJobDetail(*TrainingIndexDO) (TrainingJobDetailDO, string, error)
}

func NewTrainingRepository(mapper TrainingMapper) repository.Training {
	return training{mapper}
}

type training struct {
	mapper TrainingMapper
}

func (impl training) Save(ut *domain.UserTraining, version int) (string, error) {
	if ut.Id != "" {
		return "", errors.New("must be a new training")
	}

	do := impl.toUserTrainingDO(ut)

	v, err := impl.mapper.Insert(&do, version)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl training) Delete(info *domain.TrainingIndex) error {
	do := impl.toTrainingIndexDO(info)

	if err := impl.mapper.Delete(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl training) Get(info *domain.TrainingIndex) (obj domain.UserTraining, err error) {
	do := impl.toTrainingIndexDO(info)

	v, err := impl.mapper.Get(&do)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toUserTraining(info, &obj)
	}

	return
}

func (impl training) GetTrainingConfig(info *domain.TrainingIndex) (domain.TrainingConfig, error) {
	do := impl.toTrainingIndexDO(info)

	v, err := impl.mapper.GetTrainingConfig(&do)
	if err != nil {
		return domain.TrainingConfig{}, convertError(err)
	}

	return v.toTrainingConfig()
}

func (impl training) GetLastTrainingConfig(res *domain.ResourceIndex) (domain.TrainingConfig, error) {
	do := ToResourceIndexDO(res)
	v, err := impl.mapper.GetLastTrainingConfig(&do)
	if err != nil {
		return domain.TrainingConfig{}, convertError(err)
	}

	return v.toTrainingConfig()
}

func (impl training) List(user domain.Account, projectId string) (
	r []domain.TrainingSummary, version int, err error,
) {
	v, version, err := impl.mapper.List(user.Account(), projectId)
	if err != nil {
		err = convertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r = make([]domain.TrainingSummary, len(v))
	for i := range v {
		if err = v[i].toTrainingSummary(&r[i]); err != nil {
			return
		}
	}

	return
}

func (impl training) SaveJob(info *domain.TrainingIndex, job *domain.JobInfo) error {
	do := impl.toTrainingIndexDO(info)

	if err := impl.mapper.UpdateJobInfo(&do, job); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl training) GetJob(info *domain.TrainingIndex) (job domain.JobInfo, err error) {
	t := impl.toTrainingIndexDO(info)

	job, err = impl.mapper.GetJobInfo(&t)
	if err != nil {
		err = convertError(err)
	}

	return
}

func (impl training) GetJobDetail(info *domain.TrainingIndex) (
	job domain.JobDetail, endpoint string, err error,
) {
	t := impl.toTrainingIndexDO(info)

	job, endpoint, err = impl.mapper.GetJobDetail(&t)
	if err != nil {
		err = convertError(err)
	}

	return
}

func (impl training) UpdateJobDetail(info *domain.TrainingIndex, detail *domain.JobDetail) error {
	do := impl.toTrainingIndexDO(info)

	if err := impl.mapper.UpdateJobDetail(&do, detail); err != nil {
		return convertError(err)
	}

	return nil
}

package app

import (
	"errors"
	"strconv"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/domain/training"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

const (
	trainingStatusScheduling     = "scheduling"
	trainingStatusScheduleFailed = "schedule_failed"
)

type JobDetail = domain.JobDetail
type TrainingIndex = domain.TrainingIndex
type TrainingConfig = domain.TrainingConfig

type TrainingService interface {
	Create(*TrainingCreateCmd) (string, error)
	Recreate(*TrainingIndex) (string, error)
	UpdateJobDetail(*TrainingIndex, *JobDetail) error
	List(user domain.Account, projectId string) ([]TrainingSummaryDTO, error)
	Get(*TrainingIndex) (TrainingDTO, string, error)
	GetLastTrainingConfig(*ResourceIndexCmd) (dto TrainingConfigDTO, code string, err error)
	Delete(*TrainingIndex) error
	Terminate(*TrainingIndex) error
	GetLogDownloadURL(*TrainingIndex) (string, string, error)
	GetOutputDownloadURL(*TrainingIndex) (string, string, error)
	CreateTrainingJob(*TrainingIndex, string, bool) (bool, error)
}

func NewTrainingService(
	train training.Training,
	repo repository.Training,
	sender message.MessageProducer,
	maxTrainingRecordNum int,
) TrainingService {
	return trainingService{
		train:  train,
		repo:   repo,
		sender: sender,

		maxTrainingRecordNum: maxTrainingRecordNum,
	}
}

type trainingService struct {
	log    *logrus.Entry
	train  training.Training
	repo   repository.Training
	sender message.MessageProducer

	maxTrainingRecordNum int
}

func (s trainingService) isJobDone(status string) bool {
	return status != "" && (s.train.IsJobDone(status) || status == trainingStatusScheduleFailed)
}

func (s trainingService) Create(cmd *TrainingCreateCmd) (string, error) {
	return s.create(cmd.User, cmd.ProjectId, cmd.toTrainingConfig())
}

func (s trainingService) Recreate(info *TrainingIndex) (string, error) {
	v, err := s.repo.GetTrainingConfig(info)
	if err != nil {
		return "", err
	}

	v.Name, err = domain.NewTrainingName(
		v.Name.TrainingName() + "-" + strconv.FormatInt(utils.Now(), 10),
	)
	if err != nil {
		return "", err
	}

	return s.create(info.Project.Owner, info.Project.Id, &v)
}

func (s trainingService) create(
	user domain.Account, projectId string, config *TrainingConfig,
) (string, error) {
	v, version, err := s.repo.List(user, projectId)
	if err != nil {
		return "", err
	}

	if len(v) >= s.maxTrainingRecordNum {
		return "", ErrorExccedMaxTrainingRecord{
			errors.New("exceed max training num"),
		}
	}

	for i := range v {
		if !s.isJobDone(v[i].Status) {
			return "", ErrorOnlyOneRunningTraining{
				errors.New("a training is running"),
			}
		}
	}

	t := domain.UserTraining{
		Owner:          user,
		ProjectId:      projectId,
		CreatedAt:      utils.Now(),
		TrainingConfig: *config,
	}

	r, err := s.repo.Save(&t, version)
	if err != nil {
		return "", err
	}

	// send message
	index := TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: user,
			Id:    projectId,
		},
		TrainingId: r,
	}

	err = s.sender.SendTrainingCreated(&domain.TrainingCreatedEvent{
		Account:        user,
		TrainingIndex:  index,
		TrainingInputs: config.Inputs,
	})
	if err != nil {
		s.log.Errorf("send message of creating training failed, err:%s", err.Error())
	}

	return r, nil
}

func (s trainingService) List(user domain.Account, projectId string) ([]TrainingSummaryDTO, error) {
	v, _, err := s.repo.List(user, projectId)
	if err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]TrainingSummaryDTO, len(v))
	for i := range v {
		s.toTrainingSummaryDTO(&v[i], &r[i])
	}

	return r, nil
}

func (s trainingService) Get(info *TrainingIndex) (dto TrainingDTO, code string, err error) {
	data, err := s.repo.Get(info)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			code = ErrorTrainNotFound
		}

		return
	}

	link := ""
	if job := &data.Job; job.Endpoint != "" && job.JobId != "" {
		link, err = s.train.GetLogPreviewURL(job.Endpoint, job.JobId)
		if err != nil {
			return
		}
	}

	s.toTrainingDTO(&dto, &data, link)

	return
}

func (s trainingService) GetLastTrainingConfig(cmd *ResourceIndexCmd) (dto TrainingConfigDTO, code string, err error) {
	resourceIndex := domain.ResourceIndex(*cmd)
	data, err := s.repo.GetLastTrainingConfig(&resourceIndex)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			code = ErrorTrainNotFound
		}

		return
	}

	dto.toDTO(&data)

	return
}

func (s trainingService) UpdateJobDetail(info *TrainingIndex, v *JobDetail) error {
	return s.repo.UpdateJobDetail(info, v)
}

func (s trainingService) Delete(info *TrainingIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			return nil
		}
		return err
	}

	if job.JobId != "" {
		err = s.train.DeleteJob(job.Endpoint, job.JobId)
		if err != nil {
			// ignore 404
			return err
		}
	}

	return s.repo.Delete(info)
}

func (s trainingService) Terminate(info *TrainingIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil || job.JobId == "" {
		return err
	}

	return s.train.TerminateJob(job.Endpoint, job.JobId)
}

func (s trainingService) GetLogDownloadURL(info *TrainingIndex) (
	link string, code string, err error,
) {
	detail, endpoint, err := s.repo.GetJobDetail(info)
	if err != nil {
		return
	}

	if detail.LogPath == "" {
		code = ErrorTrainNoLog
		err = errors.New("not ready")
	} else {
		link, err = s.train.GetFileDownloadURL(endpoint, detail.LogPath)
	}

	return
}

func (s trainingService) GetOutputDownloadURL(info *TrainingIndex) (
	link string, code string, err error,
) {
	detail, endpoint, err := s.repo.GetJobDetail(info)
	if err != nil {
		return
	}

	if detail.OutputPath == "" {
		code = ErrorTrainNoOutput
		err = errors.New("not ready")
	} else {
		link, err = s.train.GetFileDownloadURL(endpoint, detail.OutputPath)
	}

	return
}

func (s trainingService) CreateTrainingJob(
	info *TrainingIndex, endpoint string, lastChance bool,
) (retry bool, err error) {
	retry, err = s.createTrainingJob(info, endpoint)
	if err == nil {
		return
	}

	if lastChance {
		err = s.repo.UpdateJobDetail(info, &JobDetail{
			Status: trainingStatusScheduleFailed,
			Error:  err.Error(),
		})
	}

	return
}

func (s trainingService) createTrainingJob(info *TrainingIndex, endpoint string) (
	retry bool, err error,
) {
	data, err := s.repo.Get(info)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			err = errorUnavailableTraining{
				errors.New("training is not exist"),
			}
		} else {
			retry = true
		}

		return
	}

	if data.Job.JobId != "" {
		return false, nil
	}

	v, err := s.train.CreateJob(endpoint, info, &data.TrainingConfig)
	if err != nil {
		retry = true

		return
	}

	if err1 := s.repo.SaveJob(info, &v); err1 != nil {
		s.log.Errorf(
			"create training(%s) job(%s) successfully, but save db err:%s",
			info.TrainingId, v.JobId, err1.Error(),
		)
	}

	return
}

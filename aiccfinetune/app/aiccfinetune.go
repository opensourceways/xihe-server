package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain/aiccfinetune"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain/message"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain/repository"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain/uploader"
	"github.com/opensourceways/xihe-server/app"

	types "github.com/opensourceways/xihe-server/domain"
	orepo "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	trainingStatusScheduling     = "scheduling"
	trainingStatusScheduleFailed = "schedule_failed"
)

type JobDetail = domain.JobDetail
type AICCFinetuneIndex = domain.AICCFinetuneIndex
type AICCFinetuneConfig = domain.AICCFinetuneConfig

type AICCFinetuneService interface {
	Create(*AICCFinetuneCreateCmd) (string, error)
	UpdateJobDetail(*AICCFinetuneIndex, *JobDetail) error
	List(user types.Account, model domain.ModelName) ([]AICCFinetuneSummaryDTO, error)
	Get(*AICCFinetuneIndex) (AICCFinetuneDTO, string, error)
	Delete(*AICCFinetuneIndex) error
	Terminate(*AICCFinetuneIndex) error
	GetLogDownloadURL(*AICCFinetuneIndex) (string, string, error)
	GetOutputDownloadURL(*AICCFinetuneIndex) (string, string, error)
	CreateAICCFinetuneJob(*AICCFinetuneIndex, string, bool) (bool, error)

	UploadData(*UploadDataCmd) (UploadDataDTO, error)
}

func NewAICCFinetuneService(
	af aiccfinetune.AICCFinetuneServer,
	sender message.AICCFinetuneMessageProducer,
	uploader uploader.DataFileUploader,
	repo repository.AICCFinetune,
	maxTrainingRecordNum int,
) AICCFinetuneService {
	return aiccFinetuneService{
		af:                   af,
		sender:               sender,
		uploader:             domain.NewUploadService(uploader),
		repo:                 repo,
		maxTrainingRecordNum: maxTrainingRecordNum,
	}
}

type aiccFinetuneService struct {
	af                   aiccfinetune.AICCFinetuneServer
	sender               message.AICCFinetuneMessageProducer
	uploader             domain.UploadService
	repo                 repository.AICCFinetune
	maxTrainingRecordNum int
}

func (s aiccFinetuneService) isJobDone(status string) bool {
	return status != "" && (s.af.IsJobDone(status) || status == trainingStatusScheduleFailed)
}

func (s aiccFinetuneService) Create(cmd *AICCFinetuneCreateCmd) (string, error) {
	return s.create(cmd.User, cmd.Model, cmd.Task, cmd.toAICCFinetuneConfig())
}

func (s aiccFinetuneService) create(
	user types.Account, model domain.ModelName, task domain.FinetuneTask, config *AICCFinetuneConfig,
) (string, error) {
	v, version, err := s.repo.List(user, model)
	if err != nil {
		return "", err
	}

	if len(v) >= s.maxTrainingRecordNum {
		return "", errors.New("exceed max finetune num")
	}

	for i := range v {
		if !s.isJobDone(v[i].Status) {
			return "", errors.New("a finetune is running")
		}
	}

	t := domain.AICCFinetune{
		User:      user,
		CreatedAt: utils.Now(),
		Model:     model,
		Task:      task,

		AICCFinetuneConfig: *config,
	}

	r, err := s.repo.Save(&t, version)
	if err != nil {
		return "", err
	}

	if err = s.sender.SendAICCFinetuneCreateMsg(&domain.AICCFinetuneCreateEvent{
		Id:    r,
		User:  user,
		Model: model.ModelName(),
		Task:  task.FinetuneTask(),
	}); err != nil {
		return "", err
	}

	return r, nil
}

func (s aiccFinetuneService) List(user types.Account, model domain.ModelName) ([]AICCFinetuneSummaryDTO, error) {
	v, _, err := s.repo.List(user, model)
	if err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]AICCFinetuneSummaryDTO, len(v))
	for i := range v {
		s.toAICCFinetuneSummaryDTO(&v[i], &r[i])
	}

	return r, nil
}

func (s aiccFinetuneService) Get(info *AICCFinetuneIndex) (dto AICCFinetuneDTO, code string, err error) {
	data, err := s.repo.Get(info)
	if err != nil {
		if orepo.IsErrorResourceNotExists(err) {
			code = app.ErrorAICCFinetuneNotFound
		}

		return
	}

	link := ""
	if job := &data.Job; job.Endpoint != "" && job.JobId != "" {
		link, err = s.af.GetLogPreviewURL(job.Endpoint, job.JobId)
		if err != nil {
			return
		}
	}

	s.toAICCFinetuneDTO(&dto, &data, link)

	return
}

func (s aiccFinetuneService) UpdateJobDetail(info *AICCFinetuneIndex, v *JobDetail) error {
	return s.repo.UpdateJobDetail(info, v)
}

func (s aiccFinetuneService) Delete(info *AICCFinetuneIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil {
		if orepo.IsErrorResourceNotExists(err) {
			return nil
		}
		return err
	}

	if job.JobId != "" {
		err = s.af.DeleteJob(job.Endpoint, job.JobId)
		if err != nil {
			// ignore 404
			return err
		}
	}

	return s.repo.Delete(info)
}

func (s aiccFinetuneService) Terminate(info *AICCFinetuneIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil || job.JobId == "" {
		return err
	}

	return s.af.TerminateJob(job.JobId)
}

func (s aiccFinetuneService) GetLogDownloadURL(info *AICCFinetuneIndex) (
	link string, code string, err error,
) {
	detail, endpoint, err := s.repo.GetJobDetail(info)
	if err != nil {
		return
	}

	if detail.LogPath == "" {
		code = app.ErrorAICCFinetuneNoLog
		err = errors.New("not ready")
	} else {
		link, err = s.af.GetFileDownloadURL(endpoint, detail.LogPath)
	}

	return
}

func (s aiccFinetuneService) GetOutputDownloadURL(info *AICCFinetuneIndex) (
	link string, code string, err error,
) {
	detail, endpoint, err := s.repo.GetJobDetail(info)
	if err != nil {
		return
	}

	if detail.OutputPath == "" {
		code = app.ErrorTrainNoOutput
		err = errors.New("not ready")
	} else {
		link, err = s.af.GetFileDownloadURL(endpoint, detail.OutputPath)
	}

	return
}

func (s aiccFinetuneService) CreateAICCFinetuneJob(
	info *AICCFinetuneIndex, endpoint string, lastChance bool,
) (retry bool, err error) {
	retry, err = s.createAICCFinetuneJob(info, endpoint)
	if err == nil {
		return
	}

	if lastChance {
		s.repo.UpdateJobDetail(info, &JobDetail{
			Status: trainingStatusScheduleFailed,
			Error:  err.Error(),
		})
	}

	return
}

func (s aiccFinetuneService) createAICCFinetuneJob(info *AICCFinetuneIndex, endpoint string) (
	retry bool, err error,
) {
	data, err := s.repo.Get(info)
	if err != nil {
		if orepo.IsErrorResourceNotExists(err) {
			err = errors.New("aicc finetune is not exist")
		} else {
			retry = true
		}

		return
	}

	if data.Job.JobId != "" {
		return false, nil
	}

	v, err := s.af.CreateJob(endpoint, info, &data)
	if err != nil {
		retry = true

		return
	}

	if err1 := s.repo.SaveJob(info, &v); err1 != nil {
		return
	}

	return
}

func (s aiccFinetuneService) UploadData(cmd *UploadDataCmd) (dto UploadDataDTO, err error) {
	err = s.uploader.Upload(cmd.Data, cmd.FileName, cmd.User.Account(), cmd.Model.ModelName(), cmd.Task.FinetuneTask())

	dto.FileName = cmd.FileName
	dto.UploadAt = utils.Now()
	dto.Status = "failed"
	if err == nil {
		dto.Status = "success"
		return
	}
	return
}

type AICCFinetuneInternalService interface {
	UpdateJobDetail(*AICCFinetuneIndex, *JobDetail) error
}

type aiccfinetuneInternalService struct {
	repo repository.AICCFinetune
}

func NewAICCFinetuneInternalService(
	repo repository.AICCFinetune,
) AICCFinetuneInternalService {
	return aiccfinetuneInternalService{
		repo: repo,
	}
}

func (s aiccfinetuneInternalService) UpdateJobDetail(info *AICCFinetuneIndex, v *JobDetail) error {
	return s.repo.UpdateJobDetail(info, v)
}

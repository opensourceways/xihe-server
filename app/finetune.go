package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/finetune"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

type FinetuneIndex = domain.FinetuneIndex
type FinetuneConfig = domain.FinetuneConfig
type FinetuneJobDetail = domain.FinetuneJobDetail

type FinetuneService interface {
	Create(*FinetuneCreateCmd) (string, string, error)
	List(user domain.Account) (UserFinetunesDTO, string, error)
	Delete(*FinetuneIndex) error
	Terminate(*FinetuneIndex) error
	GetJobInfo(*FinetuneIndex) (FinetuneJobDTO, string, error)
}

func NewFinetuneService(
	fs finetune.Finetune,
	repo repository.Finetune,
	sender message.Sender,
) FinetuneService {
	return finetuneService{
		fs:     fs,
		repo:   repo,
		sender: sender,
	}
}

type finetuneService struct {
	fs     finetune.Finetune
	repo   repository.Finetune
	sender message.Sender
}

func (s finetuneService) isJobDone(status string) bool {
	return status != "" && (s.fs.IsJobDone(status) || status == trainingStatusScheduleFailed)
}

func (s finetuneService) Create(cmd *FinetuneCreateCmd) (
	fid string, code string, err error,
) {
	user := cmd.User

	v, err := s.repo.List(user)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			code = ErrorFinetuneNoPermission
		}

		return
	}

	if utils.IsExpiry(v.Expiry) {
		code = ErrorFinetuneExpiry
		err = errors.New("it is expiry")

		return
	}

	if len(v.Datas) >= appConfig.FinetuneMaxNum {
		code = ErrorFinetuneExccedMaxNum
		err = errors.New("exceed max record num")

		return
	}

	for i := range v.Datas {
		if !s.isJobDone(v.Datas[i].Status) {
			code = ErrorFinetuneRunningJobExists
			err = errors.New("a job is running")

			return
		}
	}

	t := domain.Finetune{
		CreatedAt:      utils.Now(),
		FinetuneConfig: *cmd.toFinetuneConfig(),
	}

	if fid, err = s.repo.Save(user, &t, v.Version); err != nil {
		return
	}

	// send message
	err1 := s.sender.CreateFinetune(&FinetuneIndex{
		Owner: user,
		Id:    fid,
	})
	if err1 != nil {
		logrus.Errorf("send message of creating finetune failed, err:%s", err.Error())
	}

	return
}

func (s finetuneService) List(user domain.Account) (
	r UserFinetunesDTO, code string, err error,
) {
	v, err := s.repo.List(user)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			code = ErrorFinetuneNoPermission
		}

		return
	}

	r.Expiry = v.Expiry

	if len(v.Datas) == 0 {
		return
	}

	datas := make([]FinetuneSummaryDTO, len(v.Datas))
	for i := range v.Datas {
		s.toFinetuneSummaryDTO(&v.Datas[i], &datas[i])
	}
	r.Datas = datas

	return
}

func (s finetuneService) Delete(info *FinetuneIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			return nil
		}
		return err
	}

	if job.JobId != "" {
		if err = s.fs.DeleteJob(job.JobId); err != nil {
			// ignore 404
			return err
		}
	}

	return s.repo.Delete(info)
}

func (s finetuneService) Terminate(info *FinetuneIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil || job.JobId == "" {
		return err
	}

	if !s.fs.CanTerminate(job.Status) {
		return errors.New("can't terminate now")
	}

	return s.fs.TerminateJob(job.JobId)
}

func (s finetuneService) GetJobInfo(index *FinetuneIndex) (
	dto FinetuneJobDTO, code string, err error,
) {
	job, err := s.repo.GetJob(index)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			code = ErrorFinetuneNotFound
		}

		return
	}

	dto.IsDone = s.isJobDone(job.Status)

	if job.JobId != "" {
		dto.LogPreviewURL, err = s.fs.GetLogPreviewURL(job.JobId)
	}

	return
}

// FinetuneInternalService
type FinetuneInternalService interface {
	UpdateJobDetail(*FinetuneIndex, *FinetuneJobDetail) error
}

func NewFinetuneInternalService(
	repo repository.Finetune,
) FinetuneInternalService {
	return finetuneInternalService{
		repo: repo,
	}
}

type finetuneInternalService struct {
	repo repository.Finetune
}

func (s finetuneInternalService) UpdateJobDetail(info *FinetuneIndex, v *FinetuneJobDetail) error {
	return s.repo.UpdateJobDetail(info, v)
}

// FinetuneMessageService
type FinetuneMessageService interface {
	CreateFinetuneJob(*FinetuneIndex, bool) (bool, error)
}

func NewFinetuneMessageService(
	fs finetune.Finetune,
	repo repository.Finetune,
) FinetuneMessageService {
	return finetuneMessageService{
		fs:   fs,
		repo: repo,
	}
}

type finetuneMessageService struct {
	fs   finetune.Finetune
	repo repository.Finetune
}

func (s finetuneMessageService) CreateFinetuneJob(
	info *FinetuneIndex, lastChance bool,
) (retry bool, err error) {
	if retry, err = s.createFinetuneJob(info); err == nil || !lastChance {
		return
	}

	err = s.repo.UpdateJobDetail(info, &FinetuneJobDetail{
		Status: trainingStatusScheduleFailed,
		Error:  err.Error(),
	})

	return
}

func (s finetuneMessageService) createFinetuneJob(info *FinetuneIndex) (
	retry bool, err error,
) {
	data, err := s.repo.Get(info)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			err = errors.New("finetune is not exist")
		} else {
			retry = true
		}

		return
	}

	if data.Job.JobId != "" {
		return
	}

	v, err := s.fs.CreateJob(info, &data.FinetuneConfig)
	if err != nil {
		retry = true

		return
	}

	if err1 := s.repo.SaveJob(info, &v); err1 != nil {
		logrus.Errorf(
			"create finetune(%s) job(%s) successfully, but save db err:%s",
			info.Id, v.JobId, err1.Error(),
		)
	}

	return
}

package finetuneimpl

import (
	"github.com/opensourceways/xihe-finetune/sdk"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/finetune"
)

func NewFinetune(cfg *Config) finetune.Finetune {
	return &finetuneImpl{
		cli:                sdk.New(cfg.Endpoint),
		doneStatus:         sets.NewString(cfg.JobDoneStatus...),
		canTerminateStatus: sets.NewString(cfg.CanTerminateStatus...),
	}
}

type finetuneImpl struct {
	cli                sdk.Finetune
	doneStatus         sets.String
	canTerminateStatus sets.String
}

func (impl *finetuneImpl) IsJobDone(status string) bool {
	return impl.doneStatus.Has(status)
}

func (impl *finetuneImpl) CanTerminate(status string) bool {
	return impl.canTerminateStatus.Has(status)
}

func (impl *finetuneImpl) CreateJob(info *domain.FinetuneIndex, cfg *domain.FinetuneConfig) (
	job domain.FinetuneJobInfo, err error,
) {
	p := cfg.Param
	opt := sdk.FinetuneCreateOption{
		User:            info.Owner.Account(),
		Id:              info.Id,
		Name:            cfg.Name.FinetuneName(),
		Task:            p.Task(),
		Model:           p.Model(),
		Hyperparameters: p.Hyperparameters(),
	}

	v, err := impl.cli.Create(&opt)
	if err == nil {
		job.JobId = v.JobId
	}

	return
}

func (impl *finetuneImpl) DeleteJob(jobId string) error {
	return impl.cli.Delete(jobId)
}

func (impl *finetuneImpl) TerminateJob(jobId string) error {
	return impl.cli.Terminate(jobId)
}

func (impl *finetuneImpl) GetLogPreviewURL(jobId string) (r string, err error) {
	v, err := impl.cli.GetLogDownloadURL(jobId)
	if err == nil {
		r = v.URL
	}

	return
}

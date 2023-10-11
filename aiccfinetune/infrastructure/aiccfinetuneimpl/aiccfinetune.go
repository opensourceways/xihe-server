package aiccfinetuneimpl

import (
	"strings"

	"github.com/opensourceways/xihe-aicc-finetune/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain/aiccfinetune"
)

func NewAICCFinetune(cfg *Config) aiccfinetune.AICCFinetuneServer {
	return &aiccFinetuneImpl{
		doneStatus: sets.NewString(cfg.JobDoneStatus...),
		endpoint:   cfg.Endpoint,
	}
}

type aiccFinetuneImpl struct {
	doneStatus sets.String
	endpoint   string
}

func (impl *aiccFinetuneImpl) IsJobDone(status string) bool {
	return impl.doneStatus.Has(status)
}

func (impl *aiccFinetuneImpl) CreateJob(endpoint string, info *domain.AICCFinetuneIndex, t *domain.AICCFinetune) (
	job domain.JobInfo, err error,
) {
	opt := sdk.AICCFinetuneCreateOption{
		FinetuneId:      info.FinetuneId,
		User:            info.User.Account(),
		Model:           t.Model.ModelName(),
		Name:            t.Name.FinetuneName(),
		Desc:            t.Desc.FinetuneDesc(),
		Task:            t.Task.FinetuneTask(),
		Env:             impl.toKeyValue(t.Env),
		Hyperparameters: impl.toKeyValue(t.Hyperparameters),
	}

	logrus.Debugf(
		"create job, endpoint:%s, training:%s, opt:%#v",
		endpoint, info.FinetuneId, opt,
	)

	if t.Desc != nil {
		opt.Desc = t.Desc.FinetuneDesc()
	}

	cli := sdk.NewAICCFinetuneCenter(endpoint)

	v, err := cli.CreateAICCFinetune(&opt)
	if err != nil {
		return
	}

	job.Endpoint = endpoint
	job.JobId = v.JobId
	job.LogDir = v.LogDir
	job.OutputDir = v.OutputDir

	return
}

func (impl *aiccFinetuneImpl) DeleteJob(endpoint, jobId string) error {
	cli := sdk.NewAICCFinetuneCenter(endpoint)

	return cli.DeleteAICCFinetune(jobId)
}

func (impl *aiccFinetuneImpl) TerminateJob(jobId string) error {
	cli := sdk.NewAICCFinetuneCenter(impl.endpoint)
	jobId = "397d6e3d-e2df-4f70-834b-4f4b414cd9c4"
	return cli.TerminateeAICCFinetune(jobId)
}

func (impl *aiccFinetuneImpl) GetLogPreviewURL(endpoint, jobId string) (string, error) {
	cli := sdk.NewAICCFinetuneCenter(endpoint)

	v, err := cli.GetLogDownloadURL(jobId)
	if err != nil {
		return "", err
	}

	return v.URL, nil
}

func (impl *aiccFinetuneImpl) GetFileDownloadURL(endpoint, file string) (string, error) {
	cli := sdk.NewAICCFinetuneCenter(endpoint)

	v, err := cli.GetResultDownloadURL(
		"no_need", strings.ReplaceAll(file, "/", "%2F"),
	)
	if err != nil {
		return "", err
	}

	return v.URL, nil
}

func (impl *aiccFinetuneImpl) toKeyValue(kv []domain.KeyValue) []sdk.KeyValue {
	if len(kv) == 0 {
		return nil
	}

	r := make([]sdk.KeyValue, len(kv))

	for i := range kv {
		r[i] = sdk.KeyValue{
			Key:   kv[i].Key.CustomizedKey(),
			Value: kv[i].Value.CustomizedValue(),
		}
	}

	return r
}

// type service struct {
// 	obs obsService
// }

// func (s *aiccFinetuneImpl) Upload(data io.Reader, path string) error {
// 	return s.obs.createObject(data, path)
// }

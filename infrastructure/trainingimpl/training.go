package trainingimpl

import (
	"strings"

	"github.com/opensourceways/xihe-training-center/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/training"
)

func NewTraining(cfg *Config) training.Training {
	return &trainingImpl{
		doneStatus: sets.NewString(cfg.JobDoneStatus...),
	}
}

type trainingImpl struct {
	doneStatus sets.String
}

func (impl *trainingImpl) IsJobDone(status string) bool {
	return impl.doneStatus.Has(status)
}

func (impl *trainingImpl) CreateJob(endpoint string, info *domain.TrainingIndex, t *domain.TrainingConfig) (
	job domain.JobInfo, err error,
) {
	opt := sdk.TrainingCreateOption{
		User:            info.Project.Owner.Account(),
		ProjectId:       info.Project.Id,
		TrainingId:      info.TrainingId,
		ProjectName:     t.ProjectName.ResourceName(),
		ProjectRepoId:   t.ProjectRepoId,
		Name:            t.Name.TrainingName(),
		CodeDir:         t.CodeDir.Directory(),
		BootFile:        t.BootFile.FilePath(),
		Compute:         impl.toCompute(&t.Compute),
		Env:             impl.toKeyValue(t.Env),
		Inputs:          impl.toInput(t.Inputs),
		EnableAim:       t.EnableAim,
		EnableOutput:    t.EnableOutput,
		Hyperparameters: impl.toKeyValue(t.Hyperparameters),
	}

	logrus.Debugf(
		"create job, endpoint:%s, training:%s, opt:%#v",
		endpoint, info.TrainingId, opt,
	)

	if t.Desc != nil {
		opt.Desc = t.Desc.TrainingDesc()
	}

	cli := sdk.NewTrainingCenter(endpoint)

	v, err := cli.CreateTraining(&opt)
	if err != nil {
		return
	}

	job.Endpoint = endpoint
	job.JobId = v.JobId
	job.LogDir = v.LogDir
	job.AimDir = v.AimDir
	job.OutputDir = v.OutputDir

	return
}

func (impl *trainingImpl) DeleteJob(endpoint, jobId string) error {
	cli := sdk.NewTrainingCenter(endpoint)

	return cli.DeleteTraining(jobId)
}

func (impl *trainingImpl) TerminateJob(endpoint, jobId string) error {
	cli := sdk.NewTrainingCenter(endpoint)

	return cli.TerminateTraining(jobId)
}

func (impl *trainingImpl) GetLogPreviewURL(endpoint, jobId string) (string, error) {
	cli := sdk.NewTrainingCenter(endpoint)

	v, err := cli.GetLogDownloadURL(jobId)
	if err != nil {
		return "", err
	}

	return v.URL, nil
}

func (impl *trainingImpl) GetFileDownloadURL(endpoint, file string) (string, error) {
	cli := sdk.NewTrainingCenter(endpoint)

	v, err := cli.GetResultDownloadURL(
		"no_need", strings.ReplaceAll(file, "/", "%2F"),
	)
	if err != nil {
		return "", err
	}

	return v.URL, nil
}

func (impl *trainingImpl) toCompute(c *domain.Compute) sdk.Compute {
	return sdk.Compute{
		Type:    c.Type.ComputeType(),
		Flavor:  c.Flavor.ComputeFlavor(),
		Version: c.Version.ComputeVersion(),
	}
}

func (impl *trainingImpl) toKeyValue(kv []domain.KeyValue) []sdk.KeyValue {
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

func (impl *trainingImpl) toInput(v []domain.Input) []sdk.Input {
	r := make([]sdk.Input, len(v))

	for i := range v {
		item := &v[i]

		r[i] = sdk.Input{
			Key: item.Key.CustomizedKey(),
			Value: sdk.ResourceRef{
				Owner:  item.User.Account(),
				RepoId: item.RepoId,
				File:   item.File.InputeFilePath(),
			},
		}
	}

	return r
}

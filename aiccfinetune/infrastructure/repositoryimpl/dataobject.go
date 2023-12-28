package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

func (doc aiccFinetuneItem) toAICCFinetuneSummary(s *domain.AICCFinetuneSummary) (err error) {

	if s.Name, err = domain.NewFinetuneName(doc.Name); err != nil {
		return
	}

	if s.Desc, err = domain.NewFinetuneDesc(doc.Desc); err != nil {
		return
	}
	s.CreatedAt = doc.CreatedAt
	s.Duration = doc.JobDetail.Duration
	s.Status = doc.JobDetail.Status
	s.Id = doc.Id
	s.Error = doc.JobDetail.Error
	s.Task = doc.Task

	return
}

func (doc dJobInfo) toAICCFinetuneJobInfo(s *domain.JobInfo) (err error) {
	s.Endpoint = doc.Endpoint
	s.JobId = doc.JobId
	s.LogDir = doc.LogDir
	s.OutputDir = doc.OutputDir
	return
}

func (doc dJobDetail) toAICCFinetuneJobDetail(s *domain.JobDetail) (err error) {
	s.Status = doc.Status
	s.Duration = doc.Duration
	s.Error = doc.Error
	s.OutputPath = doc.OutputPath
	s.LogPath = doc.LogPath
	return
}

func (doc dAICCFinetune) toAICCFinetuneDO(f *domain.AICCFinetune) (err error) {
	f.Id = doc.Items[0].Id

	if f.User, err = types.NewAccount(doc.User); err != nil {
		return
	}

	if f.Task, err = domain.NewFinetuneTask(doc.Items[0].Task); err != nil {
		return
	}

	if f.Model, err = domain.NewModelName(doc.Items[0].Model); err != nil {
		return
	}

	var jobInfo domain.JobInfo
	if err = doc.Items[0].Job.toAICCFinetuneJobInfo(&jobInfo); err != nil {
		return
	}
	f.Job = jobInfo

	if f.Name, err = domain.NewFinetuneName(doc.Items[0].Name); err != nil {
		return
	}

	if f.Desc, err = domain.NewFinetuneDesc(doc.Items[0].Desc); err != nil {
		return
	}

	var jobDetail domain.JobDetail
	if err = doc.Items[0].JobDetail.toAICCFinetuneJobDetail(&jobDetail); err != nil {
		return
	}

	f.JobDetail = jobDetail

	f.CreatedAt = doc.Items[0].CreatedAt

	if f.Hyperparameters, err = toKeyValues(doc.Items[0].Hyperparameters); err != nil {
		return
	}

	return
}

func toKeyValues(kv []dKeyValue) ([]domain.KeyValue, error) {
	n := len(kv)
	if n == 0 {
		return nil, nil
	}

	r := make([]domain.KeyValue, n)

	for i, v := range kv {
		if s, err := domain.NewCustomizedKey(v.Key); err != nil {
			return nil, err
		} else {
			r[i].Key = s
		}

		if a, err := domain.NewCustomizedValue(v.Value); err != nil {
			return nil, err
		} else {
			r[i].Value = a
		}

	}

	return r, nil
}

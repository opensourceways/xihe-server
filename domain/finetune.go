package domain

import (
	"errors"

	"github.com/opensourceways/xihe-server/utils"
)

type Finetune struct {
	Id        string
	CreatedAt int64

	FinetuneConfig

	// following fileds is not under the controlling of version
	Job       FinetuneJobInfo
	JobDetail FinetuneJobDetail
}

type FinetuneConfig struct {
	Name  FinetuneName
	Param FinetuneParameter
}

type FinetuneParameter interface {
	Model() string
	Task() string
	Hyperparameters() map[string]string
}

func NewFinetuneParameter(model, task string, hyperparameters map[string]string) (
	FinetuneParameter, error,
) {
	cfg, ok := DomainConfig.Finetunes[model]
	if !ok {
		return nil, errors.New("invalid model")
	}

	// task
	bingo := false
	for _, t := range cfg.Tasks {
		if t == task {
			bingo = true

			break
		}
	}
	if !bingo {
		return nil, errors.New("invalid task")
	}

	// check: if invalid key
	f := func(key string) bool {
		for i := range cfg.Hyperparameters {
			if key == cfg.Hyperparameters[i] {
				return true
			}
		}

		return false
	}

	for k, v := range hyperparameters {
		if !f(k) {
			return nil, errors.New("invalid hyperparameter")
		}

		// check: if invalid value
		if !(utils.IsPositiveFloatPoint(v) ||
			utils.IsPositiveScientificNotation(v) ||
			utils.IsPositiveInterger(v)) {
			return nil, errors.New("invalid hyperparameter")
		}
	}

	return finetuneParameter{
		model:           model,
		task:            task,
		hyperparameters: hyperparameters,
	}, nil
}

type finetuneParameter struct {
	model           string
	task            string
	hyperparameters map[string]string
}

func (p finetuneParameter) Model() string {
	return p.model
}

func (p finetuneParameter) Task() string {
	return p.task
}

func (p finetuneParameter) Hyperparameters() map[string]string {
	return p.hyperparameters
}

type FinetuneJob struct {
	FinetuneJobInfo

	Status string
}

type FinetuneJobInfo struct {
	JobId    string
	Endpoint string
}

type FinetuneJobDetail struct {
	Error    string
	Status   string
	Duration int
}

type FinetuneSummary struct {
	Id        string
	Name      FinetuneName
	CreatedAt int64

	FinetuneJobDetail
}

type FinetuneUserInfo struct {
	Expiry int64
}

type FinetuneIndex struct {
	Id    string
	Owner Account
}

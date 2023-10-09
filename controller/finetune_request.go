package controller

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type finetuneCreateResp = trainingCreateResp

type finetuneLog struct {
	Log string `json:"log"`
}

type FinetuneCreateRequest struct {
	Name            string     `json:"name"`
	Model           string     `json:"model"`
	Task            string     `json:"task"`
	Hyperparameters []KeyValue `json:"hyperparameter"`
}

func (req *FinetuneCreateRequest) toCmd(user domain.Account) (
	cmd app.FinetuneCreateCmd, err error,
) {
	cmd.User = user

	if cmd.Name, err = domain.NewTrainingName(req.Name); err != nil {
		return
	}

	m := map[string]string{}
	if len(req.Hyperparameters) > 0 {
		for i := range req.Hyperparameters {
			if req.Hyperparameters[i].Value == "" {
				continue
			}

			item := &req.Hyperparameters[i]
			m[item.Key] = item.Value
		}
	}

	cmd.Param, err = domain.NewFinetuneParameter(req.Model, req.Task, m)
	if err == nil {
		err = cmd.Validate()
	}

	return
}

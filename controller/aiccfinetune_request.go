package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/aiccfinetune/app"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
)

type aiccFinetuneCreateResp struct {
	Id string `json:"id"`
}

type aiccFinetuneLogResp struct {
	LogURL string `json:"log_url"`
}

type aiccFinetuneDetail struct {
	app.AICCFinetuneDTO
	Log string `json:"log"`
}

type aiccFinetuneCreateRequest struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Task string `json:"task"`

	Hyperparameters []AICCKeyValue `json:"hyperparameter"`
	Env             []AICCKeyValue `json:"env"`
}

func (req *aiccFinetuneCreateRequest) toCmd(cmd *app.AICCFinetuneCreateCmd, model string) (err error) {
	if cmd.Name, err = domain.NewFinetuneName(req.Name); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewFinetuneDesc(req.Desc); err != nil {
		return
	}

	if cmd.Task, err = domain.NewFinetuneTask(req.Task); err != nil {
		return
	}

	if cmd.Model, err = domain.NewModelName(model); err != nil {
		return
	}

	if cmd.Env, err = req.toKeyValue(req.Env); err != nil {
		return
	}

	if cmd.Hyperparameters, err = req.toKeyValue(req.Hyperparameters); err != nil {
		return
	}

	return
}

func (req *aiccFinetuneCreateRequest) toKeyValue(kv []AICCKeyValue) (r []domain.KeyValue, err error) {
	n := len(kv)
	if n == 0 {
		return nil, nil
	}

	r = make([]domain.KeyValue, n)
	for i := range kv {
		if r[i], err = kv[i].toKeyValue(); err != nil {
			return
		}
	}

	return
}

type AICCKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (kv *AICCKeyValue) toKeyValue() (r domain.KeyValue, err error) {
	if kv.Key == "" {
		err = errors.New("invalid key value")

		return
	}

	if r.Key, err = domain.NewCustomizedKey(kv.Key); err != nil {
		return
	}

	r.Value, err = domain.NewCustomizedValue(kv.Value)

	return
}

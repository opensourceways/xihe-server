package app

import (
	"errors"
	"io"

	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type AICCFinetuneCreateCmd struct {
	User  types.Account
	Model domain.ModelName
	Task  domain.FinetuneTask

	domain.AICCFinetuneConfig
}

func (cmd *AICCFinetuneCreateCmd) Validate() error {
	err := errors.New("invalid cmd of creating aicc finetune")

	b := cmd.User != nil &&
		cmd.Name != nil

	if !b {
		return err
	}

	f := func(kv []domain.KeyValue) error {
		for i := range kv {
			if kv[i].Key == nil {
				return err
			}
		}

		return nil
	}

	if f(cmd.Hyperparameters) != nil {
		return err
	}

	if f(cmd.Env) != nil {
		return err
	}

	return nil
}

func (cmd *AICCFinetuneCreateCmd) toAICCFinetuneConfig() *domain.AICCFinetuneConfig {
	return &cmd.AICCFinetuneConfig
}

type AICCFinetuneSummaryDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Error     string `json:"error"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	IsDone    bool   `json:"is_done"`
	Duration  int    `json:"duration"`
	Task      string `json:"task"`
}

func (s aiccFinetuneService) toAICCFinetuneSummaryDTO(
	t *domain.AICCFinetuneSummary, dto *AICCFinetuneSummaryDTO,
) {
	status := t.Status
	if status == "" {
		status = trainingStatusScheduling
	}

	*dto = AICCFinetuneSummaryDTO{
		Id:        t.Id,
		Name:      t.Name.FinetuneName(),
		Error:     t.Error,
		Status:    status,
		IsDone:    s.isJobDone(t.Status),
		Duration:  t.Duration,
		CreatedAt: utils.ToDate(t.CreatedAt),
		Task:      t.Task,
	}

	if t.Desc != nil {
		dto.Desc = t.Desc.FinetuneDesc()
	}
}

type AICCFinetuneDTO struct {
	Id string `json:"id"`

	Name string `json:"name"`
	Desc string `json:"desc"`

	IsDone    bool   `json:"is_done"`
	Error     string `json:"error"`
	Status    string `json:"status"`
	Duration  int    `json:"duration"`
	CreatedAt string `json:"created_at"`

	LogPreviewURL string `json:"-"`
}

type ComputeDTO struct {
	ImageUrl string `json:"image_url"`
}

func (s aiccFinetuneService) toAICCFinetuneDTO(dto *AICCFinetuneDTO, ut *domain.AICCFinetune, link string) {
	t := &ut.AICCFinetuneConfig
	detail := &ut.JobDetail

	status := detail.Status
	if status == "" {
		status = trainingStatusScheduling
	}

	*dto = AICCFinetuneDTO{
		Id: ut.Id,

		Name:          t.Name.FinetuneName(),
		IsDone:        s.isJobDone(detail.Status),
		Error:         detail.Error,
		Status:        status,
		Duration:      detail.Duration,
		CreatedAt:     utils.ToDate(ut.CreatedAt),
		LogPreviewURL: link,
	}

	if t.Desc != nil {
		dto.Desc = t.Desc.FinetuneDesc()
	}
}

type UploadDataCmd struct {
	FileName string
	Data     io.Reader
	User     types.Account
	Model    domain.ModelName
	Task     domain.FinetuneTask
}

type UploadDataDTO struct {
	UploadAt int64  `json:"upload_at"`
	FileName string `json:"file_name"`
	Status   string `json:"status"`
}

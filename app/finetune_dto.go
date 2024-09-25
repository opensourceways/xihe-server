package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type FinetuneCreateCmd struct {
	User domain.Account

	domain.FinetuneConfig
}

func (cmd *FinetuneCreateCmd) Validate() error {
	b := cmd.User != nil &&
		cmd.Name != nil &&
		cmd.Param != nil

	if !b {
		return errors.New("invalid cmd of creating finetune")
	}

	return nil
}

func (cmd *FinetuneCreateCmd) toFinetuneConfig() *domain.FinetuneConfig {
	return &cmd.FinetuneConfig
}

type FinetuneSummaryDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Error     string `json:"error"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	IsDone    bool   `json:"is_done"`
	Duration  int    `json:"duration"`
}

func (s finetuneService) toFinetuneSummaryDTO(
	t *domain.FinetuneSummary, dto *FinetuneSummaryDTO,
) {
	status := t.Status
	if status == "" {
		status = trainingStatusScheduling
	}

	dto.Id = t.Id
	dto.Name = t.Name.FinetuneName()
	dto.Error = t.Error
	dto.Status = status
	dto.IsDone = s.isJobDone(status)
	dto.Duration = t.Duration
	dto.CreatedAt = utils.ToDate(t.CreatedAt)
}

type UserFinetunesDTO struct {
	Expiry int64                `json:"expiry"`
	Data   []FinetuneSummaryDTO `json:"data"`
}

type FinetuneJobDTO struct {
	IsDone        bool
	LogPreviewURL string
}

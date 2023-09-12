package app

import (
	"time"

	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type CmdToAddPointsItem struct {
	Account common.Account
	Task    string
	Desc    string
	Time    int64
}

func (cmd *CmdToAddPointsItem) dateAndTime() (string, string) {
	now := time.Now().Unix()

	if cmd.Time > now || cmd.Time < (now-minValueOfInvlidTime) {
		return "", ""
	}

	return utils.DateAndTime(cmd.Time)
}

type UserPointsDetailsDTO struct {
	Total   int               `json:"total"`
	Details []PointsDetailDTO `json:"details"`
}

type PointsDetailDTO struct {
	Task string `json:"task"`

	domain.PointsDetail
}

// TasksCompletionInfoDTO
type TasksCompletionInfoDTO struct {
	Kind  string                  `json:"kind"`
	Tasks []TaskCompletionInfoDTO `json:"tasks"`
}

func (dto *TasksCompletionInfoDTO) add(t *domain.Task, completed bool) {
	dto.Tasks = append(dto.Tasks, TaskCompletionInfoDTO{
		Name:      t.Name,
		Addr:      t.Addr,
		Points:    t.Rule.PointsPerOnce,
		Completed: completed,
	})
}

type TaskCompletionInfoDTO struct {
	Name      string `json:"name"`
	Addr      string `json:"addr"`
	Points    int    `json:"points"`
	Completed bool   `json:"completed"`
}

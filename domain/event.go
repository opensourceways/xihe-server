package domain

import (
	"encoding/json"
	"fmt"

	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/utils"
)

type Msg comsg.MsgNormal

// training
func NewCreateTrainingMsg(
	user Account,
	index TrainingIndex,
	inputs []Input,
) *Msg {
	desc := fmt.Sprintf("create training, id: %s", index.TrainingId)
	bytes, _ := json.Marshal(inputs)

	return &Msg{
		User: user.Account(),
		Desc: desc,
		Details: map[string]string{
			"project_owner": index.Project.Owner.Account(),
			"project_id":    index.Project.Id,
			"training_id":   index.TrainingId,
			"input":         string(bytes),
		},
		CreatedAt: utils.Now(),
	}
}

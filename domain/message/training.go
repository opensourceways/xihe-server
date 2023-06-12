package message

import (
	"encoding/json"

	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	MsgTypeTraningCreate = "msg_type_training_create"
)

type MsgTraining comsg.MsgNormal

func (msg *MsgTraining) TrainingCreate(
	user domain.Account,
	index domain.TrainingIndex,
	inputs []domain.Input,
) {

	bytes, err := json.Marshal(inputs)
	if err != nil {
		return
	}

	in := string(bytes)

	*msg = MsgTraining{
		Type: MsgTypeTraningCreate,
		User: user.Account(),
		Details: map[string]string{
			"project_owner": index.Project.Owner.Account(),
			"project_id":    index.Project.Id,
			"training_id":   index.TrainingId,
			"input":         in,
		},
		CreatedAt: utils.Now(),
	}
}

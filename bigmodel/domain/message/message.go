package message

import (
	"github.com/opensourceways/xihe-server/domain/message"
)

type MsgTask struct {
	Type    string            `json:"type"`
	TaskId  uint64            `json:"task_id"`
	User    string            `json:"user"`
	Status  string            `json:"status"`
	Details map[string]string `json:"details"`
}

type AsyncMessageProducer interface {
	message.Sender

	CreateWuKongTask(*MsgTask) error

	UpdateWuKongTask(*MsgTask) error
}

func (msg *MsgTask) ToMsgTask(user, desc, style string) {
	*msg = MsgTask{
		Type:   "wukong_request",
		User:   user,
		Status: "waiting",
		Details: map[string]string{
			"style": style,
			"desc":  desc,
		},
	}
}

func (msg *MsgTask) SetErrorMsgTask(tid uint64, user, errMsg string) {
	*msg = MsgTask{
		Type:   "wukong_update",
		TaskId: tid,
		User:   user,
		Status: "error",
		Details: map[string]string{
			"error": errMsg,
		},
	}
}

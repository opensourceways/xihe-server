package message

import (
	"strconv"
	"strings"

	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain/message"
)

const (
	MsgTypeWuKongInferenceStart  = "msg_type_wukong_inference_start"
	MsgTypeWuKongInferenceError  = "msg_type_wukong_inference_error"
	MsgTypeWuKongAsyncTaskStart  = "msg_type_wukong_async_task_start"
	MsgTypeWuKongAsyncTaskFinish = "msg_type_wukong_async_task_finish"
	MsgTypeWuKongPublicPicture   = "msg_type_wukong_public_picture"
)

// normal msessage
type MsgTask comsg.MsgNormal

type AsyncMessageProducer interface {
	message.Sender

	SendBigModelMsg(*MsgTask) error
	SendBigmodelPublicMsg(*MsgTask) error
}

func (msg *MsgTask) WuKongInferenceStart(user, desc, style, taskType string) {
	*msg = MsgTask{
		Type: MsgTypeWuKongInferenceStart,
		User: user,
		Details: map[string]string{
			"status":    "waiting",
			"task_type": taskType,
			"style":     style,
			"desc":      desc,
		},
	}
}

func (msg *MsgTask) WuKongInferenceError(tid uint64, user, errMsg string) {
	*msg = MsgTask{
		Type: MsgTypeWuKongInferenceError,
		User: user,
		Details: map[string]string{
			"task_id": strconv.Itoa(int(tid)),
			"status":  "error",
			"error":   errMsg,
		},
	}
}

func (msg *MsgTask) WuKongAsyncTaskStart(tid uint64, user string) {
	*msg = MsgTask{
		Type: MsgTypeWuKongAsyncTaskStart,
		User: user,
		Details: map[string]string{
			"task_id": strconv.Itoa(int(tid)),
			"status":  "running",
		},
	}
}

func (msg *MsgTask) WuKongAsyncInferenceFinish(tid uint64, user string, links map[string]string) {
	var ls string
	for k := range links { // TODO: Move it into domain.service
		ls += links[k] + ","
	}

	*msg = MsgTask{
		Type: MsgTypeWuKongAsyncTaskFinish,
		User: user,
		Details: map[string]string{
			"task_id": strconv.Itoa(int(tid)),
			"status":  "finished",
			"links":   strings.TrimRight(ls, ","),
		},
	}
}

func (msg *MsgTask) GenWuKongPublicPictureMsg(user string) {
	*msg = MsgTask{
		Type:    MsgTypeWuKongPublicPicture,
		User:    user,
		Desc:    "",
		Details: nil,
	}
}

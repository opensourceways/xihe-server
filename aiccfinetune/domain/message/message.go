package message

import (
	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	comsg "github.com/opensourceways/xihe-server/common/domain/message"
)

const (
	MsgTypeAICCFinetuneCreate = "msg_type_aicc_finetune_create"
)

type MsgTask comsg.MsgNormal

type AICCFinetuneMessageProducer interface {
	SendAICCFinetuneCreateMsg(*domain.AICCFinetuneCreateEvent) error
}

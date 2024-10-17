package message

import (
	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

const (
	MsgTypeAICCFinetuneCreate = "msg_type_aicc_finetune_create"
)

type MsgTask comsg.MsgNormal

type SpaceAppMessageProducer interface {
	SendSpaceAppCreateMsg(*domain.SpaceAppCreateEvent) error
}

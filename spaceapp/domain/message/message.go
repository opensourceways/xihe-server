package message

import (
	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

type MsgTask comsg.MsgNormal

type SpaceAppMessageProducer interface {
	SendSpaceAppCreateMsg(*domain.SpaceAppCreateEvent) error
}

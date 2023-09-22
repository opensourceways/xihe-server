package message

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	comsg "github.com/opensourceways/xihe-server/common/domain/message"
)

type MsgTask comsg.MsgNormal

type MessageProducer interface {
	// wukong
	SendWuKongInferenceStart(*domain.WuKongInferenceStartEvent) error
	SendWuKongInferenceError(*domain.WuKongInferenceErrorEvent) error
	SendWuKongAsyncTaskStart(*domain.WuKongAsyncTaskStartEvent) error
	SendWuKongAsyncInferenceFinish(*domain.WuKongAsyncInferenceFinishEvent) error
	SendWuKongPicturePublicized(*domain.WuKongPicturePublicizedEvent) error
	SendWuKongPictureLiked(*domain.WuKongPictureLikedEvent) error

	// common
	SendBigModelStarted(*domain.BigModelStartedEvent) error
}

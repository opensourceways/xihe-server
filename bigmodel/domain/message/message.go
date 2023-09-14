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
	SendUserPublicizeWuKongPicture(*domain.UserPublicizeWuKongPictureEvent) error

	// common
	SendBigModelAccessLog(*domain.BigModelAccessLogEvent) error
}

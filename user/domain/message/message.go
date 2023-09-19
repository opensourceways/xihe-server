package message

import "github.com/opensourceways/xihe-server/user/domain"

type MessageProducer interface {
	SendUserSignedUpEvent(*domain.UserSignedUpEvent) error
	SendUserAvatarSetEvent(event *domain.UserAvatarSetEvent) error
	SendUserBioSetEvent(event *domain.UserBioSetEvent) error
}

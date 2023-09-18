package message

import "github.com/opensourceways/xihe-server/user/domain"

type MessageProducer interface {
	SendUserRegisterEvent(*domain.UserSignedUpEvent) error
	SendSetAvatarIdEvent(event *domain.UserAvatarSetEvent) error
	SendSetBioEvent(event *domain.UserBioSetEvent) error
}

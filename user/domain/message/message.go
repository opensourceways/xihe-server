package message

import "github.com/opensourceways/xihe-server/user/domain"

type MessageProducer interface {
	SendUserRegisterEvent(*domain.UserRegisterEvent) error
	SendBindEmailEvent(event *domain.UserBindEmailEvent) error
	SendSetAvatarIdEvent(event *domain.UserSetAvatarIdEvent) error
	SendSetBioEvent(event *domain.UserSetBioEvent) error
}

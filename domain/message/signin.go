package message

import "github.com/opensourceways/xihe-server/domain"

type UserSignedInMessageProducer interface {
	SendUserSignedIn(*domain.UserSignedInEvent) error
}

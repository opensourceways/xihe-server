package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type UserSignedUpEvent struct {
	Account types.Account
}

type UserAvatarSetEvent struct {
	Account  types.Account
	AvatarId AvatarId
}

type UserBioSetEvent struct {
	Account types.Account
	Bio     Bio
}

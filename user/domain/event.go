package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

type UserRegisterEvent struct {
	Account types.Account
}

type UserBindEmailEvent struct {
	Account types.Account
	Email   types.Email
}

type UserSetAvatarIdEvent struct {
	Account  types.Account
	AvatarId AvatarId
}

type UserSetBioEvent struct {
	Account types.Account
	Bio     Bio
}

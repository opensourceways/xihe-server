package authing

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserInfo struct {
	Name     domain.Account
	Email    domain.Email
	Bio      domain.Bio
	AvatarId domain.AvatarId
	UserId   string
}

type Login struct {
	UserInfo

	IDToken     string
	AccessToken string
}

type User interface {
	GetByCode(code, redirectURI string) (Login, error)
	GetByAccessToken(accessToken string) (UserInfo, error)

	// email
	SendBindEmail(accessToken, email string) (err error)
	VerifyBindEmail(accessToken, email, passCode string) (err error)
}

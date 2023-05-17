package authing

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserInfo struct {
	Name     domain.Account
	Email    domain.Email
	Bio      domain.Bio
	AvatarId domain.AvatarId
}

type Login struct {
	UserInfo

	IDToken string
	AccessToken string
}

type User interface {
	GetByCode(code, redirectURI string) (Login, error)
	GetByAccessToken(accessToken string) (UserInfo, error)
}

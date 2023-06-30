package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

// user
type UserCreateCmd struct {
	Email    domain.Email
	Account  domain.Account
	Password domain.Password

	Bio      domain.Bio
	AvatarId domain.AvatarId
}

func (cmd *UserCreateCmd) Validate() error {
	b := cmd.Email != nil &&
		cmd.Account != nil &&
		cmd.Password != nil

	if !b {
		return errors.New("invalid cmd of creating user")
	}

	return nil
}

func (cmd *UserCreateCmd) toUser() domain.User {
	return domain.User{
		Email:   cmd.Email,
		Account: cmd.Account,

		Bio:      cmd.Bio,
		AvatarId: cmd.AvatarId,
	}
}

type UserDTO struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Account string `json:"account"`

	Bio      string `json:"bio"`
	AvatarId string `json:"avatar_id"`

	FollowerCount  int `json:"follower_count"`
	FollowingCount int `json:"following_count"`

	Platform struct {
		UserId      string
		Token       string
		NamespaceId string
	} `json:"-"`
}

type UpdateUserBasicInfoCmd struct {
	Bio      domain.Bio
	Email    domain.Email
	AvatarId domain.AvatarId
}

func (cmd *UpdateUserBasicInfoCmd) toUser(u *domain.User) (changed bool) {
	if cmd.AvatarId != nil && !domain.IsSameDomainValue(cmd.AvatarId, u.AvatarId) {
		u.AvatarId = cmd.AvatarId
		changed = true
	}

	if cmd.Bio != nil && !domain.IsSameDomainValue(cmd.Bio, u.Bio) {
		u.Bio = cmd.Bio
		changed = true
	}

	if cmd.Email != nil && u.Email.Email() != cmd.Email.Email() {
		u.Email = cmd.Email
		changed = true
	}

	return
}

type FollowsListCmd struct {
	User domain.Account

	repository.FollowFindOption
}

type FollowsDTO struct {
	Total int         `json:"total"`
	Data  []FollowDTO `json:"data"`
}

type FollowDTO struct {
	Account    string `json:"account"`
	AvatarId   string `json:"avatar_id"`
	Bio        string `json:"bio"`
	IsFollower bool   `json:"is_follower"`
}

// register
type UserRegisterInfoCmd domain.UserRegInfo

func (cmd *UserRegisterInfoCmd) toUserRegInfo(r *domain.UserRegInfo) {
	*r = *(*domain.UserRegInfo)(cmd)
}

type UserRegisterInfoDTO domain.UserRegInfo

func (dto *UserRegisterInfoDTO) toUserRegInfoDTO(r *domain.UserRegInfo) {
	*dto = *(*UserRegisterInfoDTO)(r)
}

type SendBindEmailCmd struct {
	User  domain.Account
	Email domain.Email
	Capt  string
}

type BindEmailCmd struct {
	User     domain.Account
	Email    domain.Email
	PassCode string
	PassWord domain.Password
}

type CreatePlatformAccountCmd struct {
	Email    domain.Email
	Account  domain.Account
	Password domain.Password
}

type PlatformInfoDTO struct {
	PlatformUser  domain.PlatformUser
	PlatformToken string
}

type UpdatePlateformInfoCmd struct {
	PlatformInfoDTO

	User  domain.Account
	Email domain.Email
}

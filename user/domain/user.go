package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
)

// user
type User struct {
	Id      string
	Email   Email
	Account Account

	Bio      Bio
	AvatarId AvatarId

	PlatformUser  PlatformUser
	PlatformToken string

	Version int

	// following fileds is not under the controlling of version
	FollowerCount  int
	FollowingCount int
}

type PlatformUser struct {
	Id          string
	NamespaceId string
}

type FollowerInfo struct {
	User     Account
	Follower Account
}

type FollowerUserInfo struct {
	Account    Account
	AvatarId   AvatarId
	Bio        Bio
	IsFollower bool
}

type UserInfo struct {
	Account  Account
	AvatarId AvatarId
}

// register
type UserRegInfo struct {
	Account  types.Account
	Name     Name
	City     City
	Email    Email
	Phone    Phone
	Identity Identity
	Province Province
	Detail   map[string]string
	Version  int
}

type UserPoints struct {
	Account Account
	Type    string
	Desc    string
	Time    int64
}

package domain

import (
	"time"

	types "github.com/opensourceways/xihe-server/domain"
)

// user
type User struct {
	Id      string
	Email   Email
	Phone   Phone
	Account Account

	Bio      Bio
	AvatarId AvatarId

	PlatformUser  PlatformUser
	PlatformToken PlatformToken

	CourseAgreement   string
	FinetuneAgreement string
	UserAgreement     string

	Version int

	// following fields is not under the controlling of version
	FollowerCount  int
	FollowingCount int

	IsAgreePrivacy bool
}

type PlatformUser struct {
	Id          string
	NamespaceId string
}

type PlatformToken struct {
	Token    string `json:"token"`
	CreateAt int64  `json:"create_at"`
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

type ModifyAccountInfo struct {
	AccountType string `json:"account_type"`
	OldAccount  string `json:"oldaccount"`
	OldCode     string `json:"oldcode"`
	Account     string `json:"account"`
	Code        string `json:"code"`
}

// whitelist
type WhiteListInfo struct {
	Account   types.Account
	Type      WhiteListType
	Enabled   bool
	StartTime int64
	EndTime   int64
}

func (w WhiteListInfo) Enable() bool {
	t := time.Now()

	return w.Enabled && t.After(time.Unix(w.StartTime, 0)) && t.Before(time.Unix(w.EndTime, 0))
}

func (u *User) RevokePrivacy() {
	u.IsAgreePrivacy = false
}

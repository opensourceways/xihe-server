package repository

import (
	"github.com/opensourceways/xihe-server/user/domain"
)

type FollowFindOption struct {
	Follower domain.Account

	CountPerPage int
	PageNum      int
}

type FollowerUserInfos struct {
	Users []domain.FollowerUserInfo
	Total int
}

type UserSearchOption struct {
	// can't define Name as domain.Account
	// because the Name can be subpart of the real account
	Name   string
	TopNum int
}

type UserSearchResult struct {
	Top []domain.Account

	Total int
}

type User interface {
	Save(*domain.User) (domain.User, error)
	GetByAccount(domain.Account) (domain.User, error)
	GetByFollower(owner, follower domain.Account) (domain.User, bool, error)
	FindUsersInfo([]domain.Account) ([]domain.UserInfo, error)
	GetUserAvatarId(domain.Account) (domain.AvatarId, error)

	Search(*UserSearchOption) (UserSearchResult, error)

	AddFollowing(*domain.FollowerInfo) error
	RemoveFollowing(*domain.FollowerInfo) error
	FindFollowing(domain.Account, *FollowFindOption) (FollowerUserInfos, error)

	AddFollower(*domain.FollowerInfo) error
	RemoveFollower(*domain.FollowerInfo) error
	FindFollower(domain.Account, *FollowFindOption) (FollowerUserInfos, error)
}

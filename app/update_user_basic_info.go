package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure"
)

type UpdateUserBasicInfoCmd struct {
	NickName domain.Nickname
	AvatarId domain.AvatarId
	Bio      domain.Bio
}

func (cmd *UpdateUserBasicInfoCmd) validate() error {
	if cmd.NickName == nil || cmd.AvatarId == nil || cmd.Bio == nil {
		return errors.New("invalid cmd of updating user's basic info")
	}

	return nil
}

func (cmd *UpdateUserBasicInfoCmd) toUser(u *domain.User) (changed bool) {
	set := func() {
		if !changed {
			changed = true
		}
	}
	if cmd.NickName.Nickname() != u.Nickname.Nickname() {
		u.Nickname = cmd.NickName
		set()
	}

	if cmd.AvatarId.AvatarId() != u.AvatarId.AvatarId() {
		u.AvatarId = cmd.AvatarId
		set()
	}

	if cmd.Bio.Bio() != u.Bio.Bio() {
		u.Bio = cmd.Bio
		set()
	}

	return
}

type UserService interface {
	UpdateBasicInfo(userId string, cmd UpdateUserBasicInfoCmd) error
	RecordLikeProject(userId, projectId string) error
}

func NewUserService() UserService {
	var repo infrastructure.UserMapper
	return userService{repo}
}

type userService struct {
	repo infrastructure.UserMapper
}

func (s userService) UpdateBasicInfo(userId string, cmd UpdateUserBasicInfoCmd) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	user, err := s.repo.Get(userId)
	if err != nil {
		return err
	}

	// if b := cmd.toUser(&user); !b {
	// 	return nil
	// }

	return s.repo.Save(user)
}
func (s userService) RecordLikeProject(userId, projectId string) error {
	s.repo.LikeProject(projectId)
	return nil

}

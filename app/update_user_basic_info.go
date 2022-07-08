package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
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

	if cmd.NickName.Nickname() != cmd.NickName.Nickname() {
		u.Nickname = cmd.NickName
		set()
	}

	if cmd.AvatarId.AvatarId() != u.AvatarId.AvatarId() {
		u.AvatarId = cmd.AvatarId
		set()
	}

	if cmd.Bio.Bio() != cmd.Bio.Bio() {
		u.Bio = cmd.Bio
		set()
	}

	return
}

type UserService interface {
	UpdateBasicInfo(userId string, cmd UpdateUserBasicInfoCmd) error
}

func NewUserService(repo repository.User) UserService {
	return userService{repo}
}

type userService struct {
	repo repository.User
}

func (s userService) UpdateBasicInfo(userId string, cmd UpdateUserBasicInfoCmd) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	user, err := s.repo.Get(userId)
	if err != nil {
		return err
	}

	if b := cmd.toUser(&user); !b {
		return nil
	}

	return s.repo.Save(user)
}

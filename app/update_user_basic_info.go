package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

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

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return err
	}

	if b := cmd.toUser(&user); !b {
		return nil
	}

	_, err = s.repo.Save(&user)

	return err
}

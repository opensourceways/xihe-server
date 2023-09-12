package app

import "github.com/opensourceways/xihe-server/user/domain"

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return err
	}

	if b := cmd.toUser(&user); !b {
		return nil
	}

	_, err = s.repo.Save(&user)
	if cmd.AvatarId != nil {
		return s.sender.SetAvatarId(account, cmd.AvatarId)
	}
	if cmd.Bio != nil {
		return s.sender.SetBio(account, cmd.Bio)
	}

	return err
}

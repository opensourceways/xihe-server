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

	if _, err = s.repo.Save(&user); err != nil {
		return err
	}

	if cmd.avatarChanged {
		_ = s.sender.SendUserAvatarSetEvent(&domain.UserAvatarSetEvent{
			Account:  account,
			AvatarId: cmd.AvatarId,
		})
	}

	if cmd.bioChanged {
		_ = s.sender.SendUserBioSetEvent(&domain.UserBioSetEvent{
			Account: account,
			Bio:     cmd.Bio,
		})
	}

	return nil
}

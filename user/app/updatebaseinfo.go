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

	if cmd.avatarChanged == true {
		return s.sender.SendUserAvatarSetEvent(&domain.UserAvatarSetEvent{
			Account:  account,
			AvatarId: cmd.AvatarId,
		})
	}

	if cmd.bioChanged == true {
		return s.sender.SendUserBioSetEvent(&domain.UserBioSetEvent{
			Account: account,
			Bio:     cmd.Bio,
		})
	}

	return err
}

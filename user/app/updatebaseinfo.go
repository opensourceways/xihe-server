package app

import "github.com/opensourceways/xihe-server/user/domain"

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return err
	}

	b := cmd.toUser(&user)
	if !b {
		return nil
	}

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

	_, err = s.repo.Save(&user)

	return err

}

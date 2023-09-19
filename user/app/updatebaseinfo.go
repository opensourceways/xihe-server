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

	_, err = s.repo.Save(&user)
	if err != nil {
		return err
	} else if b && cmd.AvatarId != nil && cmd.AvatarId != user.AvatarId {
		return s.producer.SendUserAvatarSetEvent(&domain.UserAvatarSetEvent{
			Account:  account,
			AvatarId: cmd.AvatarId,
		})
	} else if b && cmd.Bio != nil && cmd.Bio != user.Bio {
		return s.producer.SendUserBioSetEvent(&domain.UserBioSetEvent{
			Account: account,
			Bio:     cmd.Bio,
		})
	}
	return err
}

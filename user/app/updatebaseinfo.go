package app

import "github.com/opensourceways/xihe-server/user/domain"

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return err
	}

	if b := cmd.toUser(&user); !b {
		return nil
	} else if b && cmd.AvatarId != nil && cmd.AvatarId != user.AvatarId {
		return s.producer.SendSetAvatarIdEvent(&domain.UserAvatarSetEvent{
			Account:  account,
			AvatarId: cmd.AvatarId,
		})

	} else if b && cmd.Bio != nil && cmd.Bio != user.Bio {
		return s.producer.SendSetBioEvent(&domain.UserBioSetEvent{
			Account: account,
			Bio:     cmd.Bio,
		})
	}

	_, err = s.repo.Save(&user)

	return err
}

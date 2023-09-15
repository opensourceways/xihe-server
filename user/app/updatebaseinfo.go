package app

import "github.com/opensourceways/xihe-server/user/domain"

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return err
	}

	if b := cmd.toUser(&user); !b {
		return nil
	} else if b && cmd.AvatarId != nil {
		return s.producer.SendSetAvatarIdEvent(&domain.UserSetAvatarIdEvent{
			Account:  account,
			AvatarId: cmd.AvatarId,
		})

	} else if b && cmd.Bio != nil {
		return s.producer.SendSetBioEvent(&domain.UserSetBioEvent{
			Account: account,
			Bio:     cmd.Bio,
		})
	}

	_, err = s.repo.Save(&user)

	return err
}

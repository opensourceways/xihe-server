package app

import (
	"github.com/opensourceways/xihe-server/user/domain"
)

const (
	auditTitle = "title"
)

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	//sdk text audit
	bio := cmd.Bio.Bio()
	if bio != "" {
		if err := s.audit.TextAudit(bio, auditTitle); err != nil {
			return err
		}
	}

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

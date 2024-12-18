package app

import (
	"golang.org/x/xerrors"

	auditapi "github.com/opensourceways/xihe-audit-sync-sdk/audit/api"
	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/user/domain"
)

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	//sdk text audit
	bio := cmd.Bio.Bio()
	if bio != "" {
		resp, _, err := auditapi.Text(bio, "profile")
		if err != nil {
			return allerror.New(
				allerror.ErrorCodeFailToModerate,
				resp.Result, err)
		} else if resp.Result != "pass" {
			e := xerrors.Errorf("moderate unpass")
			return allerror.New(
				allerror.ErrorCodeModerateUnpass,
				resp.Result, e)
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

package controller

import (
	"github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/user/domain"
)

type userBasicInfoUpdateRequest struct {
	AvatarId *string `json:"avatar_id"`
	Bio      *string `json:"bio"`
}

func (req *userBasicInfoUpdateRequest) toCmd() (
	cmd app.UpdateUserBasicInfoCmd,
	err error,
) {
	if req.Bio != nil {
		if cmd.Bio, err = domain.NewBio(*req.Bio); err != nil {
			return
		}
	}

	if req.AvatarId != nil {
		cmd.AvatarId, err = domain.NewAvatarId(*req.AvatarId)
	}

	return
}

type userCreateRequest struct {
	Account  string `json:"account"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	AvatarId string `json:"avatar_id"`
}

func (req *userCreateRequest) toCmd() (cmd app.UserCreateCmd, err error) {
	if cmd.Account, err = domain.NewAccount(req.Account); err != nil {
		return
	}

	if cmd.Email, err = domain.NewEmail(req.Email); err != nil {
		return
	}

	if cmd.Bio, err = domain.NewBio(req.Bio); err != nil {
		return
	}

	if cmd.AvatarId, err = domain.NewAvatarId(req.AvatarId); err != nil {
		return
	}

	if cmd.Password, err = domain.NewPassword(apiConfig.DefaultPassword); err != nil {
		return
	}

	err = cmd.Validate()

	return
}

type followingCreateRequest struct {
	Account string `json:"account" required:"true"`
}

type userDetail struct {
	*app.UserDTO

	IsFollower bool `json:"is_follower"`
}

type EmailCode struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (req *EmailCode) toCmd(user domain.Account) (cmd app.BindEmailCmd, err error) {
	if cmd.Email, err = domain.NewEmail(req.Email); err != nil {
		return
	}

	cmd.PassCode = req.Code
	cmd.User = user

	if cmd.PassWord, err = domain.NewPassword(apiConfig.DefaultPassword); err != nil {
		return
	}

	return
}

type EmailSend struct {
	Email string `json:"email"`
	Capt  string `json:"capt"`
}

func (req *EmailSend) toCmd(user domain.Account) (cmd app.SendBindEmailCmd, err error) {
	if cmd.Email, err = domain.NewEmail(req.Email); err != nil {
		return
	}

	cmd.User = user
	cmd.Capt = req.Capt

	return
}

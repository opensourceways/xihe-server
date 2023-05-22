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

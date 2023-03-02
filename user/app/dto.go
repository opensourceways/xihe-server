package app

import (
	"github.com/opensourceways/xihe-server/user/domain"
)

type UserRegisterInfoCmd domain.UserRegInfo

func (cmd *UserRegisterInfoCmd) toUserRegInfo(r *domain.UserRegInfo) {
	*r = *(*domain.UserRegInfo)(cmd)
}

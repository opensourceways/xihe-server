package app

import "github.com/opensourceways/xihe-server/domain"

func (s loginService) toLoginDTO(u *domain.Login, dto *LoginDTO) {
	dto.Info = u.Info
	dto.Email = u.Email.Email()
	dto.UserId = u.UserId
}

package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/course/domain"
	userApp "github.com/opensourceways/xihe-server/user/app"
)

// player apply
type PlayerApplyCmd domain.Player

func (cmd *PlayerApplyCmd) Validate() error {
	b := cmd.Student.Account != nil &&
		cmd.Student.Name != nil &&
		cmd.Student.Email != nil &&
		cmd.Student.Identity != nil

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *PlayerApplyCmd) toPlayer() (p domain.Player) {
	return *(*domain.Player)(cmd)
}

// userRegInfo
type StudentInfoDTO = userApp.UserRegisterInfoCmd

func toStudentDTO(s *domain.Student) (dto *StudentInfoDTO) {
	return &StudentInfoDTO{
		Account:  s.Account.Account(),
		Name:     s.Name.StudentName(),
		City:     s.City.City(),
		Email:    s.Email.Email(),
		Phone:    s.Phone.Phone(),
		Identity: s.Identity.StudentIdentity(),
		Province: s.Province.Province(),
		Detail:   s.Detail,
	}
}

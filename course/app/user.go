package app

import (
	"errors"

	userDomain "github.com/opensourceways/xihe-server/user/domain"
)

func (s *courseService) Apply(cmd *PlayerApplyCmd) (err error) {
	course, err := s.courseRepo.FindCourse(cmd.CourseId)
	if err != nil {
		return
	}

	if course.IsOver() {
		err = errors.New("course is over")

		return
	}

	if course.IsPreliminary() {
		err = errors.New("course is preparing")

		return
	}

	p := cmd.toPlayer()
	p.CreateToday()

	// save player
	if err = s.playerRepo.SavePlayer(&p); err != nil {
		return
	}

	// save register info
	var regInfo userDomain.UserRegInfo
	if err = regInfo.NewRegFromStudent(&p.Student); err != nil {
		return
	}
	if err = s.userRepo.AddUserRegInfo(&regInfo); err != nil {
		return
	}

	return
}

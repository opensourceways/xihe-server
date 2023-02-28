package app

import (
	"errors"
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
	p.NewId()

	if err = s.playerRepo.SavePlayer(&p); err != nil {
		return
	}

	if err = s.user.AddUserRegInfo(&p.Student); err != nil {
		return
	}

	return
}

package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/course/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

func (s *courseService) Apply(cmd *PlayerApplyCmd) (code string, err error) {
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
		if repoerr.IsErrorDuplicateCreating(err) {
			code = errorDuplicateApply
		}

		return
	}

	if err = s.userCli.AddUserRegInfo(&p.Student); err != nil {
		return
	}

	return
}

func (s *courseService) AddReleatedProject(cmd *CourseAddReleatedProjectCmd) (
	code string, err error,
) {
	// check phase
	course, err := s.courseRepo.FindCourse(cmd.Cid)
	if course.IsOver() {
		err = errors.New("course is over")

		return
	}

	if course.IsPreliminary() {
		err = errors.New("course is preparing")

		return
	}
	// check permission
	player, err := s.playerRepo.FindPlayer(cmd.Cid, cmd.User)

	if !course.IsApplyed(&player.Player) {
		code = errorNoPermission
		return
	}

	if cmd.Project.Owner != cmd.User {
		code = errorDoesnotOwnProject
		err = errors.New("the user does not own the project")
		return
	}

	repo := domain.NewCourseProject(cmd.User, cmd.repo())

	err = s.playerRepo.SaveRepo(cmd.Cid, &repo, player.Version)

	return
}

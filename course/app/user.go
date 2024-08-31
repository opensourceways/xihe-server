package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/agreement/app"
	"github.com/opensourceways/xihe-server/course/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/sirupsen/logrus"
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

	if err = p.CreateToday(); err != nil {
		return
	}

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

	ver := app.GetCurrentCourseAgree()
	user, err := s.userRepo.GetByAccount(cmd.Account)
	if err != nil {
		return
	}

	if user.CourseAgreement != ver {
		user.CourseAgreement = ver
		logrus.Debugf("User %s user agreement updated from %s to %s ",
			user.Account.Account(), user.CourseAgreement, ver)
		if _, err = s.userRepo.Save(&user); err != nil {
			return
		}
	}

	if err = s.producer.SendCourseAppliedEvent(&domain.CourseAppliedEvent{
		Account:    cmd.Account,
		CourseName: course.Name,
	}); err != nil {
		return
	}

	return
}

func (s *courseService) AddRelatedProject(cmd *CourseAddRelatedProjectCmd) (
	code string, err error,
) {
	// check phase
	course, err := s.courseRepo.FindCourse(cmd.Cid)
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
	// check permission
	player, err := s.playerRepo.FindPlayer(cmd.Cid, cmd.User)

	if !course.IsApplied(&player.Player) {
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

func (s *courseService) GetCertification(cmd *CourseGetCmd) (dto CertInfoDTO, err error) {
	p, err := s.playerRepo.FindPlayer(cmd.Cid, cmd.User)
	if err != nil {
		return
	}
	c, err := s.courseRepo.FindCourse(cmd.Cid)
	if err != nil {
		return
	}
	if !c.IsApplied(&p.Player) {
		return
	}

	asg, err := s.courseRepo.FindAssignments(cmd.Cid)
	if err != nil {
		return
	}

	var score float32
	for i := range asg {
		w, err := s.workRepo.GetWork(cmd.Cid, cmd.User, asg[i].Id, nil)
		if err != nil {
			break
		}
		score += w.Score
	}

	var pass bool
	if score >= c.PassScore.CoursePassScore() {
		pass = true
	}

	toCertInfoDTO(cmd.User, &c, pass, &dto)

	return
}

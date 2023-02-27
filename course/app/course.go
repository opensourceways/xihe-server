package app

import (
	"github.com/opensourceways/xihe-server/course/domain/repository"
	"github.com/opensourceways/xihe-server/course/domain/user"
)

type CourseService interface {
	// player
	Apply(*PlayerApplyCmd) error
}

func NewCourseService(
	userCli user.User,

	courseRepo repository.Course,
	playerRepo repository.Player,
) *courseService {
	return &courseService{
		userCli:    userCli,
		courseRepo: courseRepo,
		playerRepo: playerRepo,
	}
}

type courseService struct {
	userCli user.User

	courseRepo repository.Course
	playerRepo repository.Player
}

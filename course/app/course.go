package app

import (
	"github.com/opensourceways/xihe-server/course/domain/repository"
	"github.com/opensourceways/xihe-server/course/domain/user"
)

type CourseService interface {
	// player
	Apply(*PlayerApplyCmd) error

	// course
}

func NewCourseService(
	user user.User,
	courseRepo repository.Course,
	playerRepo repository.Player,
) *courseService {
	return &courseService{
		user:       user,
		courseRepo: courseRepo,
		playerRepo: playerRepo,
	}
}

type courseService struct {
	user       user.User
	courseRepo repository.Course
	playerRepo repository.Player
}

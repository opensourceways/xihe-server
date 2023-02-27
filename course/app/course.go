package app

import (
	"github.com/opensourceways/xihe-server/course/domain/repository"
	user "github.com/opensourceways/xihe-server/user/domain/repository"
)

type CourseService interface {
	// player
	Apply(*PlayerApplyCmd) error

	// course
}

func NewCourseService(
	userRepo user.User,
	courseRepo repository.Course,
	playerRepo repository.Player,
) *courseService {
	return &courseService{
		userRepo:   userRepo,
		courseRepo: courseRepo,
		playerRepo: playerRepo,
	}
}

type courseService struct {
	userRepo user.User

	courseRepo repository.Course
	playerRepo repository.Player
}

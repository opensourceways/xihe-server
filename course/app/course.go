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
	userregRepo user.UserReg,
	courseRepo repository.Course,
	playerRepo repository.Player,
) *courseService {
	return &courseService{
		userRepo:   userregRepo,
		courseRepo: courseRepo,
		playerRepo: playerRepo,
	}
}

type courseService struct {
	userRepo user.UserReg

	courseRepo repository.Course
	playerRepo repository.Player
}

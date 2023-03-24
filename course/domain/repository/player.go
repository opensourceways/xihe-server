package repository

import (
	"github.com/opensourceways/xihe-server/course/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type PlayerVersion struct {
	domain.Player
	Version int
}

type Player interface {
	FindPlayer(cid string, account types.Account) (PlayerVersion, error)
	SavePlayer(*domain.Player) error
	PlayerCount(cid string) (int, error)
	SaveRepo(course_id string, a *domain.CourseProject, version int) error
	FindCoursesUserApplied(types.Account) ([]string, error)
}
